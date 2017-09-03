package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	defaultMaxIdle     = 500
	defaultMaxActive   = 0
	defaultIdleTimeout = 20 * time.Second
)

type Redis struct {
	host       string
	port       int
	password   string
	db         int
	enablePool bool
	timeout    time.Duration
	pool       *redis.Pool
}

// New Redis instance
func NewRedis(redisInfo map[string]interface{}, enablePool bool, timeout time.Duration) (*Redis, error) {
	if redisInfo["host"] == "" || redisInfo["port"] == 0 {
		return nil, errors.New("redisInfo ip or port should not be empty")
	}

	r := &Redis{
		host:       redisInfo["host"].(string),
		port:       redisInfo["port"].(int),
		timeout:    timeout,
		enablePool: enablePool,
	}

	if redisInfo["password"] != nil {
		r.password = redisInfo["password"].(string)
	}
	if redisInfo["db"] != nil {
		r.db = redisInfo["db"].(int)
	}

	if enablePool {
		r.pool = &redis.Pool{
			MaxIdle:     defaultMaxIdle,
			MaxActive:   defaultMaxActive,
			IdleTimeout: defaultIdleTimeout,
			Wait:        true,
			Dial:        r.dial,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}

	return r, nil
}

// Connect
func (this *Redis) Connect() (redis.Conn, error) {
	if this.enablePool {
		return this.pool.Get(), nil
	}

	return this.dial()
}

func (this *Redis) dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", this.host, this.port),
		redis.DialConnectTimeout(this.timeout),
		redis.DialReadTimeout(this.timeout),
		redis.DialWriteTimeout(this.timeout),
		redis.DialDatabase(this.db),
		redis.DialPassword(this.password))

	if err != nil {
		return nil, err
	}

	return c, nil
}
