package redis

import (
	"context"
	"net"
	"time"

	"github.com/Alibay/go-kit/logger"
	redis "github.com/go-redis/redis/v8"
)

const (
	NotFound = redis.Nil
)

type Redis struct {
	Instance *redis.Client
	Ttl      time.Duration
	log      logger.CLoggerFunc
}

// Config redis config
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Db       int
	Ttl      uint
}

func (r *Redis) l() logger.CLogger {
	return r.log().Cmp("redis")
}

func Open(ctx context.Context, params *Config, log logger.CLoggerFunc) (*Redis, error) {

	l := log().Cmp("redis").Mth("open")

	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(params.Host, params.Port),
		Username: params.Username,
		Password: params.Password,
		DB:       params.Db,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, ErrRedisPingErr(err)
	}

	l.Inf("ok")
	return &Redis{
		Instance: client,
		Ttl:      time.Duration(params.Ttl) * time.Second,
		log:      log,
	}, nil
}

func (r *Redis) Close() {
	l := r.l().Mth("close")
	if r.Instance != nil {
		_ = r.Instance.Close()
	}
	l.Inf("ok")
}
