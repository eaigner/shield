package shield

import "testing"

var (
	redisStore = NewRedisStore("127.0.0.1:6379", "", logger, "redis")
)

func TestRedisLearn(t *testing.T) {
	sh := newShield(redisStore)
	testLearn(t, sh)
}

func TestRedisDecrement(t *testing.T) {
	sh := newShield(redisStore)
	testDecrement(t, sh)
}
