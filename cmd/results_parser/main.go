package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseFlags() string {
	var file string

	flag.StringVar(&file, "file", "", "File from which to parse results")
	flag.StringVar(&file, "f", "", "File from which to parse results")
	flag.Parse()

	return file
}

func main() {
	fileName := parseFlags()
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)
	text := []string{}
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	f.Close()

	cumulativeHashTime := int64(0)
	cumulativeUploadTime := int64(0)
	cumulativeFindMissingTime := int64(0)
	cumulativeExecuteTime := int64(0)
	cumulativeUploadTimeServer := int64(0)
	totalRuns := 0
	totalHashes := 0
	totalUploaded := int64(0)
	totalHashed := int64(0)
	for _, line := range text {
		if strings.HasPrefix(line, "Bytes") {
			tokens := strings.Split(line, " ")
			switch tokens[1] {
			case "Uploaded:":
				uploaded, err := strconv.ParseInt(tokens[2], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing bytes uploaded: %v", err)
				}
				totalUploaded += uploaded
				break
			case "Hashed:":
				hashed, err := strconv.ParseInt(tokens[2], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing bytes hashed: %v", err)
				}
				totalHashed += hashed
				break
			}

		} else if strings.HasPrefix(line, "Time") {
			tokens := strings.Split(line, " ")
			switch tokens[1] {
			case "Hashing:":
				d, err := time.ParseDuration(tokens[2])
				if err != nil {
					log.Fatalf("Error parsing duration: %v", err)
				}
				cumulativeHashTime += d.Nanoseconds()
				totalHashes++
				break
			case "Uploading:":
				d, err := time.ParseDuration(tokens[2])
				if err != nil {
					log.Fatalf("Error parsing duration: %v", err)
				}
				cumulativeUploadTime += d.Nanoseconds()
				break
			case "Finding":
				d, err := time.ParseDuration(tokens[4])
				if err != nil {
					log.Fatalf("Error parsing duration: %v", err)
				}
				cumulativeFindMissingTime += d.Nanoseconds()
				break
			}
		} else {
			tokens := strings.Split(line, " ")
			switch tokens[0] {
			case "Execute:":
				pieces := strings.Split(tokens[1], "s")
				executeStartSec, err := strconv.ParseInt(pieces[0], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				executeStartNSec, err := strconv.ParseInt(pieces[1], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				executeStart := time.Unix(executeStartSec, executeStartNSec)
				pieces = strings.Split(tokens[2], "s")
				executeEndSec, err := strconv.ParseInt(pieces[0], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				executeEndNSec, err := strconv.ParseInt(pieces[1], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				executeEnd := time.Unix(executeEndSec, executeEndNSec)
				cumulativeExecuteTime += (executeEnd.UnixNano() - executeStart.UnixNano())
				totalRuns++
				break
			case "Upload:":
				pieces := strings.Split(tokens[1], "s")
				uploadStartSec, err := strconv.ParseInt(pieces[0], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				uploadStartNSec, err := strconv.ParseInt(pieces[1], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				uploadStart := time.Unix(uploadStartSec, uploadStartNSec)
				pieces = strings.Split(tokens[2], "s")
				uploadEndSec, err := strconv.ParseInt(pieces[0], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				uploadEndNSec, err := strconv.ParseInt(pieces[1], 10, 64)
				if err != nil {
					log.Fatalf("Error parsing time: %v", err)
				}
				uploadEnd := time.Unix(uploadEndSec, uploadEndNSec)
				cumulativeUploadTimeServer += (uploadEnd.UnixNano() - uploadStart.UnixNano())
				break
			}
		}
	}
	hashTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeHashTime)/float64(totalHashes)))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	uploadTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeUploadTime)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	findMissingTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeFindMissingTime)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	executeTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeExecuteTime)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	uploadTimeServer, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeUploadTimeServer)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}
	fmt.Printf("Client:\nHash time: %v\nUpload Time: %v\nFind Missing Time: %v\n", hashTime, uploadTime, findMissingTime)
	fmt.Printf("Bytes Hashed: %v\nBytes Uploaded: %v\n", totalHashed, totalUploaded)
	fmt.Printf("Server:\nExecute time: %v\nUpload Time: %v\n", executeTime, uploadTimeServer)
}
