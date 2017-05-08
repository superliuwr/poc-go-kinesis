package stream

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type ShardReader struct {
	clientAPI         kinesisiface.KinesisAPI
	streamName        string
	shardId           string
	readInterval      time.Duration
	batchSize         *int64
	channelBufferSize int
	err               error
}

// NewShardReader defines a ShardReader.
func NewShardReader(kinesisAPI kinesisiface.KinesisAPI, streamName string, shardId string, streamReadInterval time.Duration, readBatchSize int, channelBufferSize int) (*ShardReader, error) {
	sr := &ShardReader{
		clientAPI:         kinesisAPI,
		streamName:        streamName,
		shardId:           shardId,
		readInterval:      streamReadInterval,
		batchSize:         aws.Int64(int64(readBatchSize)),
		channelBufferSize: channelBufferSize,
	}

	return sr, nil
}

// Records consumes a shard from the last untrimmed record. It returns a read-only buffered channel via which results
// are delivered.
func (sr *ShardReader) GetRecords() <-chan *kinesis.Record {
	ch := make(chan *kinesis.Record, sr.channelBufferSize)

	shardIteratorType := aws.String(kinesis.ShardIteratorTypeTrimHorizon)

	iteratorInput := &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(sr.streamName),
		ShardId:           aws.String(sr.shardId),
		ShardIteratorType: shardIteratorType,
	}

	iterator, err := sr.clientAPI.GetShardIterator(iteratorInput)
	if err != nil {
		sr.err = err
		close(ch)
		return ch
	}

	go sr.consumeStream(ch, iterator.ShardIterator)

	return ch
}

func (sr *ShardReader) consumeStream(ch chan *kinesis.Record, shardIterator *string) {
	for {
		out, err := sr.clientAPI.GetRecords(&kinesis.GetRecordsInput{
			Limit:         sr.batchSize,
			ShardIterator: shardIterator,
		})

		if err != nil {
			sr.err = err
			close(ch)
			return
		}

		shardIterator = out.NextShardIterator

		if len(out.Records) == 0 {
			continue
		}

		for _, record := range out.Records {
			ch <- record
		}

		time.Sleep(sr.readInterval)
	}
}
