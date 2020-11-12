package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-storage/pkg/digest"
	bpb "google.golang.org/genproto/googleapis/bytestream"
	"google.golang.org/grpc"
)

var blake bool

func parseFlags() (string, string) {
	var address string
	var inputRoot string

	flag.StringVar(&address, "address", "", "Address of gRPC endpoint for Buildbarn frontend")
	flag.StringVar(&address, "a", "", "Address of gRPC endpoint for Buildbarn frontend (shorthand)")
	flag.BoolVar(&blake, "blake", true, "True for BLAKE3ZCC, false for SHA256")
	flag.BoolVar(&blake, "b", true, "True for BLAKE3ZCC, false for SHA256 (shorthand)")
	flag.StringVar(&inputRoot, "directory", "", "The directory to be the input root of the action")
	flag.StringVar(&inputRoot, "d", "", "The directory to be the input root of the action (shorthand)")
	flag.Parse()

	return address, inputRoot
}

func main() {
	address, inputRoot := parseFlags()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	casClient := remoteexecution.NewContentAddressableStorageClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	blobUploader := blobUploader.NewBlobUploader(conn, 100)
	instanceName, err := digest.NewInstanceName("blake-client")
	if err != nil {
		log.Fatalf("Error creating instance name: %v", err)
	}
	hash := sha256.New()
	if blake {
		// Random digest to get access to hash function
		d, err := instanceName.NewDigest(
			"B3Z:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			123)
		if err != nil {
			log.Fatalf("Error creating digest for hash function access: %v", err)
		}
		hash = d.NewHasher()
	}
	file, err := os.Open("/home/scott/buildbarn/bb-storage/randomfile")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("%v", err)
	}

	b := make([]byte, stat.Size())
	_, err = file.Read(b)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalf("%v", err)
	}

	var digest digest.Digest
	if blake {
		digest, err = instanceName.NewDigest(fmt.Sprintf("B3Z:%s", hex.EncodeToString(hash.Sum(nil))), stat.Size())
	} else {
		digest, err = instanceName.NewDigest(hex.EncodeToString(hash.Sum(nil)), stat.Size())
	}
	if err != nil {
		log.Fatalf("Error creating digest for file: %v", err)
	}
	//log.Printf("Digest %v", digest)

	findMissingRequest := &remoteexecution.FindMissingBlobsRequest{
		InstanceName: instanceName.String(),
		BlobDigests:  []*remoteexecution.Digest{digest.GetProto()},
	}

	findMissingResponse, err := casClient.FindMissingBlobs(ctx, findMissingRequest)
	if err != nil {
		log.Fatalf("Error with FindMissingBlobs: %v", err)
	}
	//log.Printf("Find Missing Response: %v", findMissingResponse)

	if len(findMissingResponse.MissingBlobDigests) != 0 {
		updateRequest := &remoteexecution.BatchUpdateBlobsRequest{
			InstanceName: instanceName.String(),
			Requests: []*remoteexecution.BatchUpdateBlobsRequest_Request{{
				Digest: digest.GetProto(),
				Data:   b,
			}},
		}

		_, err = casClient.BatchUpdateBlobs(ctx, updateRequest)
		if err != nil {
			log.Fatalf("Error with BatchUpdateBlobs: %v", err)
		}
	}

	//log.Printf("Response %v", updateResponse)

	//readResponse, err := casClient.BatchReadBlobs(ctx, &remoteexecution.BatchReadBlobsRequest{
	//InstanceName: "test",
	//Digests:      []*remoteexecution.Digest{digest.GetProto()},
	//})
	//if err != nil {
	//log.Fatalf("Error with BatchReadBlobs: %v", err)
	//}

	//newFile, err := os.Create("/home/scott/output.png")
	//if err != nil {
	//log.Fatalf("Failed to create output file: %v", err)
	//}

	//_, err = newFile.Write(readResponse.Responses[0].Data)
	//if err != nil {
	//log.Fatalf("Failed to write to output file: %v", err)
	//}
	//if manifestDigest, _, ok := digest.ToManifest(int64(8 * 1024)); ok {
	//log.Printf("Manifest Size: %v", manifestDigest.GetSizeBytes())
	//}
}
