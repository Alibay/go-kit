package centrifugo

import (
	"context"
	"net"
	"time"

	"github.com/Alibay/go-kit"
	apiproto "github.com/Alibay/go-kit/centrifugo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	// Connect connects to grpc server
	Connect(ctx context.Context) error
	// Close closes connection
	Close(ctx context.Context) error
	// Publish publishes message to channel
	Publish(ctx context.Context, channel string, msg any) error
	// GetPresence retrieves presence by channel
	GetPresence(ctx context.Context, channel string) (*apiproto.PresenceResult, error)
}

type ServerConfig struct {
	Host   string // Host grpc host
	Port   string // Port grpc port
	ApiKey string `mapstructure:"api_key"` // ApiKey to connect to grpc server
	Secret string // Secret used for token generating
}

func NewServer(cfg *ServerConfig, logger kit.CLoggerFunc) Server {
	return &serverImpl{
		logger: logger,
		cfg:    cfg,
	}
}

type serverImpl struct {
	cfg    *ServerConfig
	logger kit.CLoggerFunc
	conn   *grpc.ClientConn
	client apiproto.CentrifugoApiClient
}

func (s *serverImpl) l() kit.CLogger {
	return s.logger().Cmp("centrifugo-server")
}

type keyAuth struct {
	key string
}

func (t keyAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "apikey " + t.key,
	}, nil
}

func (t keyAuth) RequireTransportSecurity() bool {
	return false
}

func (s *serverImpl) Connect(ctx context.Context) error {
	l := s.l().C(ctx).F(kit.KV{"host": s.cfg.Host, "port": s.cfg.Port})

	// configure grpc connection
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// if api key specified, passed it
	if s.cfg.ApiKey != "" {
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(keyAuth{s.cfg.ApiKey}))
	}

	// connect
	var err error
	s.conn, err = grpc.Dial(net.JoinHostPort(s.cfg.Host, s.cfg.Port), dialOptions...)
	if err != nil {
		return ErrGrpcServerConnect(ctx, err)
	}

	s.client = apiproto.NewCentrifugoApiClient(s.conn)

	l.Dbg("ok")

	return nil

}

func (s *serverImpl) Close(ctx context.Context) error {
	if s.client != nil {
		_ = s.conn.Close()
		s.client = nil
	}
	return nil
}

func (s *serverImpl) Publish(ctx context.Context, channel string, payload any) error {
	s.l().C(ctx).Mth("publish").F(kit.KV{"chan": channel}).Dbg().TrcObj("%v", payload)

	bytes, _ := kit.JsonEncode(payload)

	rs, err := s.client.Publish(ctx, &apiproto.PublishRequest{
		Channel: channel,
		Data:    bytes,
	})
	if err != nil {
		return ErrGrpcServerPublish(ctx, err)
	}
	if rs.Error != nil {
		return ErrGrpcServerPublishRs(ctx, rs.Error)
	}

	return nil
}

func (s *serverImpl) GetPresence(ctx context.Context, channel string) (*apiproto.PresenceResult, error) {
	s.l().C(ctx).Mth("get-presence").F(kit.KV{"chan": channel}).Dbg()

	rs, err := s.client.Presence(ctx, &apiproto.PresenceRequest{
		Channel: channel,
	})
	if err != nil {
		return nil, ErrGrpcServerPresence(ctx, err)
	}
	if rs.Error != nil {
		return nil, ErrGrpcServerPresenceRs(ctx, rs.Error)
	}

	return rs.Result, nil
}

// GenerateSubscribeToken generates a client subscribe token
func GenerateSubscribeToken(ctx context.Context, secret, userId, channel string, ttl time.Duration) (string, error) {
	return kit.GenJwtToken(ctx, &kit.JwtRequest{
		UserId:   userId,
		Secret:   []byte(secret),
		ExpireAt: time.Now().Add(ttl),
		Claims: map[string]any{
			"expired_at": time.Now().Add(ttl).Unix(),
			"created_at": time.Now().Unix(),
			"channel":    channel,
		},
	})
}

// GenerateConnectToken generates a client connect token
func GenerateConnectToken(ctx context.Context, secret, userId string, autoSubscribeChannels []string, ttl time.Duration, info any) (string, error) {

	claims := map[string]any{
		"expired_at": time.Now().Add(ttl).Unix(),
		"created_at": time.Now().Unix(),
	}

	if len(autoSubscribeChannels) > 0 {
		claims["channels"] = autoSubscribeChannels
	}

	if info != nil {
		claims["info"] = info
	}

	return kit.GenJwtToken(ctx, &kit.JwtRequest{
		UserId:   userId,
		Secret:   []byte(secret),
		ExpireAt: time.Now().Add(ttl),
		Claims:   claims,
	})
}
