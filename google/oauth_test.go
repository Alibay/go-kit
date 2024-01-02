//go:build example

package google

import (
	"encoding/json"
	"testing"

	"github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
)

type oauthTestSuite struct {
	kit.Suite
}

func (s *oauthTestSuite) SetupSuite() {
	s.Suite.Init(nil)
}

func TestOAuthSuite(t *testing.T) {
	suite.Run(t, new(oauthTestSuite))
}

type cred struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
}

type googleCfg struct {
	Web       *cred `json:"web"`
	Installed *cred `json:"installed"`
}

// Use https://developers.google.com/oauthplayground/ to obtain clientId and token for test
func (s *oauthTestSuite) Test() {
	cfg := &googleCfg{
		Web: &cred{
			ClientID:     "407408718192.apps.googleusercontent.com",
			ClientSecret: "",
			RedirectURIs: []string{"https://google.com"},
			AuthURI:      "",
			TokenURI:     "",
		},
	}
	v, _ := json.Marshal(cfg)
	o := NewOAuth(&Config{
		JsonConfiguration: string(v),
	}, s.L())
	ui, err := o.GetGoogleUser(s.Ctx, "<put-our-token-here>")
	s.NoError(err)
	s.NotEmpty(ui)
}
