package grpc

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/Alibay/go-kit"
	"github.com/Alibay/go-kit/goroutine"
	"github.com/Alibay/go-kit/monitoring"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

// ServerAuthConfig authorization config
type ServerAuthConfig struct {
	Enabled bool   // Enabled if true auth is applied
	Secret  string // Secret key
}

// ServerConfig represents gRPC server configuration
type ServerConfig struct {
	Host  string
	Port  string
	Trace bool
	Auth  ServerAuthConfig
}

type Server struct {
	healthpb.HealthServer
	monitoring.MetricsProvider
	Srv     *grpc.Server
	Service string
	logger  kit.CLoggerFunc
	config  *ServerConfig
	ln      net.Listener
}

func NewServer(service string, logger kit.CLoggerFunc, config *ServerConfig) (*Server, error) {

	s := &Server{
		Service:      service,
		HealthServer: NewHealthServer(),
		logger:       logger,
		config:       config,
	}

	// setup server middleware
	opts := []grpc_recovery.Option{grpc_recovery.WithRecoveryHandlerContext(s.PanicFunc)}
	unaryInterceptors := []grpc.UnaryServerInterceptor{s.unaryServerInterceptor(), grpc_recovery.UnaryServerInterceptor(opts...)}
	streamInterceptors := []grpc.StreamServerInterceptor{s.streamServerInterceptor(), grpc_recovery.StreamServerInterceptor(opts...)}

	// authorization
	if s.config.Auth.Enabled {
		authFunc := buildAuthFunc([]byte(config.Auth.Secret))
		unaryInterceptors = append(unaryInterceptors, grpc_auth.UnaryServerInterceptor(authFunc))
		streamInterceptors = append(streamInterceptors, grpc_auth.StreamServerInterceptor(authFunc))
	}

	// build a new server
	s.Srv = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(unaryInterceptors...),
		grpc_middleware.WithStreamServerChain(streamInterceptors...),
	)

	// register health server
	healthpb.RegisterHealthServer(s.Srv, s)

	return s, nil
}

func (s *Server) Listen(ctx context.Context) error {
	l := s.logger().Cmp(s.Service).Pr("grpc").Mth("listen").F(kit.KV{"port": s.config.Port}).Inf("start listening")

	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// Enable SO_REUSEADDR
				err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				if err != nil {
					l.E(err).St().Err("could not set SO_REUSEADDR socket option")
				}
				// Enable SO_REUSEPORT
				err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				if err != nil {
					l.E(err).St().Err("Could not set SO_REUSEPORT socket option")
				}
			})
		},
	}

	// Start Listener
	var err error
	s.ln, err = lc.Listen(ctx, "tcp", fmt.Sprint(":", s.config.Port))
	if err != nil {
		return ErrGrpcSrvListen(err)
	}

	err = s.Srv.Serve(s.ln)
	if err != nil {
		return ErrGrpcSrvServe(err)
	}

	return nil

}

func (s *Server) ListenAsync(ctx context.Context) {
	goroutine.New().
		WithLoggerFn(s.logger).
		Cmp("grpc").
		Mth("listen").
		WithRetry(goroutine.Unrestricted).
		Go(context.Background(), func() {
		start:
			err := s.Listen(ctx)
			if err != nil {
				s.logger().E(err).St().Err()
				time.Sleep(time.Second * 5)
				goto start
			}
		})
}

func (s *Server) Close() {
	s.Srv.Stop()
	_ = s.ln.Close()
}

// this middleware is applied on server side
// it retrieves gRPC metadata and puts it to the context
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// convert metadata to request context
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = kit.FromGrpcMD(ctx, md)
		}
		resp, err := handler(ctx, req)

		// tracing
		if s.config.Trace {
			s.logger().C(ctx).Pr("grpc").
				Cmp(s.Service).
				Mth(info.FullMethod).
				C(ctx).
				TrcObj("rq: %+v, rs: %+v", req, resp)
		}

		// logging
		if err != nil {
			// log errors
			s.logger().C(ctx).Pr("grpc").Cmp(s.Service).Mth(info.FullMethod).E(err).St().Err()
		}

		// convert to grpc status
		if err != nil {
			err = toGrpcStatus(err)
		}

		return resp, err
	}
}

// this middleware is applied on server side
// it retrieves gRPC metadata and puts it to the context
func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		err := handler(srv, ss)

		// logging errors
		if err != nil {
			s.logger().Pr("grpc").Cmp(s.Service).Mth(info.FullMethod).E(err).St().Err()
		}

		// convert to grpc status
		if err != nil {
			err = toGrpcStatus(err)
		}

		return err
	}
}

func (s *Server) PanicFunc(ctx context.Context, panicCause interface{}) error {
	// convert metadata to request context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ctx = kit.FromGrpcMD(ctx, md)
	} else {
		ctx = context.Background()
	}
	err := kit.ErrPanic(ctx, panicCause)
	// log panic
	s.logger().Pr("grpc").Cmp(s.Service).E(err).St().Err()
	return err
}
