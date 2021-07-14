package redigo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/tk42/victolinux/logging"
	"go.uber.org/zap"
)

// Pool contains dockertest and redis connection pool
type Pool struct {
	logger    *logging.Logger
	redisPool *redis.Pool
	dockerRes *dockertest.Resource
}

var (
	once          = new(sync.Once)
	redisConnPool *Pool
	PoolCache     map[uint32]*Pool
)

func GetRedisConnPool(config PoolConfiguration) *Pool {
	var pool *Pool
	once.Do(func() {
		pool = getRedisConnPoolByDB(config)
	})
	return pool
}

func GetRedisConnPoolByDB(config PoolConfiguration) *Pool {
	if _, ok := PoolCache[config.db]; !ok {
		PoolCache[config.db] = getRedisConnPoolByDB(config)
	}
	return PoolCache[config.db]
}

func getRedisConnPoolByDB(config PoolConfiguration) *Pool {
	logger := logging.GetLogger("RedisConn")

	var dockerRes *dockertest.Resource

	var dialFunc func() (redis.Conn, error)
	if config.isMock {
		dockerPool, err := dockertest.NewPool("")
		if err != nil {
			logger.Fatal("could not connect to docker", zap.Error(err))
		}
		dockerRes, err = dockerPool.Run("redis", "5.0", nil)
		if err != nil {
			logger.Fatal("could not start resource", zap.Error(err))
		}
		logger.Info("Loaded Test Redis")
		dialFunc = func() (redis.Conn, error) {
			return redis.DialURL(fmt.Sprintf("redis://localhost:%s/%d", dockerRes.GetPort("6379/tcp"), config.db))
		}
	} else {
		address := fmt.Sprintf("%v:%v/%v", config.host, config.port, config.db)
		if config.host == "localhost" {
			logger.Warn("Loaded Redis(localhost) address", zap.String("address", address))
		} else {
			logger.Info("Loaded Redis address", zap.String("address", address))
		}
		dialFunc = func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		}
	}

	logger.Info("Starting to connect to Redis")

	return &Pool{
		logger: logger,
		redisPool: &redis.Pool{
			MaxIdle:     config.maxIdleConnections,
			MaxActive:   config.maxActiveConnections,
			Wait:        false, // true: blocking until the number of connections is under MaxActive
			IdleTimeout: 240 * time.Second,
			Dial:        dialFunc,
		},
		dockerRes: dockerRes,
	}
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
	conn, err := p.GetContext(nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("FLUSHALL")
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
	if p.dockerRes != nil {
		if err := p.dockerRes.Close(); err != nil {
			errs = append(errs, err)
		}
		once = new(sync.Once)
	}
	if len(errs) > 0 {
		log.Fatalf("unexpected error: %v", errs[0])
	}
}
