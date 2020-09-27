package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func Test_GetConn(t *testing.T) {
	os.Setenv("LOGLEVEL", "DEBUG")

	p := GetRedisConnPool()
	c := p.Get()
	s, _ := redis.Strings(c.Do("KEYS", "*"))
	fmt.Println(s)
	// Output:
	// [ ]
}
