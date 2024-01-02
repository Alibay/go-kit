package cluster

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Alibay/go-kit"
	"github.com/Alibay/go-kit/storages/clickhouse"
	"github.com/Alibay/go-kit/storages/migration"
	kitStorage "github.com/Alibay/go-kit/storages/pg"
	"github.com/spf13/cobra"
)

// Bootstrap interface needs to be implemented by service instance implementation to handle service's lifecycle
type Bootstrap interface {
	// Init initializes the service
	Init(ctx context.Context, cfg any) error
	// Start executes all background processes
	Start(ctx context.Context) error
	// Close closes the service
	Close(ctx context.Context)
}

const (
	configFlagName             = "config"
	configDefaultPath          = "./config.yml"
	migrationSourceFlagName    = "source"
	migrationSourceDefaultPath = "./db/migrations"

	ErrCodeConfigPathNotSpecified     = "SVS-001"
	ErrCodeMigrationSourcePathInvalid = "SVS-002"
)

var (
	ErrConfigPathNotSpecified = func() error {
		return kit.NewAppErrBuilder(ErrCodeConfigPathNotSpecified, "").Business().Err()
	}
	ErrMigrationSourcePathInvalid = func() error {
		return kit.NewAppErrBuilder(ErrCodeMigrationSourcePathInvalid, "migration source path isn't valid").Business().Err()
	}
)

type ServiceInstance[TCfg any] struct {
	svcCode      string         // svcCode unique identifier of the service
	nodeId       string         // nodeId service node ID
	instanceId   string         // instanceId service instance when multiple replicas are running
	bootstrap    Bootstrap      // bootstrap implementation of Bootstrap interface
	rootCmd      *cobra.Command // rootCmd root command
	confPathEnv  *string        // confPathEnv env var name of config path (optional)
	migSourceEnv *string        // migSourceEnv env var name of migration path (optional)
	logger       *kit.Logger
}

func New[TCfg any](svcCode string, bootstrap Bootstrap) *ServiceInstance[TCfg] {

	s := &ServiceInstance[TCfg]{
		bootstrap: bootstrap,
		svcCode:   svcCode,
		logger:    kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel, Format: kit.FormatterJson}),
	}

	// init root command
	s.rootCmd = &cobra.Command{
		Use: svcCode,
	}
	flags := s.rootCmd.PersistentFlags()
	flags.String(
		configFlagName,
		configDefaultPath,
		"--config <path-to-file>",
	)

	// app command
	appCmd := &cobra.Command{
		Use: "app",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.executeAppCmd(cmd, args)
		},
	}
	s.rootCmd.AddCommand(appCmd)

	return s
}

func (s *ServiceInstance[TCfg]) l() kit.CLogger {
	return s.GetLogger()()
}

// WithConfigPathEnv allows to specify env which is used to obtain a config path. It has priority over command flags
func (s *ServiceInstance[TCfg]) WithConfigPathEnv(env string) *ServiceInstance[TCfg] {
	if env != "" {
		s.confPathEnv = &env
	}
	return s
}

// WithMigrationSourceEnv allows to specify env which is used to obtain a migration source folder path. It has priority over command flags
func (s *ServiceInstance[TCfg]) WithMigrationSourceEnv(env string) *ServiceInstance[TCfg] {
	if env != "" {
		s.migSourceEnv = &env
	}
	return s
}

func (s *ServiceInstance[TCfg]) WithDbMigration(getDbConfigFn func(cfg *TCfg) (any, error)) *ServiceInstance[TCfg] {

	// create migration commands
	dbUpCmd := &cobra.Command{
		Use: "db-up",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.executePgCmd(cmd, getDbConfigFn, true)
		},
	}
	dbDownCmd := &cobra.Command{
		Use: "db-down",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.executePgCmd(cmd, getDbConfigFn, false)
		},
	}

	s.rootCmd.AddCommand(dbUpCmd, dbDownCmd)

	// set flags
	flags := dbUpCmd.PersistentFlags()
	flags.String(
		migrationSourceFlagName,
		migrationSourceDefaultPath,
		"--source <path-to-migration-folder>",
	)

	flags = dbDownCmd.PersistentFlags()
	flags.String(
		migrationSourceFlagName,
		migrationSourceDefaultPath,
		"--source <path-to-migration-folder>",
	)

	return s
}

func (s *ServiceInstance[TCfg]) WithClickHouseMigration(getClickConfigFn func(cfg *TCfg) (any, error)) *ServiceInstance[TCfg] {

	// create migration commands
	dbUpCmd := &cobra.Command{
		Use: "db-up",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.executeClickCmd(cmd, getClickConfigFn, true)
		},
	}
	dbDownCmd := &cobra.Command{
		Use: "db-down",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.executeClickCmd(cmd, getClickConfigFn, false)
		},
	}

	s.rootCmd.AddCommand(dbUpCmd, dbDownCmd)

	// set flags
	flags := dbUpCmd.PersistentFlags()
	flags.String(
		migrationSourceFlagName,
		migrationSourceDefaultPath,
		"--source <path-to-migration-folder>",
	)
	flags = dbDownCmd.PersistentFlags()
	flags.String(
		migrationSourceFlagName,
		migrationSourceDefaultPath,
		"--source <path-to-migration-folder>",
	)

	return s
}

func (s *ServiceInstance[TCfg]) GetCode() string {
	return s.svcCode
}

func (s *ServiceInstance[TCfg]) NodeId() string {
	return s.nodeId
}

func (s *ServiceInstance[TCfg]) Execute() error {
	return s.rootCmd.Execute()
}

func (s *ServiceInstance[TCfg]) GetLogger() kit.CLoggerFunc {
	return func() kit.CLogger {
		return kit.L(s.logger).Srv(s.svcCode).Node(s.nodeId)
	}
}

func (s *ServiceInstance[TCfg]) loadConfig(cmd *cobra.Command) (*TCfg, error) {
	l := s.l().Mth("load-config")

	configLoader := kit.NewConfigLoader[TCfg]()
	var path string

	// get from cmd flag
	if f := cmd.Flag(configFlagName); f != nil {
		path = f.Value.String()
		if path != "" {
			configLoader = configLoader.WithPath(path)
		}
	} else if s.confPathEnv != nil {
		// try to load config from the path passed by env
		configLoader = configLoader.WithEnv(*s.confPathEnv)
		path = os.Getenv(*s.confPathEnv)
	} else {
		return nil, ErrConfigPathNotSpecified()
	}

	// try to load from passed path
	absPath, _ := filepath.Abs(path)
	config, err := configLoader.Load()
	if err != nil {
		return nil, err
	}
	if config != nil {
		l.DbgF("found: %s", absPath).TrcObj("%v", config)
	}

	return config, nil
}

func (s *ServiceInstance[TCfg]) loadMigrationSource(cmd *cobra.Command) (string, error) {

	var absPath string

	// get from cmd flag
	if f := cmd.Flag(migrationSourceFlagName); f != nil {
		source := f.Value.String()
		if source != "" {
			absPath, _ = filepath.Abs(source)
		}
	}

	// try to load migration source from the path passed by env
	if absPath == "" && s.migSourceEnv != nil {
		source := os.Getenv(*s.migSourceEnv)
		if source != "" {
			absPath, _ = filepath.Abs(source)
		}
	}

	if absPath == "" {
		return "", ErrMigrationSourcePathInvalid()
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", ErrMigrationSourcePathInvalid()
	}

	return absPath, nil
}

func (s *ServiceInstance[TCfg]) executeAppCmd(cmd *cobra.Command, args []string) error {
	l := s.l().Mth("app")

	// load config
	config, err := s.loadConfig(cmd)
	if err != nil {
		return err
	}

	// init context
	ctx, cancelFn := context.WithCancel(kit.NewRequestCtx().Empty().WithNewRequestId().ToContext(context.Background()))
	defer cancelFn()

	// init service
	if err = s.bootstrap.Init(ctx, config); err != nil {
		l.E(err).Err("init: fail")
		return err
	}
	l.Inf("init: ok")

	// start listening
	if err := s.bootstrap.Start(ctx); err != nil {
		l.C(ctx).E(err).Err("start: fail")
		return err
	}
	l.Inf("started: ok")

	// close
	defer func() {
		s.bootstrap.Close(ctx)
		l.Inf("graceful shutdown")
	}()

	// handle app close
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	return nil
}

func (s *ServiceInstance[TCfg]) executePgCmd(cmd *cobra.Command, getDbConfigFn func(cfg *TCfg) (any, error), up bool) error {

	// load config
	config, err := s.loadConfig(cmd)
	if err != nil {
		return err
	}

	// extract config
	dbConfig, err := getDbConfigFn(config)
	if err != nil {
		return err
	}

	// load migration source
	src, err := s.loadMigrationSource(cmd)
	if err != nil {
		return err
	}

	// build function opening database
	openDb := func() (*sql.DB, error) {
		pg, err := kitStorage.Open(dbConfig.(*kitStorage.DbConfig), s.GetLogger())
		if err != nil {
			return nil, err
		}
		db, _ := pg.Instance.DB()
		return db, nil
	}

	// we count on that migrations is located either in "src" folder or in "src/pg" folder
	return s.executeMigrationCmd(openDb, migration.DialectPostgres, []string{fmt.Sprintf("%s/pg", src), src}, up)
}

func (s *ServiceInstance[TCfg]) executeClickCmd(cmd *cobra.Command, getClickConfigFn func(cfg *TCfg) (any, error), up bool) error {

	// load config
	config, err := s.loadConfig(cmd)
	if err != nil {
		return err
	}

	// extract config
	dbConfig, err := getClickConfigFn(config)
	if err != nil {
		return err
	}

	// load migration source
	src, err := s.loadMigrationSource(cmd)
	if err != nil {
		return err
	}

	// build function opening database
	openDb := func() (*sql.DB, error) {
		db, err := clickhouse.OpenDb(dbConfig.(*clickhouse.Config), s.GetLogger())
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	// we count on that migrations is located either in "src" folder or in "src/click" folder
	return s.executeMigrationCmd(openDb, migration.DialectClickHouse, []string{fmt.Sprintf("%s/click", src), src}, up)
}

func (s *ServiceInstance[TCfg]) executeMigrationCmd(openDbFn func() (*sql.DB, error), dialect string, srcPaths []string, up bool) error {

	// check paths
	var srcPath string
	for _, p := range srcPaths {
		// check folder exists
		absPath, _ := filepath.Abs(p)
		if _, err := os.Stat(absPath); err == nil {
			srcPath = absPath
			break
		}
	}

	// open db
	sqlDb, err := openDbFn()
	if err != nil {
		return err
	}
	defer func() { _ = sqlDb.Close() }()

	// migration
	m := migration.NewMigration(sqlDb, srcPath, s.GetLogger(), dialect)

	// run migration command
	if up {
		return m.Up()
	}
	return m.Down()
}
