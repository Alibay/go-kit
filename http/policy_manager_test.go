package http

import (
	"context"
	"net/http"
	"testing"

	kit "github.com/Alibay/go-kit"

	kitTesting "github.com/Alibay/go-kit/testing"
	"github.com/stretchr/testify/suite"
)

type policyManagerTestSuite struct {
	kitTesting.Suite
	policyManager ResourcePolicyManager
}

func (s *policyManagerTestSuite) SetupSuite() {
	s.Suite.Init(logf)
	s.policyManager = NewResourcePolicyManager()
}

func (s *policyManagerTestSuite) SetupTest() {
}

func TestPolicyManagerSuite(t *testing.T) {
	suite.Run(t, new(policyManagerTestSuite))
}

func (s *policyManagerTestSuite) Test_WhenSingleResourceWithPositiveCondition() {
	routeId := kit.NewRandString()
	resourceCode := "resource"
	resource := Resource(resourceCode, "rwxd").When(func(context.Context, *http.Request) (bool, error) { return true, nil })

	s.policyManager.RegisterResourceMapping(routeId, resource)

	authResources, err := s.policyManager.GetRequestedResources(s.Ctx, routeId, nil)
	s.NoError(err)
	s.NotEmpty(authResources)
	s.Equal(1, len(authResources))
	s.Equal(resourceCode, authResources[0].Resource)
	s.Equal(4, len(authResources[0].Permissions))
}

func (s *policyManagerTestSuite) Test_WhenMultipleResources_PositiveAndNegativeConditions() {
	routeId := kit.NewRandString()
	resourceCode1 := "resource1"
	resourceCode2 := "resource2"
	resource1 := Resource(resourceCode1, "rwxd").When(func(context.Context, *http.Request) (bool, error) { return true, nil })
	resource2 := Resource(resourceCode2, "rwxd").WhenNot(func(context.Context, *http.Request) (bool, error) { return true, nil })

	s.policyManager.RegisterResourceMapping(routeId, resource1, resource2)

	authResources, err := s.policyManager.GetRequestedResources(s.Ctx, routeId, nil)
	s.NoError(err)
	s.NotEmpty(authResources)
	s.Equal(1, len(authResources))
	s.Equal(resourceCode1, authResources[0].Resource)
	s.Equal(4, len(authResources[0].Permissions))
}

func (s *policyManagerTestSuite) Test_WhenMultipleResources_NoConditions() {
	routeId := kit.NewRandString()
	resourceCode1 := "resource1"
	resourceCode2 := "resource2"
	resource1 := Resource(resourceCode1, "r")
	resource2 := Resource(resourceCode2, "w")

	s.policyManager.RegisterResourceMapping(routeId, resource1, resource2)

	authResources, err := s.policyManager.GetRequestedResources(s.Ctx, routeId, nil)
	s.NoError(err)
	s.NotEmpty(authResources)
	s.Equal(2, len(authResources))
}

func (s *policyManagerTestSuite) Test_WithoutResources() {
	routeId := kit.NewRandString()

	s.policyManager.RegisterResourceMapping(routeId)

	authResources, err := s.policyManager.GetRequestedResources(s.Ctx, routeId, nil)
	s.NoError(err)
	s.Empty(authResources)
}
