package profile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/Alibay/go-kit"
	"github.com/Alibay/go-kit/goroutine"
)

const (
	ErrCodeProfileHttpError = "PRF-001"
)

var (
	ErrProfileHttpError = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeProfileHttpError, "").Wrap(cause).Err()
	}
)

// Server exposes profile dumps
type Server interface {
	// Init initializes server with given opts
	Init(config *Config) error
	// Listen starts async listening
	Listen()
	// Close closes connection
	Close()
}

type Config struct {
	Enabled bool
	Port    string
}

type profileServerImpl struct {
	logger kit.CLoggerFunc
	config *Config
}

func New(logger kit.CLoggerFunc) Server {
	return &profileServerImpl{
		logger: logger,
	}
}

func (p *profileServerImpl) l() kit.CLogger {
	return p.logger().Cmp("profile")
}

func (p *profileServerImpl) Init(config *Config) error {
	p.l().Mth("init").Dbg()
	p.config = config
	return nil
}

func (p *profileServerImpl) Listen() {
	goroutine.New().
		WithLoggerFn(p.logger).
		WithRetry(goroutine.Unrestricted).
		Go(context.Background(),
			func() {
				l := p.l().Mth("listen").Inf("start listening")
				if err := http.ListenAndServe(fmt.Sprintf(":%s", p.config.Port), nil); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						l.E(ErrProfileHttpError(err)).St().Err()
					} else {
						l.Inf("closed")
					}
				}
			},
		)
}

func (p *profileServerImpl) Close() {
	p.l().Mth("close").Inf("closed")
}
