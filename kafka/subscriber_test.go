package kafka

import (
	"testing"

	"github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
)

type subscriberTestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
}

func (s *subscriberTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func TestSubscriberSuite(t *testing.T) {
	suite.Run(t, new(subscriberTestSuite))
}

func (s *subscriberTestSuite) Test_IndexByKey() {

	test := func(workers int, keys []string, exp ...int) {
		sub := &subscriber{workers: workers}
		for i, k := range keys {
			s.Equal(exp[i], sub.chanIndexByKey([]byte(k)))
		}
	}

	test(1, []string{"1", "2", "33244", kit.NewRandString(), "AAaaFFFff"}, 0, 0, 0, 0, 0)
	test(2, []string{"1", "1", "2", "2"}, 0, 0, 1, 1)
	test(2, []string{"aaFFaaFF", "bbCCbbCD", "aaFFaaFF", "bbCCbbCD"}, 1, 0, 1, 0)

	randKey := kit.NewRandString()
	sub := &subscriber{workers: 10}
	s.Equal(sub.chanIndexByKey([]byte(randKey)), sub.chanIndexByKey([]byte(randKey)))

}
