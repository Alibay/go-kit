package redis

import kit "github.com/Alibay/go-kit"

const (
	ErrCodeRedisPingErr = "RDS-001"
)

var (
	ErrRedisPingErr = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeRedisPingErr, "").Wrap(cause).Err()
	}
)
