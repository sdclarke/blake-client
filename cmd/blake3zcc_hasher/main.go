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
	"path"
	"time"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/golang/protobuf/ptypes"
	"github.com/sdclarke/blake-client/pkg/blobuploader"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	commands = []*remoteexecution.Command{
		{
			Arguments:   []string{"cp", "in0", "out"},
			OutputPaths: []string{"out"},
		},
		{
			Arguments:   []string{"sh", "-c", "cat in0 in0 > out"},
			OutputPaths: []string{"out"},
		},
		{
			Arguments:   []string{"sh", "-c", "cat in0 in1 > out"},
			OutputPaths: []string{"out"},
		},
		{
			Arguments:   []string{"sh", "-c", "cp in0 out && truncate -r in1 out"},
			OutputPaths: []string{"out"},
		},
		{
			Arguments:   []string{"sh", "-c", "cp in0 out && printf ' ' | dd of=out bs=1 count=1 seek=%d conv=notrunc"},
			OutputPaths: []string{"out"},
		},
	}
)

var blake bool

func parseFlags() (string, string, int) {
	var address string
	var inputRoot string
	var commandNumber int

	flag.StringVar(&address, "address", "", "Address of gRPC endpoint for Buildbarn frontend")
	flag.StringVar(&address, "a", "", "Address of gRPC endpoint for Buildbarn frontend (shorthand)")
	flag.BoolVar(&blake, "blake", true, "True for BLAKE3ZCC, false for SHA256")
	flag.BoolVar(&blake, "b", true, "True for BLAKE3ZCC, false for SHA256 (shorthand)")
	flag.StringVar(&inputRoot, "directory", "", "The directory to be the input root of the action")
	flag.StringVar(&inputRoot, "d", "", "The directory to be the input root of the action (shorthand)")
	flag.IntVar(&commandNumber, "command", 0, fmt.Sprintf("The number of the command to be run (0-%d)", len(commands)-1))
	flag.IntVar(&commandNumber, "c", 0, fmt.Sprintf("The number of the command to be run (0-%d)", len(commands)-1))
	flag.Parse()

	return address, inputRoot, commandNumber
}

func main() {
	address, inputRoot, commandNumber := parseFlags()
	if commandNumber >= len(commands) || commandNumber < 0 {
		log.Fatalf("Command number is invalid, must be in range 0-%d", len(commands)-1)
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
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
	if err != nil {
		log.Fatalf("Something went wrong uploading input root: %v", err)
	}

	command := commands[commandNumber]
	if commandNumber == 4 {
		file, err := os.Open(path.Join(inputRoot, "in0"))
		if err != nil {
			log.Fatalf("Error opening input file: %v", err)
		}

		stat, err := file.Stat()
		if err != nil {
			log.Fatalf("Error stating input file: %v", err)
		}
		file.Close()
		command.Arguments[2] = fmt.Sprintf(command.Arguments[2], stat.Size()/2)
		log.Printf("Command: %v", command)
	}
	commandDigest, err := blobUploader.AddProto(ctx, command)
	if err != nil {
		log.Fatalf("Error uploading command: %v", err)
	}

	action := createAction(commandDigest, inputRootDigest)
	actionDigest, err := blobUploader.AddProto(ctx, action)
	if err != nil {
		log.Fatalf("Error uploading action: %v", err)
	}

	hashingDuration, err := time.ParseDuration(fmt.Sprintf("%dns", blobUploader.GetTimeHashing()))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	log.Printf("Time Hashing: %v", hashingDuration)
	log.Printf("Bytes Hashed: %v", blobUploader.GetBytesHashed())
	err = finaliser(ctx)
	if err != nil {
		log.Fatalf("Error finalising blob uploads %v", err)
	}
	log.Printf("Action Digest: %v %v %v %d", actionDigest.GetHashBlake3Zcc(), hex.EncodeToString(actionDigest.GetHashBlake3Zcc()), actionDigest.GetHashOther(), actionDigest.GetSizeBytes())
	uploadingDuration, err := time.ParseDuration(fmt.Sprintf("%dns", blobUploader.GetTimeUploading()))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	findingMissingDuration, err := time.ParseDuration(fmt.Sprintf("%dns", blobUploader.GetTimeFindingMissing()))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	log.Printf("Bytes Uploaded: %v", blobUploader.GetBytesUploaded())
	log.Printf("Time Uploading: %v", uploadingDuration)
	log.Printf("Time Finding Missing Blobs: %v", findingMissingDuration)

	executionClient := remoteexecution.NewExecutionClient(conn)

	stream, err := executionClient.Execute(ctx, &remoteexecution.ExecuteRequest{
		InstanceName: instanceName.String(),
		ActionDigest: actionDigest,
	})
	if err != nil {
		log.Fatalf("Error calling Execute(): %v", err)
	}

	executeResponse, err := receiveResults(stream)
	if err != nil {
		log.Fatalf("Error receiving results: %v", err)
	}
	log.Printf("Execute response: %v", executeResponse)
}

func createAction(commandDigest, inputRootDigest *remoteexecution.Digest) *remoteexecution.Action {
	return &remoteexecution.Action{
		CommandDigest:   commandDigest,
		InputRootDigest: inputRootDigest,
	}
}

func receiveResults(stream remoteexecution.Execution_ExecuteClient) (*remoteexecution.ExecuteResponse, error) {
	resp := &remoteexecution.ExecuteResponse{}
	for {
		a, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if a.GetDone() {
			err = ptypes.UnmarshalAny(a.GetResponse(), resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
	return nil, status.Errorf(codes.Internal, "Something is on fire")
}
