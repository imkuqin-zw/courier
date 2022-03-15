package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	conf := &Config{
		Addrs: []string{"localhost:6379"},
	}
	redisCli := conf.Build()
	err := redisCli.Ping()
	require.NoError(t, err, "ping redislock")
	err = redisCli.Set("_test_redis", 1, time.Second)
	require.NoError(t, err, "set")
	val, err := redisCli.Get("_test_redis").Int()
	require.NoError(t, err, "get")
	require.Equal(t, 1, val, "get")
}
