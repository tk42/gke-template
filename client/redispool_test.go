package client

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func Test_GetConn(t *testing.T) {
	once = new(sync.Once)

	os.Setenv("LOGLEVEL", "DEBUG")
	os.Setenv("TEST_ENV", "false")

	p := GetRedisConnPool()
	c := p.Get()
	s, _ := redis.Strings(c.Do("KEYS", "*"))
	fmt.Println(s)
	// Output:
	// [ ]
}

func Test_GetConn_Test(t *testing.T) {
	once = new(sync.Once)

	os.Setenv("LOGLEVEL", "DEBUG")
	os.Setenv("TEST_ENV", "true")

	p := GetRedisConnPool()
	c := p.Get()
	s, _ := redis.Strings(c.Do("KEYS", "*"))
	fmt.Println(s)
	// Output:
	// [ ]
}
