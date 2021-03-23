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
	"time"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-storage/pkg/digest"
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
	bytestreamClient   bpb.ByteStreamClient
	casClient          remoteexecution.ContentAddressableStorageClient
	blobs              map[string]*blob
	instanceName       digest.InstanceName
	maxSize            int
	hash               hash.Hash
	blake              bool
	decompose          bool
	decomposeSize      int64
	bytesUploaded      int64
	bytesHashed        int64
	timeHashing        int64
	timeUploading      int64
	timeFindingMissing int64
	missingBlobs       int64
	totalBlobs         int64
}

func NewBlobUploader(conn grpc.ClientConnInterface, instanceName digest.InstanceName, maxSize int, hash hash.Hash, blake, decompose bool, decomposeSize int) (*blobUploader, func(context.Context) error) {
	bu := &blobUploader{
		bytestreamClient: bpb.NewByteStreamClient(conn),
		casClient:        remoteexecution.NewContentAddressableStorageClient(conn),
		blobs:            make(map[string]*blob),
		instanceName:     instanceName,
		maxSize:          maxSize,
		hash:             hash,
		blake:            blake,
		decompose:        decompose,
		decomposeSize:    int64(decomposeSize),
	}
	return bu, func(ctx context.Context) error {
		return bu.findMissingAndUpload(ctx)
	}
}

func (bu *blobUploader) GetBytesUploaded() int64 {
	return bu.bytesUploaded
}

func (bu *blobUploader) GetBytesHashed() int64 {
	return bu.bytesHashed
}

func (bu *blobUploader) GetTimeHashing() int64 {
	return bu.timeHashing
}

func (bu *blobUploader) GetTimeUploading() int64 {
	return bu.timeUploading
}

func (bu *blobUploader) GetTimeFindingMissing() int64 {
	return bu.timeFindingMissing
}

func (bu *blobUploader) GetMissingBlobs() int64 {
	return bu.missingBlobs
}

func (bu *blobUploader) GetTotalBlobs() int64 {
	return bu.totalBlobs
}

func (bu *blobUploader) add(ctx context.Context, digest *remoteexecution.Digest, b []byte) error {
	hash := digest.GetHashOther()
	if hash == "" {
		if len(digest.GetHashBlake3Zcc()) != 0 {
			hash = fmt.Sprintf("B3Z:%s", hex.EncodeToString(digest.GetHashBlake3Zcc()))
		} else {
			hash = fmt.Sprintf("B3ZM:%s", hex.EncodeToString(digest.GetHashBlake3ZccManifest()))
		}
	}
	//log.Printf("Hash: %v", hash)
	bu.totalBlobs++
	if _, ok := bu.blobs[hash]; ok {
		return nil
	}
	bu.blobs[hash] = &blob{digest: digest, b: b}
	return nil
}

func (bu *blobUploader) Add(ctx context.Context, digest *remoteexecution.Digest, b []byte) error {
	if err := bu.add(ctx, digest, b); err != nil {
		return err
	}
	if len(bu.blobs) == bu.maxSize {
		err := bu.findMissingAndUpload(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bu *blobUploader) AddProto(ctx context.Context, m proto.Message) (*remoteexecution.Digest, error) {
	digest, bytes, err := bu.protoToDigest(m)
	if err != nil {
		return nil, err
	}
	return digest, bu.Add(ctx, digest, bytes)
}

func (bu *blobUploader) findMissingAndUpload(ctx context.Context) error {
	findMissingRequest := &remoteexecution.FindMissingBlobsRequest{
		InstanceName: bu.instanceName.String(),
		BlobDigests:  []*remoteexecution.Digest{},
	}
	for _, blob := range bu.blobs {
		if bu.decompose {
			bbDigest, err := bu.instanceName.NewDigestFromProto(blob.digest)
			if err != nil {
				return err
			}
			if bbManifestDigest, manifestParser, ok := bbDigest.ToManifest(bu.decomposeSize); ok {
				bu.totalBlobs--
				manifestSize := bbManifestDigest.GetSizeBytes()
				if manifestSize > (9 * 1024 * 1024) {
					return status.Errorf(codes.InvalidArgument,
						"Buffer requires a manifest that is %d bytes in size, while a maximum of %d bytes is permitted",
						manifestSize,
						9*1024*1024)
				}

				manifest := make([]byte, 0, manifestSize)
				data := blob.b
				for {
					if len(data) == 0 {
						break
					}
					block := make([]byte, int(bu.decomposeSize))
					n := copy(block, data)
					data = data[n:]
					block = block[:n]
					timeBefore := time.Now().UnixNano()
					bbBlockDigest := manifestParser.AppendBlockDigest(&manifest, block)
					timeAfter := time.Now().UnixNano()
					bu.timeHashing += (timeAfter - timeBefore)
					blockDigest := bbBlockDigest.GetProto()
					bu.add(ctx, blockDigest, block)
					findMissingRequest.BlobDigests = append(findMissingRequest.BlobDigests, blockDigest)
				}
				manifestDigest := bbManifestDigest.GetProto()
				bu.add(ctx, manifestDigest, manifest)
				findMissingRequest.BlobDigests = append(findMissingRequest.BlobDigests, manifestDigest)
			} else {
				findMissingRequest.BlobDigests = append(findMissingRequest.BlobDigests, blob.digest)
			}
		} else {
			findMissingRequest.BlobDigests = append(findMissingRequest.BlobDigests, blob.digest)
		}
	}
	timeBeforeFindMissing := time.Now().UnixNano()
	findMissingResponse, err := bu.casClient.FindMissingBlobs(ctx, findMissingRequest)
	if err != nil {
		return err
	}
	timeAfterFindMissing := time.Now().UnixNano()
	bu.timeFindingMissing += (timeAfterFindMissing - timeBeforeFindMissing)
	missing := findMissingResponse.GetMissingBlobDigests()
	bu.missingBlobs += int64(len(missing))
	for _, digest := range missing {
		wr, err := bu.bytestreamClient.Write(ctx)
		if err != nil {
			return err
		}
		uuid := uuid.New()

		size := digest.GetSizeBytes()
		hash := digest.GetHashOther()
		if hash == "" {
			if len(digest.GetHashBlake3Zcc()) != 0 {
				hash = fmt.Sprintf("B3Z:%s", hex.EncodeToString(digest.GetHashBlake3Zcc()))
			} else {
				hash = fmt.Sprintf("B3ZM:%s", hex.EncodeToString(digest.GetHashBlake3ZccManifest()))
			}
		}

		resourceNameEnd := fmt.Sprintf("blobs/%s/%d", hash, size)
		resourceName := path.Join(bu.instanceName.String(), "uploads", uuid.String(), resourceNameEnd)

		writeOffset := int64(0)
		blobBytes := bu.blobs[hash].b
		bytes := make([]byte, maxSize)
		n := copy(bytes, blobBytes)
		timeBefore := time.Now().UnixNano()
		for n > 0 {
			// Write request for non-terminating chunk
			writeRequest := &bpb.WriteRequest{
				ResourceName: resourceName,
				WriteOffset:  writeOffset,
				FinishWrite:  false,
				Data:         bytes[:n],
			}

			if err = wr.Send(writeRequest); err != nil {
				//writeResponse, err := wr.CloseAndRecv()
				//log.Printf("Write Response: %v", writeResponse)
				return err
			}
			resourceName = ""
			writeOffset += int64(n)
			blobBytes = blobBytes[n:]
			n = copy(bytes, blobBytes)
		}
		// Write request for terminating chunk
		writeRequest := &bpb.WriteRequest{
			ResourceName: resourceName,
			WriteOffset:  writeOffset,
			FinishWrite:  true,
		}
		if err = wr.Send(writeRequest); err != nil {
			return err
		}
		writeResponse, err := wr.CloseAndRecv()
		if err != nil {
			return err
		}
		timeAfter := time.Now().UnixNano()
		bu.timeUploading += (timeAfter - timeBefore)
		committedSize := writeResponse.GetCommittedSize()
		if committedSize < size {
			return status.Errorf(codes.Unknown, "Committed size was %d, expected %d", committedSize, size)
		}
		bu.bytesUploaded += committedSize
		delete(bu.blobs, hash)
	}
	bu.blobs = make(map[string]*blob)
	return nil
}

func (bu *blobUploader) UploadDirectory(ctx context.Context, directory string) (*remoteexecution.Digest, error) {
	dir, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
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
	inputRootDigest, bytes, err := bu.protoToDigest(directory)
	if err != nil {
		return nil, err
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
	bu.bytesHashed += stat.Size()

	bu.hash.Reset()

	timeBefore := time.Now().UnixNano()
	size, err := io.Copy(bu.hash, file)
	if err != nil {
		return nil, err
	}
	hashBytes := bu.hash.Sum(nil)
	timeAfter := time.Now().UnixNano()
	bu.timeHashing += (timeAfter - timeBefore)
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

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

func (bu *blobUploader) protoToDigest(m proto.Message) (*remoteexecution.Digest, []byte, error) {
	bytes, err := proto.Marshal(m)
	if err != nil {
		return nil, nil, err
	}
	bu.bytesHashed += int64(len(bytes))
	bu.hash.Reset()
	timeBefore := time.Now().UnixNano()
	_, err = bu.hash.Write(bytes)
	if err != nil {
		return nil, nil, err
	}
	hashBytes := bu.hash.Sum(nil)
	timeAfter := time.Now().UnixNano()
	bu.timeHashing += (timeAfter - timeBefore)
	if bu.blake {
		return &remoteexecution.Digest{
			HashBlake3Zcc: hashBytes,
			SizeBytes:     int64(len(bytes)),
		}, bytes, nil
	}
	return &remoteexecution.Digest{
		HashOther: hex.EncodeToString(hashBytes),
		SizeBytes: int64(len(bytes)),
	}, bytes, nil
}
