package redigo

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/qjw/kelly/sessions/cache"
	parse "github.com/qjw/url"
)

type rediGo struct {
	conn redis.Conn
}

func (this rediGo) Get(key string) ([]byte, error) {
	return redis.Bytes(this.conn.Do("GET", key))
}

func (this rediGo) Set(key string, value []byte) error {
	_, err := this.conn.Do("SET", key, value)
	return err
}

func (this rediGo) SetWithExpire(key string, value []byte, expire time.Duration) error {
	_, err := this.conn.Do("SET", key, value, "EX", int64(expire.Seconds()))
	return err
}

func (this rediGo) Delete(key string) error {
	_, err := this.conn.Do("DEL", key)
	return err
}

func (this rediGo) Expire(key string, expire time.Duration) error {
	_, err := this.conn.Do("EXPIRE", key, expire.Seconds())
	return err
}

func (this rediGo) Exisits(key string) error {
	e, err := redis.Int(this.conn.Do("EXISTS", key))
	if err == nil && e > 0 {
		return nil
	} else {
		if err != nil {
			return err
		} else {
			return fmt.Errorf("not exists")
		}
	}
}

func (this *rediGo) Close() error {
	if this.conn != nil {
		err := this.conn.Close()
		if err != nil {
			return err
		}
		this.conn = nil
	}
	return nil
}

func NewCache(url string) (cache.Cache, error) {
	config, err := parse.ParseRedis(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url %s,error %s", url, err.Error())
	}

	// 连接redis
	params := []redis.DialOption{
		redis.DialDatabase(config.Db),
		redis.DialConnectTimeout(5 * time.Second),
		redis.DialReadTimeout(5 * time.Second),
		redis.DialWriteTimeout(5 * time.Second),
	}
	if config.Password != nil {
		params = append(params, redis.DialPassword(*config.Password))
	}
	c, err := redis.Dial("tcp", config.Host, params...)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil, err
	}

	return &rediGo{conn: c}, nil
}
