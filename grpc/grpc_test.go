//go:build integration

package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
)

type grpcTestSuite struct {
	kit.Suite
}

func (s *grpcTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func TestGrpcSuite(t *testing.T) {
	suite.Run(t, new(grpcTestSuite))
}

type srvImpl struct {
	UnimplementedTestServiceServer
}

func (s *srvImpl) WithError(ctx context.Context, rq *WithErrorRequest) (*WithErrorResponse, error) {
	e := kit.NewAppErrBuilder("TST-123", "%s happens", "shit").GrpcSt(uint32(codes.AlreadyExists)).C(ctx).F(kit.KV{"id": "123"}).Err()
	return nil, e
}

func (s *srvImpl) WithPanic(ctx context.Context, rq *WithPanicRequest) (*WithPanicResponse, error) {
	panic("JUST A PANIC, BRO.....")
}

func (s *srvImpl) Do(ctx context.Context, rq *Empty) (*Empty, error) {
	logger := func() kit.CLogger {
		return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel, Context: true, Format: "json"}))
	}
	logger().C(ctx).Trc("log")
	return &Empty{}, nil
}

func (s *grpcTestSuite) Test_WhenAppError() {

	srv, _ := NewServer("test", s.L, &ServerConfig{Port: "55556"})
	defer srv.Close()
	RegisterTestServiceServer(srv.Srv, &srvImpl{})

	go func() {
		if err := srv.Listen(s.Ctx); err != nil {
			return
		}
	}()

	time.Sleep(time.Millisecond * 200)

	cl, err := NewClient(&ClientConfig{Host: "localhost", Port: "55556"})
	s.NoError(err)
	defer func() { _ = cl.Conn.Close() }()

	ctx := kit.NewRequestCtx().WithApp("test").ToContext(context.Background())
	svc := NewTestServiceClient(cl.Conn)
	_, err = svc.WithError(ctx, &WithErrorRequest{})
	if err != nil {
		if appErr, ok := kit.IsAppErr(err); ok {
			ctx := appErr.Fields()["ctx"].(map[string]interface{})
			s.Equal(ctx["_ctx.app"], "test")
			s.Equal(appErr.Fields()["id"], "123")
			s.L().E(err).Err()
		} else {
			s.Fatal("not app error")
		}
	}
}

func (s *grpcTestSuite) Test_WhenPanicRecover() {
	port := "55557"
	srv, _ := NewServer("test", s.L, &ServerConfig{Port: port})
	defer srv.Close()
	RegisterTestServiceServer(srv.Srv, &srvImpl{})

	go func() {
		if err := srv.Listen(s.Ctx); err != nil {
			return
		}
	}()

	time.Sleep(time.Millisecond * 200)

	cl, err := NewClient(&ClientConfig{Host: "localhost", Port: port})
	s.NoError(err)
	defer cl.Conn.Close()

	ctx := kit.NewRequestCtx().ToContext(context.Background())
	svc := NewTestServiceClient(cl.Conn)
	_, err = svc.WithPanic(ctx, &WithPanicRequest{})
	s.AssertAppErr(err, kit.ErrCodePanic)
}

func (s *grpcTestSuite) Test_WithAuth() {
	port := "55557"
	secret := "secret"
	srv, _ := NewServer("test", s.L, &ServerConfig{
		Port: port,
		Auth: ServerAuthConfig{
			Enabled: true,
			Secret:  secret,
		},
	})
	defer srv.Close()
	RegisterTestServiceServer(srv.Srv, &srvImpl{})

	go func() {
		if err := srv.Listen(s.Ctx); err != nil {
			return
		}
	}()

	time.Sleep(time.Millisecond * 200)

	cl, err := NewClient(&ClientConfig{
		Host: "localhost",
		Port: port,
		Auth: ClientAuthConfig{
			Enabled:     true,
			TokenSecret: secret,
			TokenTTL:    100,
			Caller:      "test",
		},
	})
	s.NoError(err)
	defer cl.Conn.Close()

	ctx := kit.NewRequestCtx().
		WithNewRequestId().
		WithClientIp("123.123.123.123").
		WithUser(kit.NewId(), "user").ToContext(context.Background())
	svc := NewTestServiceClient(cl.Conn)
	rs, err := svc.Do(ctx, &Empty{})
	s.NoError(err)
	s.NotEmpty(rs)
}

func (s *grpcTestSuite) Test_WithoutAuth() {
	port := "55557"
	srv, _ := NewServer("test", s.L, &ServerConfig{
		Port: port,
		Auth: ServerAuthConfig{
			Enabled: false,
		},
	})
	defer srv.Close()
	RegisterTestServiceServer(srv.Srv, &srvImpl{})

	go func() {
		if err := srv.Listen(s.Ctx); err != nil {
			return
		}
	}()

	time.Sleep(time.Millisecond * 200)

	cl, err := NewClient(&ClientConfig{
		Host: "localhost",
		Port: port,
		Auth: ClientAuthConfig{
			Enabled: false,
		},
	})
	s.NoError(err)
	defer cl.Conn.Close()

	ctx := kit.NewRequestCtx().ToContext(context.Background())
	svc := NewTestServiceClient(cl.Conn)
	rs, err := svc.Do(ctx, &Empty{})
	s.NoError(err)
	s.NotEmpty(rs)
}
