package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"poc-go-kinesis/lib/log"
	"poc-go-kinesis/lib/stream"
	"time"
)

const (
	defaultStreamName                      = ""
	defaultReaderName                      = "default-reader"
	defaultReadBatchSize                   = 10
	defaultChannelBufferSize               = 100
)

func main() {
	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		if VersionPrerelease != "" {
			fmt.Println("Version PreRelease:", VersionPrerelease)
		}
		return
	}

	awsConfig := &aws.Config{
		Region:      aws.String("ap-southeast-2"),
		Credentials: credentials.NewEnvCredentials(),
	}

	processor, err := stream.NewProcessor(awsConfig, "blue-content", 1000 * time.Millisecond, 5, 5)

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Start processing")

	for record := range processor.GetRecords() {
		log.Print("Data: ", string(record.Data))
	}
}
