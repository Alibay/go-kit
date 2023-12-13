//go:build dev

package sqs

import (
	"testing"

	"github.com/Alibay/go-kit/logger"

	kit "github.com/Alibay/go-kit"
	kitAws "github.com/Alibay/go-kit/aws"
	kitTesting "github.com/Alibay/go-kit/testing"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/suite"
)

type s3TestSuite struct {
	kitTesting.Suite
	log logger.CLoggerFunc
}

func (s *s3TestSuite) SetupSuite() {
	s.log = func() logger.CLogger { return logger.L(logger.InitLogger(&logger.LogConfig{Level: logger.TraceLevel})) }
	s.Suite.Init(s.log)
}

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(s3TestSuite))
}

var (
	awsCfg = &kitAws.Config{
		Region:              "eu-central-1",
		AccessKeyId:         "access_key_id",
		SecretAccessKey:     "secret_access_key",
		SharedConfigProfile: "chatlab/dev",
	}
)

func (s *s3TestSuite) Test_Init() {
	// init client
	client := NewClient(awsCfg, s.log)
	s.NoError(client.Init(s.Ctx))
	s.NotEmpty(client.sqsClient)

	_, err := client.GetQueueURL(s.Ctx, &sqs.GetQueueUrlInput{
		QueueName:              kit.StringPtr("ext-storage-dev"),
		QueueOwnerAWSAccountId: nil,
	})
	s.NoError(err)
}
