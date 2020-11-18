package blobuploader

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"strings"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	bpb "google.golang.org/genproto/googleapis/bytestream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxSize = 65536

type blob struct {
	digest *remoteexecution.Digest
	b      []byte
}

type blobUploader struct {
	bytestreamClient bpb.ByteStreamClient
	casClient        remoteexecution.ContentAddressableStorageClient
	blobs            map[string]*blob
	instanceName     string
	maxSize          int
	hash             hash.Hash
	blake            bool
}

func NewBlobUploader(conn grpc.ClientConnInterface, instanceName string, maxSize int, hash hash.Hash, blake bool) (*blobUploader, func(context.Context) error) {
	bu := &blobUploader{
		bytestreamClient: bpb.NewByteStreamClient(conn),
		casClient:        remoteexecution.NewContentAddressableStorageClient(conn),
		blobs:            make(map[string]*blob),
		instanceName:     instanceName,
		maxSize:          maxSize,
		hash:             hash,
		blake:            blake,
	}
	return bu, func(ctx context.Context) error {
		return bu.findMissingAndUpload(ctx)
	}
}

func (bu *blobUploader) Add(ctx context.Context, digest *remoteexecution.Digest, b []byte) error {
	hash := digest.GetHashOther()
	if hash == "" {
		hash = fmt.Sprintf("B3Z:%s", hex.EncodeToString(digest.GetHashBlake3Zcc()))
	}
	if _, ok := bu.blobs[hash]; ok {
		return nil
	}
	bu.blobs[hash] = &blob{digest: digest, b: b}
	if len(bu.blobs) == bu.maxSize {
		err := bu.findMissingAndUpload(ctx)
		if err != nil {
			return err
		}
		bu.blobs = make(map[string]*blob)
	}
	return nil
}

func (bu *blobUploader) findMissingAndUpload(ctx context.Context) error {
	findMissingRequest := &remoteexecution.FindMissingBlobsRequest{
		InstanceName: bu.instanceName,
		BlobDigests:  []*remoteexecution.Digest{},
	}
	for _, blob := range bu.blobs {
		findMissingRequest.BlobDigests = append(findMissingRequest.BlobDigests, blob.digest)
	}
	findMissingResponse, err := bu.casClient.FindMissingBlobs(ctx, findMissingRequest)
	if err != nil {
		return err
	}
	missing := findMissingResponse.GetMissingBlobDigests()
	for _, digest := range missing {
		wr, err := bu.bytestreamClient.Write(ctx)
		if err != nil {
			return err
		}
		uuid := uuid.New()

		size := digest.GetSizeBytes()
		hash := digest.GetHashOther()
		if hash == "" {
			hash = fmt.Sprintf("B3Z:%s", hex.EncodeToString(digest.GetHashBlake3Zcc()))
		}

		resourceNameEnd := fmt.Sprintf("blobs/%s/%d", hash, size)
		resourceName := path.Join(bu.instanceName, "uploads", uuid.String(), resourceNameEnd)

		writeOffset := int64(0)
		blobBytes := bu.blobs[hash].b
		for {
			bytes := make([]byte, maxSize)
			n := copy(bytes, blobBytes)
			for n > 0 {
				// Write request for non-terminating chunk of file
				writeRequest_file := &bpb.WriteRequest{
					ResourceName: resourceName,
					WriteOffset:  writeOffset,
					FinishWrite:  false,
					Data:         bytes[:n],
				}

				if err = wr.Send(writeRequest_file); err != nil {
					return err
				}
				resourceName = ""
				writeOffset += int64(n)
				blobBytes = blobBytes[n:]
				n = copy(bytes, blobBytes)
			}
			// Write request for terminating chunk of file
			writeRequest_file := &bpb.WriteRequest{
				ResourceName: resourceName,
				WriteOffset:  writeOffset,
				FinishWrite:  true,
			}
			if err = wr.Send(writeRequest_file); err != nil {
				return err
			}
			break
		}
		writeResponse, err := wr.CloseAndRecv()
		if err != nil {
			return err
		}
		if committedSize := writeResponse.GetCommittedSize(); committedSize < size {
			return status.Errorf(codes.Unknown, "Committed size was %d, expected %d", committedSize, size)
		}
	}
	return nil
}

func (bu *blobUploader) UploadDirectory(ctx context.Context, directory string) (*remoteexecution.Digest, error) {
	dir, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	fileNames, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	files := []*remoteexecution.FileNode{}
	directories := []*remoteexecution.DirectoryNode{}
	for _, fileName := range fileNames {
		node, err := bu.createFileNode(ctx, path.Join(directory, fileName.Name()))
		if err != nil {
			childDirectoryDigest, err := bu.UploadDirectory(ctx, path.Join(directory, fileName.Name()))
			if err != nil {
				return nil, err
			}
			directories = append(directories, &remoteexecution.DirectoryNode{
				Name:   fileName.Name(),
				Digest: childDirectoryDigest,
			})
		} else {
			files = append(files, node)
		}
	}
	return bu.createDirectory(ctx, files, directories)
}

func (bu *blobUploader) createDirectory(ctx context.Context, files []*remoteexecution.FileNode, directories []*remoteexecution.DirectoryNode) (*remoteexecution.Digest, error) {
	directory := &remoteexecution.Directory{
		Files:       files,
		Directories: directories,
	}
	bytes, err := proto.Marshal(directory)
	if err != nil {
		return nil, err
	}
	bu.hash.Reset()
	_, err = bu.hash.Write(bytes)
	if err != nil {
		return nil, err
	}
	hashBytes := bu.hash.Sum(nil)
	var inputRootDigest *remoteexecution.Digest
	if bu.blake {
		inputRootDigest = &remoteexecution.Digest{
			HashBlake3Zcc: hashBytes,
			SizeBytes:     int64(len(bytes)),
		}
	} else {
		inputRootDigest = &remoteexecution.Digest{
			HashOther: hex.EncodeToString(hashBytes),
			SizeBytes: int64(len(bytes)),
		}
	}
	bu.Add(ctx, inputRootDigest, bytes)
	return inputRootDigest, nil
}

func (bu *blobUploader) createFileNode(ctx context.Context, fileName string) (*remoteexecution.FileNode, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, status.Error(codes.InvalidArgument, "File is a directory")
	}

	bu.hash.Reset()

	size, err := io.Copy(bu.hash, file)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	hashBytes := bu.hash.Sum(nil)

	fileNode := &remoteexecution.FileNode{
		Name:         stat.Name(),
		IsExecutable: strings.Contains(stat.Mode().Perm().String(), "x"),
	}
	if bu.blake {
		fileNode.Digest = &remoteexecution.Digest{
			HashBlake3Zcc: hashBytes,
			SizeBytes:     int64(size),
		}
	} else {
		fileNode.Digest = &remoteexecution.Digest{
			HashOther: hex.EncodeToString(hashBytes),
			SizeBytes: int64(size),
		}
	}
	bytes := make([]byte, size)
	_, err = io.ReadFull(file, bytes)
	if err != nil {
		return nil, err
	}
	bu.Add(ctx, fileNode.Digest, bytes)
	return fileNode, nil
}
