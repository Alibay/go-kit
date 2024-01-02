//go:build dev

package s3

import (
	"fmt"
	"testing"

	"github.com/Alibay/go-kit"
	kitAws "github.com/Alibay/go-kit/aws"
	"github.com/stretchr/testify/suite"
)

type s3TestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
}

func (s *s3TestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(s3TestSuite))
}

var (
	s3Cfg = &Config{
		PublicBucketName:  "ext.storage.dev.chatlab",
		PrivateBucketName: "int.storage.dev.chatlab",
		PresignedLinkTTL:  60,
	}
	awsCfg = &kitAws.Config{
		Region:              "eu-central-1",
		AccessKeyId:         "access_key_id",
		SecretAccessKey:     "secret_access_key",
		SharedConfigProfile: "chatlab/dev",
	}
)

func (s *s3TestSuite) Test_S3() {

	// init client
	client := NewClient(awsCfg, s3Cfg, s.logger)
	s.NoError(client.Init(s.Ctx))
	s.NotEmpty(client.s3Client)

	// get new upload link
	ownerId := kit.NewId()
	fn := fmt.Sprintf("%s.png", kit.NewRandString())
	url, key, err := client.GetNewFileUploadLink(s.Ctx, false, false, ownerId, fn, "test")
	s.NoError(err)
	s.NotEmpty(url)
	s.NotEmpty(key)

	// update
	url, err = client.GetUpdateFileUploadLink(s.Ctx, false, key)
	s.NoError(err)
	s.NotEmpty(url)

	// get
	url, err = client.GetGetFileLink(s.Ctx, false, key)
	s.NoError(err)
	s.NotEmpty(url)
	s.L().DbgF("url: %s", url)

	// delete
	s.NoError(client.DeleteFileByKey(s.Ctx, false, key))
}
