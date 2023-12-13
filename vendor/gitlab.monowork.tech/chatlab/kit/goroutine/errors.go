package goroutine

import (
	"context"
	"gitlab.monowork.tech/chatlab/kit"
)

const (
	ErrCodeGoroutineNoLogger = "GORTN-001"
)

var (
	ErrGoroutineNoLogger = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeGoroutineNoLogger, "either logger or logger func must be specified").C(ctx).Err()
	}
)
