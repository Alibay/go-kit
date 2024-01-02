package grpc

import (
	"math"
	"testing"

	"github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	kit.Suite
}

func (s *testSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) Test() {
	accessToken, err := kit.GenerateInternalAccessToken(
		s.Ctx,
		[]byte("YunduPOjwY28mJMaQDJi371IUBSaGSKw"),
		math.MaxInt,
		"test",
	)
	s.NoError(err)
	s.NotEmpty(accessToken)
	s.L().DbgF("token: %s", accessToken)
}
