package constant

import "time"

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("+OK\r\n")
var RespZero = []byte(":0\r\n")
var RespOne = []byte(":1\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
var ActiveExpireFrequency = 100 * time.Millisecond
var ActiveExpireSampleSize = 20 // https://blog.x.com/engineering/en_us/topics/infrastructure/2019/improving-key-expiration-in-redis
var ActiveExpireThreshold = 0.25 // https://engineering.grab.com/a-key-expired-in-redis-you-wont-believe-what-happened-next
var MaxActiveExpireExecutionTime = 25 * time.Millisecond // https://groups.google.com/g/redis-db/c/tF1cIg-bXS0
var IOMultiplexerTimeout = 50 * time.Millisecond 
var DefaultBPlusTreeDegree = 64 // https://timmastny.com/blog/tuning-b-plus-trees/