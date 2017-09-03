package dbsvc

import (
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"

	redisc "scheduler/common/redis"
)

var (
	redisInfo = map[string]interface{}{
		"host": "your ip",
		"port": 6666, //your redis port
	}
	r *redisc.Redis
)

func init() {
	var err error
	if r, err = redisc.NewRedis(redisInfo, true, 2*time.Second); err != nil {
		fmt.Println("error occured when New redis: ", err)
		os.Exit(1)
	}
}

func Setex(key string, seconds int64, value string) (int, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return 0, err
	}
	if _, serr := conn.Do("SETEX", key, seconds, value); serr != nil {
		return 0, serr
	}
	return 1, nil
}

func Exists(key string) (int, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return 0, err
	}
	e, eerr := redis.Int(conn.Do("EXISTS", key))
	if eerr != nil {
		return 0, eerr
	}
	return e, nil
}

func Pop(list string) (string, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return "", err
	}
	ele, perr := redis.String(conn.Do("LPOP", list))
	if perr != nil {
		return "", perr
	}
	return ele, nil
}

func Push(list, ele string) (int, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return 0, err
	}
	size, perr := redis.Int(conn.Do("RPUSH", list, ele))
	if perr != nil {
		return 0, perr
	}
	return size, nil
}

func Setnx(key, value string) (int, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return 0, err
	}
	ok, serr := redis.Int(conn.Do("SETNX", key, value))
	if serr != nil {
		return 0, serr
	}
	return ok, nil
}

func Del(key string) (int, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return 0, err
	}
	_, derr := conn.Do("DEL", key)
	if derr != nil {
		return 0, derr
	}
	return 1, nil
}

func Get(key string) (string, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return "", err
	}
	value, gerr := redis.String(conn.Do("GET", key))
	if gerr != nil {
		return "", gerr
	}
	return value, nil
}

func GetSet(key string, value string) (string, error) {
	conn, err := r.Connect()
	defer conn.Close()
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			conn, err = r.Connect()
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return "", err
	}
	_value, gerr := redis.String(conn.Do("GETSET", key, value))
	if gerr != nil {
		return "", gerr
	}
	return _value, nil
}
