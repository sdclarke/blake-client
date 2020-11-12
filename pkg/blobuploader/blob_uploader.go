package blobuploader

import (
	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	bpb "google.golang.org/genproto/googleapis/bytestream"
)

type blob struct {
	digest *remoteexecution.Digest
	b      []byte
}

type blobUploader struct {
	bytestreamClient bpb.ByteStreamClient
	blobs            []*blob
	maxSize          int
}

func NewBlobUploader(conn grpc.ClientConnInterface, maxSize int) (*blobUploader, func(context.Context) error) {
	bu := &blobUploader{
		bytestreamClient: bpb.NewByteStreamClient(conn),
		blobs:            []*blob{},
		maxSize:          maxSize,
	}
	return bu, func(ctx context.Context) error {
		return bu.findMissingAndUpload(ctx)
	}
}

func (bu *blobUploader) Add(digest *remoteexecution.Digest, b []byte) error {
	bu.blobs = append(bu.blobs, &blob{digest, b})
	if len(bu.blobs) == maxSize {
		err := bu.findMissingAndUpload()
		if err != nil {
			return err
		}
		bu.blobs = []*blob{}
	}
	return nil
}

func (bu *blobUploader) findMissingAndUpload() error {
	return nil
}
