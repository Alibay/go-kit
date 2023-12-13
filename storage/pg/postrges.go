package pg

import (
	"fmt"
	"time"

	"github.com/Alibay/go-kit/logger"

	kit "github.com/Alibay/go-kit"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Storage struct {
	Instance *gorm.DB
	DBName   string
	log      logger.CLoggerFunc
}

// DbConfig database configuration
type DbConfig struct {
	ConnectionString string `mapstructure:"connection_string"` // ConnectionString if specified overrides all other connection params
	User             string
	Password         string
	DBName           string
	Port             string
	Host             string
}

// DbClusterConfig configuration of database cluster
type DbClusterConfig struct {
	Master *DbConfig // Master database
	Slave  *DbConfig // Slave database
}

func Open(config *DbConfig, log logger.CLoggerFunc) (*Storage, error) {

	s := &Storage{
		DBName: config.DBName,
		log:    log,
	}

	dsn := config.ConnectionString
	if dsn == "" {
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s",
			config.User,
			config.Password,
			config.DBName,
			config.Port,
			config.Host,
		)
	}

	// uncomment to log all queries
	cfg := &gorm.Config{
		Logger: gormLogger.New(
			log(),
			//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold: time.Second * 10, // Slow SQL threshold
				LogLevel:      gormLogger.Info,  // Log level
				Colorful:      false,            // Disable color
			},
		),
		NowFunc: func() time.Time { return kit.Now() },
	}

	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, ErrPostgresOpen(err)
	}

	log().Pr("db").Cmp(config.User).Inf("ok")

	s.Instance = db

	return s, nil
}

func (s *Storage) Close() {
	if s.Instance != nil {
		db, _ := s.Instance.DB()
		_ = db.Close()
	}
}