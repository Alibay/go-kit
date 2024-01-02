package google

import (
	"context"

	"github.com/Alibay/go-kit"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"

	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
	oauthv2 "google.golang.org/api/oauth2/v2"
)

const (
	ErrCodeOAuthConfigRead = "OAUTH-001"
	ErrCodeOAuthService    = "OAUTH-002"
	ErrCodeOAuthGetUser    = "OAUTH-003"
)

var (
	ErrOAuthConfigRead = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeOAuthConfigRead, "reading config").Wrap(cause).C(ctx).Business().Err()
	}
	ErrOAuthService = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeOAuthService, "service").Wrap(cause).C(ctx).Business().Err()
	}
	ErrOAuthGetUser = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeOAuthGetUser, "get user").Wrap(cause).C(ctx).Business().Err()
	}
)

type OAuth2 interface {
	// GetGoogleUser retrieves google user info
	GetGoogleUser(ctx context.Context, token string) (*oauthv2.Userinfo, error)
}

type oauth struct {
	cfg       *Config
	client    *http.Client
	logger    kit.CLogger
	lazy      sync.Once
	googleCfg *oauth2.Config
}

func NewOAuth(cfg *Config, logger kit.CLogger) OAuth2 {
	return &oauth{
		cfg:    cfg,
		logger: logger,
		client: &http.Client{Timeout: time.Duration(cfg.ClientTimeout)},
	}
}

func (o *oauth) l() kit.CLogger {
	return o.logger.Cmp("oauth")
}

func (o *oauth) GetGoogleUser(ctx context.Context, token string) (*oauthv2.Userinfo, error) {
	o.l().Mth("get-user").C(ctx).Dbg()

	// load oauth config
	var err error
	o.lazy.Do(func() {
		o.googleCfg, err = google.ConfigFromJSON([]byte(o.cfg.JsonConfiguration))
	})
	if err != nil {
		return nil, ErrOAuthConfigRead(ctx, err)
	}

	// client with token
	at := &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}

	httpClient := o.googleCfg.Client(ctx, at)
	httpClient.Timeout = time.Duration(o.cfg.ClientTimeout)

	// prepare google service
	gService, err := oauthv2.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, ErrOAuthService(ctx, err)
	}

	// execute
	ui, err := gService.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return nil, ErrOAuthGetUser(ctx, err)
	}

	return ui, nil
}
