package stream

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"poc-go-kinesis/lib/log"
)

type Processor struct {
	clientAPI kinesisiface.KinesisAPI

	streamName         string
	streamReadInterval time.Duration
	readBatchSize      int
	channelBufferSize  int

	err         error
	recordsChan chan *kinesis.Record

	consumers []*ShardReader
}

func NewProcessor(awsConfig *aws.Config, streamName string, streamReadInterval time.Duration, readBatchSize int, channelBufferSize int) (*Processor, error) {
	r := &Processor{
		clientAPI:          kinesis.New(session.New(awsConfig)),
		streamName:         streamName,
		streamReadInterval: streamReadInterval,
		readBatchSize:      readBatchSize,
		channelBufferSize:  channelBufferSize,

		recordsChan: make(chan *kinesis.Record),
		consumers:   []*ShardReader{},
	}

	return r, nil
}

func (sp *Processor) GetRecords() chan *kinesis.Record {
	go sp.consumeRecords()
	return sp.recordsChan
}

func (sp *Processor) consumeRecords() {
	desc, err := sp.GetStreamDescription(sp.streamName)
	if err != nil {
		sp.err = err
		log.Fatal(err)
		return
	}

	for _, shard := range desc.Shards {
		shardReader, err := NewShardReader(sp.clientAPI, sp.streamName, *shard.ShardId, sp.streamReadInterval, sp.readBatchSize, sp.channelBufferSize)
		if err != nil {
			sp.err = err
			log.Fatal(err)
			return
		}

		go sp.consumeShard(shardReader)
	}
}

func (sp *Processor) GetStreamDescription(streamName string) (*kinesis.StreamDescription, error) {
	out, err := sp.clientAPI.DescribeStream(&kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	})
	if err != nil {
		return nil, err
	}

	return out.StreamDescription, nil
}

func (sp *Processor) consumeShard(sr *ShardReader) {
	for record := range sr.GetRecords() {
		sp.recordsChan <- record
	}
}
