package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jimako1989/gke-template/env"
	"github.com/jimako1989/gke-template/logging"
	dockertest "github.com/ory/dockertest/v3"
	"go.uber.org/zap"
)

// Pool contains dockertest and redis connection pool
type Pool struct {
	logger    zap.Logger
	redisPool *redis.Pool
}

var (
	once          = new(sync.Once)
	redisConnPool *Pool
)

func GetRedisConnPool() *Pool {
	once.Do(func() {
		logger := logging.GetLogger("RedisConn")
		isTest := env.GetBoolean("TEST_ENV", true)

		var dialFunc func() (redis.Conn, error)
		if isTest {
			dockerPool, err := dockertest.NewPool("")
			if err != nil {
				logger.Fatal("could not connect to docker", zap.Error(err))
			}
			dockerRes, err := dockerPool.Run("redis", "5.0", nil)
			if err != nil {
				logger.Fatal("could not start resource", zap.Error(err))
			}
			logger.Info("Loaded Test Redis")
			dialFunc = func() (redis.Conn, error) {
				return redis.DialURL(fmt.Sprintf("redis://localhost:%s", dockerRes.GetPort("6379/tcp")))
			}
		} else {
			address := env.GetString("REDIS_HOST", "localhost") + ":" + env.GetString("REDIS_PORT", "6379")
			logger.Info("Loaded Redis address", zap.String("address", address))
			dialFunc = func() (redis.Conn, error) {
				return redis.Dial("tcp", address)
			}
		}

		redisConnPool = &Pool{
			redisPool: &redis.Pool{
				MaxIdle:     env.GetInt("REDIS_MAX_IDLE_NUM", 20),
				MaxActive:   env.GetInt("REDIS_MAX_ACTIVE_NUM", 20),
				Wait:        false, // true: blocking until the number of connections is under MaxActive
				IdleTimeout: 240 * time.Second,
				Dial:        dialFunc,
			},
		}

		logger.Info("Starting to connect to Redis")
	})
	return redisConnPool
}

// Get gets a connection with redis
func (p *Pool) Get() redis.Conn {
	var redisConn redis.Conn
	for {
		redisConn = p.redisPool.Get()
		if redisConn.Err() != nil {
			p.logger.Error("Failed to get a connection from Redis pool", zap.Error(redisConn.Err()))
			time.Sleep(1 * time.Minute)
			continue
		}
		if redisConn == nil {
			p.logger.Error("Failed to get redis connection")
			time.Sleep(1 * time.Minute)
			continue
		}
		break
	}
	return redisConn
}

// GetContext gets a connection with redis
func (p *Pool) GetContext(ctx context.Context) (redis.Conn, error) {
	return p.redisPool.GetContext(ctx)
}

// GetPool gets a connection pool with redis
func (p *Pool) GetPool() *redis.Pool {
	return p.redisPool
}

// Cleanup remove all data in redis
func (p *Pool) Cleanup() error {
	conn := p.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

// Close closes redis connection pool and dockertest pool
func (p *Pool) Close() {
	var errs []error
	if err := p.Cleanup(); err != nil {
		errs = append(errs, err)
	}
	if err := p.redisPool.Close(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		log.Fatalf("unexpected error: %v", errs[0])
	}
}

func GetConn(dbNo string) redis.Conn {
	conn := GetRedisConnPool().Get()
	conn.Do("SELECT", dbNo)
	return conn
}
