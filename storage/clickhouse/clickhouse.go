package clickhouse

import (
	"database/sql"
	"fmt"

	"github.com/Alibay/go-kit/logger"

	kit "github.com/Alibay/go-kit"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouse struct {
	Instance clickhouse.Conn
	cfg      *Config
	log      logger.CLoggerFunc
}

// Config configuration parameters
type Config struct {
	User     string // User username
	Password string // Password password
	Database string // Database database name
	Port     string // Port connection
	Host     string // Host connection
	Debug    bool   // Debug if debug mode enabled
}

func Open(config *Config, log logger.CLoggerFunc) (*ClickHouse, error) {
	s := &ClickHouse{
		log: log,
		cfg: config,
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},
		Debug: config.Debug,
		Debugf: func(format string, v ...interface{}) {
			log().Cmp("click").Mth("debug").DbgF(format, v...)
		},
	})
	if err != nil {
		return nil, ErrClickOpen(err)
	}
	s.Instance = conn
	v, err := conn.ServerVersion()
	if err != nil {
		return nil, ErrClickGetVer(err)
	}
	log().Pr("click").Cmp(config.User).Mth("open").F(kit.KV{"version": v}).Inf("ok")
	return s, nil
}

func OpenDb(config *Config, logger logger.CLoggerFunc) (*sql.DB, error) {

	// make connection
	conn := clickhouse.OpenDB(cfgToOptions(config, logger))

	// ping
	err := conn.Ping()
	if err != nil {
		return nil, ErrClickPing(err)
	}

	return conn, nil
}

func (s *ClickHouse) l() logger.CLogger {
	return s.log().Cmp("click")
}

func (s *ClickHouse) Close() {
	if s.Instance != nil {
		_ = s.Instance.Close()
		s.Instance = nil
	}
	s.log().Cmp("click").Mth("close").Inf("ok")
}

func cfgToOptions(config *Config, logger logger.CLoggerFunc) *clickhouse.Options {
	return &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},
		Debug: config.Debug,
		Debugf: func(format string, v ...interface{}) {
			logger().Cmp("click").Mth("debug").DbgF(format, v...)
		},
	}
}
