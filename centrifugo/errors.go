package centrifugo

import (
	"context"

	"github.com/Alibay/go-kit"
	apiproto "github.com/Alibay/go-kit/centrifugo/proto"
)

const (
	ErrCodeCentrifugoConnect     = "CTRF-001"
	ErrCodeCentrifugoPublish     = "CTRF-002"
	ErrCodeCentrifugeInternal    = "CTRF-004"
	ErrCodeCentrifugoSubscribing = "CTRF-005"
	ErrCodeCentrifugoSubscribe   = "CTRF-006"
	ErrCodeGrpcServerConnect     = "CTRF-007"
	ErrCodeGrpcServerPublish     = "CTRF-008"
	ErrCodeGrpcServerPublishRs   = "CTRF-009"
)

var (
	ErrCentrifugoConnect = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCentrifugoConnect, "centrifugo: connect").Wrap(cause).Err()
	}
	ErrCentrifugoGrpcPublish = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCentrifugoPublish, "centrifugo: publish").Wrap(cause).Err()
	}
	ErrCentrifugeInternal = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCentrifugeInternal, "centrifugo error").Wrap(cause).Err()
	}
	ErrCentrifugoSubscribing = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCentrifugoSubscribing, "centrifugo: subscribing").Wrap(cause).Err()
	}
	ErrCentrifugoSubscribe = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCentrifugoSubscribe, "centrifugo: subscribe").Wrap(cause).Err()
	}
	ErrGrpcServerConnect = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcServerConnect, "centrifugo server: connect").Wrap(cause).Err()
	}
	ErrGrpcServerPublish = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcServerPublish, "centrifugo server: publish").Wrap(cause).Err()
	}
	ErrGrpcServerPublishRs = func(ctx context.Context, cause *apiproto.Error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcServerPublishRs, "centrifugo server: publish %s (%d)", cause.Message, cause.Code).Err()
	}
	ErrGrpcServerPresence = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcServerPublish, "centrifugo server: presence").Wrap(cause).Err()
	}
	ErrGrpcServerPresenceRs = func(ctx context.Context, cause *apiproto.Error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcServerPublishRs, "centrifugo server: presence %s (%d)", cause.Message, cause.Code).Err()
	}
)
