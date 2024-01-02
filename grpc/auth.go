package grpc

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/Alibay/go-kit"
	"google.golang.org/grpc/metadata"
)

const (
	Headers    = "headers"
	AuthHeader = "access-token"
)

type authContextValue struct {
	CallerName  string `json:"caller_name"`
	AccessToken string `json:"access_token"`
}

type internalAuthConfig struct {
	InternalAccessTokenSecret string `mapstructure:"internalAccessTokenSecret"`
	InternalAccessTokenTTL    int    `mapstructure:"internalAccessTokenTTL"`
}

func buildAuthFunc(internalAccessTokenSecret []byte) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {

		// check metadata passed
		md, isMdExists := metadata.FromIncomingContext(ctx)
		if !isMdExists {
			return nil, ErrGrpcAuthNoMd(ctx)
		}

		// check header exists
		authHeaders := md.Get(AuthHeader)
		if len(authHeaders) == 0 {
			return nil, ErrGrpcAuthNoHeader(ctx)
		}

		callerName, err := kit.ParseInternalAccessToken(ctx, internalAccessTokenSecret, authHeaders[0])
		if err != nil {
			return nil, ErrGrpcAuthParseToken(ctx, err)
		}

		// return modified context
		return context.WithValue(ctx, Headers, authContextValue{
			CallerName:  callerName,
			AccessToken: authHeaders[0],
		}), nil

	}
}

func buildAccessToken(ctx context.Context, cfg *internalAuthConfig, serviceName string) (string, error) {
	accessToken, err := kit.GenerateInternalAccessToken(
		ctx,
		[]byte(cfg.InternalAccessTokenSecret),
		cfg.InternalAccessTokenTTL,
		serviceName,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
