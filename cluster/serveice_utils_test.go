package cluster

import (
	"testing"

	kitTesting "github.com/Alibay/go-kit/testing"
	"github.com/stretchr/testify/suite"
)

type serviceUtilsTestSuite struct {
	kitTesting.Suite
}

func (s *serviceUtilsTestSuite) SetupSuite() {
	s.Suite.Init(nil)
}

func TestServiceUtilsSuite(t *testing.T) {
	suite.Run(t, new(serviceUtilsTestSuite))
}

func (s *serviceUtilsTestSuite) Test_GetServiceRootPath_WhenEmptyInput() {
	s.Empty(GetServiceRootPath(""))
}

func (s *serviceUtilsTestSuite) Test_GetServiceRootPath_WhenNotExistent() {
	s.Empty(GetServiceRootPath("some"))
}

func (s *serviceUtilsTestSuite) Test_GetServiceRootPath_WhenKit() {
	s.NotEmpty(GetServiceRootPath("kit"))
}
