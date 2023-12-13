//go:build integration

package clickhouse

import (
	"fmt"
	"testing"

	kit "github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
)

type clickHouseTestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
}

func (s *clickHouseTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func TestClickHouseSuite(t *testing.T) {
	suite.Run(t, new(clickHouseTestSuite))
}

var (
	config = &Config{
		User:     "admin",
		Password: "admin",
		Database: "ev2go",
		Port:     "19000",
		Host:     "127.0.0.1",
	}
)

func (s *clickHouseTestSuite) Test_CreateDropTable_Select() {
	// open database
	ch, err := Open(config, s.logger)
	s.NoError(err)
	// close database
	defer func() {
		ch.Close()
	}()
	// ping
	s.NoError(ch.Instance.Ping(s.Ctx))
	// create table
	s.NoError(ch.Instance.Exec(s.Ctx, "DROP TABLE IF EXISTS _test"))
	s.NoError(ch.Instance.Exec(s.Ctx, "CREATE TABLE IF NOT EXISTS _test (i Int64) ENGINE=MergeTree() order by (i)"))
	// drop table
	defer func() {
		s.NoError(ch.Instance.Exec(s.Ctx, "DROP TABLE _test"))
	}()
	// insert into table
	s.NoError(ch.Instance.Exec(s.Ctx, fmt.Sprintf("INSERT INTO _test VALUES(%d)", 1)))
	// select from table
	var result struct {
		Col1  int64  `ch:"col"`
		Count uint64 `ch:"count"`
	}
	s.NoError(ch.Instance.QueryRow(s.Ctx, "SELECT i as col, count() as count FROM _test GROUP BY i").ScanStruct(&result))
	s.NotEmpty(result)
	s.Equal(int64(1), result.Col1)
	s.Equal(uint64(1), result.Count)
}
