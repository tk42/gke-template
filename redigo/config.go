package redigo

type PoolConfiguration struct {
	isMock               bool
	host                 string
	port                 string
	db                   uint32
	maxIdleConnections   int
	maxActiveConnections int
}

type Option func(*PoolConfiguration)

func IsMock(isMock bool) Option {
	return func(c *PoolConfiguration) {
		c.isMock = isMock
	}
}

func Host(host string) Option {
	return func(c *PoolConfiguration) {
		c.host = host
	}
}

func Port(port string) Option {
	return func(c *PoolConfiguration) {
		c.port = port
	}
}

func DB(db uint32) Option {
	return func(c *PoolConfiguration) {
		c.db = db
	}
}

func MaxIdleConnections(maxIdleConnections int) Option {
	return func(c *PoolConfiguration) {
		c.maxIdleConnections = maxIdleConnections
	}
}

func MaxActiveConnections(maxActiveConnections int) Option {
	return func(c *PoolConfiguration) {
		c.maxActiveConnections = maxActiveConnections
	}
}

func PoolConfig(ops ...Option) PoolConfiguration {
	cfg := PoolConfiguration{
		isMock:               false,
		host:                 "localhost",
		port:                 "6379",
		db:                   0,
		maxIdleConnections:   20,
		maxActiveConnections: 20,
	}
	for _, option := range ops {
		option(&cfg)
	}
	return cfg
}
