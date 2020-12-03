package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"log"

	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/sdclarke/blake-client/pkg/blobuploader"
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
	//casClient := remoteexecution.NewContentAddressableStorageClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	blobUploader, finaliser := blobuploader.NewBlobUploader(conn, instanceName.String(), 100, hash, blake)

	inputRootDigest, err := blobUploader.UploadDirectory(ctx, inputRoot)
	if err != nil {
		log.Fatalf("Something went wrong uploading input root: %v", err)
	}
	err = finaliser(ctx)
	if err != nil {
		log.Fatalf("Something went wrong uploading input root: %v", err)
	}
	log.Printf("Input root digest: %v", inputRootDigest)
}
