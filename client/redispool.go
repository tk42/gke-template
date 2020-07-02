package client

import (
	"sync"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/jimako1989/gke-template/env"
	"github.com/jimako1989/gke-template/logging"
	"go.uber.org/zap"
)

type RedigoPool interface {
	Get() redigo.Conn
	Close() error
}

// MEMO: Pool handles connection pool. cf:https://qiita.com/riverplus/items/12f9b37cf1795d9bdbb1
type RedisConnPool struct {
	logger zap.Logger
	pool   RedigoPool
}

var once sync.Once
var redisConn *RedisConnPool

func GetRedisConnPool() *RedisConnPool {
	once.Do(func() {
		logger := logging.GetLogger("RedisConn")
		address := env.GetString("REDIS_HOST", "localhost") + ":" + env.GetString("REDIS_PORT", "6379")
		logger.Info("Loaded Redis address", zap.String("address", address))
		pool := &redigo.Pool{
			MaxIdle:     env.GetInt("REDIS_MAX_IDLE_NUM", 20),
			MaxActive:   env.GetInt("REDIS_MAX_ACTIVE_NUM", 20),
			Wait:        false, // true: blocking until the number of connections is under MaxActive
			IdleTimeout: 240 * time.Second,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", address)
			},
		}
		logger.Info("Starting to connect to Redis", zap.String("address", address))

		redisConn = &RedisConnPool{
			logger: logger,
			pool:   pool,
		}
	})
	return redisConn
}

func (rcp *RedisConnPool) getConn() redigo.Conn {
	var redisConn redigo.Conn
	for {
		redisConn = rcp.pool.Get()
		if redisConn.Err() != nil {
			rcp.logger.Error("Failed to get a connection from Redis pool", zap.Error(redisConn.Err()))
			time.Sleep(1 * time.Minute)
			continue
		}
		if redisConn == nil {
			rcp.logger.Error("Failed to get redis connection")
			time.Sleep(1 * time.Minute)
			continue
		}
		break
	}
	return redisConn
}

func GetConn(tableNo string) redigo.Conn {
	rcp := GetRedisConnPool()
	conn := rcp.getConn()
	conn.Do("SELECT", tableNo)
	return conn
}
