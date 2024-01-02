package sqs

import (
	"context"
	"time"

	"github.com/Alibay/go-kit"
	"github.com/Alibay/go-kit/goroutine"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Subscriber struct {
	logger    kit.CLoggerFunc
	client    *Client
	queueName string
	config    *Config
	receiver  chan types.Message
}

func NewSubscriber(client *Client, cfg *Config, queueName string, receiver chan types.Message, logger kit.CLoggerFunc) *Subscriber {
	return &Subscriber{
		logger:    logger,
		client:    client,
		queueName: queueName,
		config:    cfg,
		receiver:  receiver,
	}
}

func (s *Subscriber) l() kit.CLogger {
	return s.logger().Cmp("sqs-sub")
}

func (s *Subscriber) Run(ctx context.Context) error {
	l := s.l().C(ctx).Mth("run").Dbg()

	queueURL, err := s.client.GetQueueURL(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(s.queueName),
	})
	if err != nil {
		return err
	}

	goroutine.New().WithLogger(l).Go(ctx, func() {

		ticker := time.NewTicker(time.Duration(s.config.FetchInterval))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				messages, err := s.client.GetMessages(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            queueURL.QueueUrl,
					VisibilityTimeout:   s.config.VisibilityTimeout,
					MaxNumberOfMessages: s.config.MaxMessages,
				})
				if err != nil {
					s.l().C(ctx).Mth("run").E(ErrSQSSubGetMessage(ctx, err)).St().Err()
					continue
				}
				for _, message := range messages.Messages {
					s.receiver <- message
				}
			case <-ctx.Done():
				return
			}
		}
	})

	return nil

}
