package transports

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type CloudWatchTransport struct {
	client        *cloudwatchlogs.Client
	logGroup      string
	logStream     string
	sequenceToken *string
}

func NewCloudWatchTransport(cfg aws.Config) *CloudWatchTransport {
	return &CloudWatchTransport{
		client:    cloudwatchlogs.NewFromConfig(cfg),
		logGroup:  "test-group",
		logStream: "test-stream",
	}
}

func (t *CloudWatchTransport) Write(p []byte) (n int, err error) {
	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(t.logGroup),
		LogStreamName: aws.String(t.logStream),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(string(p)),
				Timestamp: aws.Int64(time.Now().UnixMilli()),
			},
		},
		SequenceToken: t.sequenceToken,
	}

	fmt.Printf("cloudwatch: %v\n", input)
	resp, err := t.client.PutLogEvents(context.Background(), input)
	if err != nil {
		fmt.Printf("cloudwatch error: %v\n", err)
		return 0, err
	}

	t.sequenceToken = resp.NextSequenceToken
	return len(p), nil
}

func (t *CloudWatchTransport) Close() error {
	return nil
}
