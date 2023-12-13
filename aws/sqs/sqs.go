package sqs

import (
	"context"

	"github.com/Alibay/go-kit/logger"

	kit "github.com/Alibay/go-kit"
	kitAws "github.com/Alibay/go-kit/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	ErrCodeSQSGetUrl        = "SQS-001"
	ErrCodeSQSGetMessages   = "SQS-002"
	ErrCodeSQSSubGetMessage = "SQS-003"
)

var (
	ErrSQSGetUrl = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeSQSGetUrl, "get url").C(ctx).Wrap(cause).Err()
	}
	ErrSQSGetMessages = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeSQSGetMessages, "get messages").C(ctx).Wrap(cause).Err()
	}
	ErrSQSSubGetMessage = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeSQSSubGetMessage, "subscriber: get message").C(ctx).Wrap(cause).Err()
	}
)

type Config struct {
	FetchInterval     int64 `mapstructure:"fetch_interval"`
	MaxMessages       int32 `mapstructure:"max_messages"`
	VisibilityTimeout int32 `mapstructure:"visibility_timeout"`
}

type Client struct {
	log       logger.CLoggerFunc
	awsCfg    *kitAws.Config
	sqsClient *sqs.Client
}

func NewClient(awsCfg *kitAws.Config, log logger.CLoggerFunc) *Client {
	return &Client{
		log:    log,
		awsCfg: awsCfg,
	}
}

func (c *Client) Init(ctx context.Context) error {
	awsConfig, err := kitAws.GetAwsConfig(ctx, c.awsCfg)
	if err != nil {
		return err
	}
	c.sqsClient = sqs.NewFromConfig(*awsConfig)
	return nil
}

func (c *Client) GetQueueURL(ctx context.Context, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	r, err := c.sqsClient.GetQueueUrl(ctx, input)
	if err != nil {
		return nil, ErrSQSGetUrl(ctx, err)
	}
	return r, nil
}

func (c *Client) GetMessages(ctx context.Context, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	r, err := c.sqsClient.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, ErrSQSGetMessages(ctx, err)
	}
	return r, nil
}
