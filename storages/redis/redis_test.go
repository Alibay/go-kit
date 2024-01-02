///go:build integration

package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
)

type redisTestSuite struct {
	kit.Suite
}

func (s *redisTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(redisTestSuite))
}

var (
	config = &Config{
		Host: "localhost",
		Port: "6379",
		Ttl:  0,
	}
)

func (s *redisTestSuite) Test_Range() {

	cl, err := Open(s.Ctx, config, s.L)
	s.NoError(err)
	defer cl.Close()

	key := kit.NewRandString()
	jsons, err := cl.Instance.LRange(s.Ctx, key, 0, -1).Result()
	s.NoError(err)
	fmt.Println(jsons)

	pipe := cl.Instance.Pipeline()
	pipe.Expire(s.Ctx, key, time.Second*10)

	s.NoError(cl.Instance.RPush(s.Ctx, key, "1").Err())
	s.NoError(cl.Instance.RPush(s.Ctx, key, "2").Err())
	s.NoError(cl.Instance.RPush(s.Ctx, key, "3").Err())

	_, err = pipe.Exec(s.Ctx)
	s.NoError(err)

	jsons, err = cl.Instance.LRange(s.Ctx, key, 0, -1).Result()
	s.NoError(err)
	s.Equal(3, len(jsons))
}

func (s *redisTestSuite) Test_Distributed_Lock() {
	cl, err := Open(s.Ctx, config, s.L)
	s.NoError(err)
	defer cl.Close()

	key, unlockId := kit.NewRandString(), kit.NewRandString()

	// apply lock
	locked, err := cl.Lock(s.Ctx, key, unlockId, time.Second*10)
	s.NoError(err)
	s.True(locked)

	// apply lock again
	locked, err = cl.Lock(s.Ctx, key, unlockId, time.Second*10)
	s.NoError(err)
	s.False(locked)

	// try to lock with another unlockId
	locked, err = cl.Lock(s.Ctx, key, kit.NewRandString(), time.Second*10)
	s.NoError(err)
	s.False(locked)

	// try to unlock with another unlock ID
	unlocked, err := cl.UnLock(s.Ctx, key, kit.NewRandString())
	s.NoError(err)
	s.False(unlocked)

	// try to unlock with another unlock ID
	unlocked, err = cl.UnLock(s.Ctx, key, unlockId)
	s.NoError(err)
	s.True(unlocked)

}
