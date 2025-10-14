package constant

import "time"

var (
	RespNil                      = []byte("$-1\r\n")
	RespOk                       = []byte("+OK\r\n")
	RespZero                     = []byte(":0\r\n")
	RespOne                      = []byte(":1\r\n")
	TtlKeyNotExist               = []byte(":-2\r\n")
	TtlKeyExistNoExpire          = []byte(":-1\r\n")
	ActiveExpireFrequency        = 100 * time.Millisecond
	ActiveExpireSampleSize       = 20                    // https://blog.x.com/engineering/en_us/topics/infrastructure/2019/improving-key-expiration-in-redis
	ActiveExpireThreshold        = 0.25                  // https://engineering.grab.com/a-key-expired-in-redis-you-wont-believe-what-happened-next
	MaxActiveExpireExecutionTime = 25 * time.Millisecond // https://groups.google.com/g/redis-db/c/tF1cIg-bXS0
	IOMultiplexerTimeout         = 50 * time.Millisecond
	DefaultBPlusTreeDegree       = 64 // https://timmastny.com/blog/tuning-b-plus-trees/
)

const (
	BfDefaultInitCapacity = 100
	BfDefaultErrRate      = 0.01
)

const (
	ServerStatusIdle         = 1
	ServerStatusBusy         = 2
	ServerStatusShuttingDown = 3
)
