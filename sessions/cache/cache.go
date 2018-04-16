package cache

import (
	"fmt"
	"time"
)

var CacheNotExistError = fmt.Errorf("key not found")

type Cache interface {
	Get(string) ([]byte, error)
	Set(key string, value []byte) error
	SetWithExpire(key string, value []byte, expire time.Duration) error
	Delete(string) error
	Expire(string, time.Duration) error
	Exisits(string) error
	Close() error
}
