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
					log.Fatalf("Error parsing hashing duration: %v, %v", tokens[2], err)
				}
				cumulativeHashTime += d.Nanoseconds()
				totalHashes++
				break
			case "Uploading:":
				d, err := time.ParseDuration(tokens[2])
				if err != nil {
					log.Fatalf("Error parsing uploading duration: %v, %v", tokens[2], err)
				}
				cumulativeUploadTime += d.Nanoseconds()
				break
			case "Finding":
				d, err := time.ParseDuration(tokens[4])
				if err != nil {
					log.Fatalf("Error parsing finding duration: %v, %v", tokens[4], err)
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
	log.Printf("Runs: %v, Hashes: %v", totalRuns, totalHashes)
	hashTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeHashTime)/float64(totalHashes)))
	if err != nil {
		log.Fatalf("Error parsing hash duration: %v, %v", float64(cumulativeHashTime)/float64(totalHashes), err)
	}
	totalHashTime, err := time.ParseDuration(fmt.Sprintf("%dns", cumulativeHashTime))
	if err != nil {
		log.Fatalf("Error parsing total hash duration: %v", err)
	}
	uploadTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeUploadTime)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing upload duration: %v, %v", float64(totalRuns), err)
	}
	findMissingTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeFindMissingTime)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing find missing duration: %v, %v", float64(cumulativeFindMissingTime)/float64(totalRuns), err)
	}
	executeTime, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeExecuteTime)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing execute duration: %v, %v", float64(cumulativeExecuteTime)/float64(totalRuns), err)
	}
	uploadTimeServer, err := time.ParseDuration(fmt.Sprintf("%.0fns", float64(cumulativeUploadTimeServer)/float64(totalRuns)))
	if err != nil {
		log.Fatalf("Error parsing server upload duration: %v, %v", float64(cumulativeUploadTimeServer)/float64(totalRuns), err)
	}
	fmt.Printf("Client:\nAverage Hash time: %v\nTotal Hash Time: %v\nUpload Time: %v\nFind Missing Time: %v\n", hashTime, totalHashTime, uploadTime, findMissingTime)
	fmt.Printf("Bytes Hashed: %v\nBytes Uploaded: %v\n", totalHashed, totalUploaded)
	var hashRate string
	if rateOfHashing := float64(totalHashed) / (float64(cumulativeHashTime) / (1000 * 1000 * 1000)); rateOfHashing < 1024 {
		hashRate = fmt.Sprintf("%fB/s", rateOfHashing)
	} else if rateOfHashing < 1024*1024 {
		hashRate = fmt.Sprintf("%fkB/s", rateOfHashing/1024)
	} else if rateOfHashing < 1024*1024*1024 {
		hashRate = fmt.Sprintf("%fMB/s", rateOfHashing/(1024*1024))
	} else {
		hashRate = fmt.Sprintf("%fGB/s", rateOfHashing/(1024*1024*1024))
	}
	fmt.Printf("Rate of Hashing: %s\n", hashRate)
	fmt.Printf("Server:\nExecute time: %v\nUpload Time: %v\n", executeTime, uploadTimeServer)
}
