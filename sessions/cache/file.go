package cache

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type encoder interface {
	Encode(f *os.File, d *p) error
	Decode(f *os.File) (*p, error)
}

type fileCache struct {
	path string
	d    encoder
	m    sync.Mutex
}

type p struct {
	Expire int64
	Value  []byte
}

func (this fileCache) getPath(key string) string {
	return fmt.Sprintf("%s/%s", this.path, key)
}

func (this fileCache) Get(key string) ([]byte, error) {
	this.m.Lock()
	defer this.m.Unlock()

	path := this.getPath(key)
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return nil, CacheNotExistError
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	obj, err := this.d.Decode(file)
	if err != nil {
		return nil, err
	}
	if obj.Expire == 0 {
		return obj.Value, nil
	}

	if time.Now().Unix() > obj.Expire {
		// expire
		this.Delete(key)
		return nil, CacheNotExistError
	}

	return obj.Value, nil
}

func (this fileCache) Set(key string, value []byte) error {
	this.m.Lock()
	defer this.m.Unlock()
	file, err := os.Create(this.getPath(key))
	if err != nil {
		return err
	}
	defer file.Close()
	return this.d.Encode(file, &p{
		Value: value,
	})
}

func (this fileCache) SetWithExpire(key string, value []byte, expire time.Duration) error {
	this.m.Lock()
	defer this.m.Unlock()
	file, err := os.Create(this.getPath(key))
	if err != nil {
		return err
	}
	defer file.Close()

	expireTime := time.Now().Add(expire)
	return this.d.Encode(file, &p{
		Value:  value,
		Expire: expireTime.Unix(),
	})
}

func (this fileCache) Delete(key string) error {
	this.m.Lock()
	defer this.m.Unlock()
	return os.Remove(this.getPath(key))
}

func (this fileCache) Expire(key string, expire time.Duration) error {
	this.m.Lock()
	defer this.m.Unlock()

	path := this.getPath(key)
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return CacheNotExistError
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	obj, err := this.d.Decode(file)
	if err != nil {
		return err
	}

	expireTime := time.Now().Add(expire)
	return this.d.Encode(file, &p{
		Value:  obj.Value,
		Expire: expireTime.Unix(),
	})
}

func (this fileCache) Exisits(key string) error {
	this.m.Lock()
	defer this.m.Unlock()

	path := this.getPath(key)
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return CacheNotExistError
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	obj, err := this.d.Decode(file)
	if err != nil {
		return err
	}
	if obj.Expire == 0 {
		return nil
	}
	if time.Now().Unix() > obj.Expire {
		return CacheNotExistError
	}
	return nil
}

func (this *fileCache) Close() error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
type gobEncoder struct {
}

func (this gobEncoder) Encode(f *os.File, d *p) error {
	encoder := gob.NewEncoder(f)
	return encoder.Encode(d)
}

func (this gobEncoder) Decode(f *os.File) (*p, error) {
	decoder := gob.NewDecoder(f)
	var obj p
	err := decoder.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func NewCache(path string) (Cache, error) {
	return newCache(path, &gobEncoder{})
}

///////////////////////////////////////////////////////////////////////////////////
type jsonEncoder struct {
}

func (this jsonEncoder) Encode(f *os.File, d *p) error {
	encoder := json.NewEncoder(f)
	return encoder.Encode(d)
}

func (this jsonEncoder) Decode(f *os.File) (*p, error) {
	decoder := json.NewDecoder(f)
	var obj p

	if err := decoder.Decode(&obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func NewJsonCache(path string) (Cache, error) {
	return newCache(path, &jsonEncoder{})
}

///////////////////////////////////////////////////////////////////////////////////
func newCache(path string, d encoder) (Cache, error) {
	rand.Seed(time.Now().UnixNano())

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("path '%s' is NOT dir", path)
	}

	return &fileCache{
		path: path,
		d:    d,
	}, nil
}
