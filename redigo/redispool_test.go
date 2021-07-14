package redigo

import (
	"fmt"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/tk42/victolinux/env"
)

var (
	config = PoolConfig(
		IsMock(env.GetBoolean("TEST_ENV", true)),
		Host(env.GetString("REDIS_HOST", "localhost")),
		Port(env.GetString("REDIS_PORT", "6379")),
		MaxIdleConnections(env.GetInt("REDIS_MAX_IDLE_NUM", 20)),
		MaxActiveConnections(env.GetInt("REDIS_MAX_ACTIVE_NUM", 20)),
	)
)

func Test_GetConn(t *testing.T) {
	os.Setenv("LOGLEVEL", "DEBUG")
	os.Setenv("TEST_ENV", "false")

	p := GetRedisConnPool(config)
	c, _ := p.GetContext(nil)
	s, _ := redis.Strings(c.Do("KEYS", "*"))
	fmt.Println(s)
	// Output:
	// [ ]
}

func Test_GetConn_Test(t *testing.T) {
	os.Setenv("LOGLEVEL", "DEBUG")
	os.Setenv("TEST_ENV", "true")

	p := GetRedisConnPool(config)
	c, _ := p.GetContext(nil)
	s, _ := redis.Strings(c.Do("KEYS", "*"))
	fmt.Println(s)
	// Output:
	// [ ]
}
