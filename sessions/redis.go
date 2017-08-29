// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2012 Brian "bojo" Jones. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sessions

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/securecookie"
	"gopkg.in/redis.v5"
	"net/http"
	"strings"
	"time"
	"github.com/qjw/kelly"
)

// Amount of time for cookies/redis keys to expire.
var sessionExpire = 86400 * 30

// SessionSerializer provides an interface hook for alternative serializers
type SessionSerializer interface {
	Deserialize(d []byte, session *SessionImp) error
	Serialize(session *SessionImp) ([]byte, error)
}

// JSONSerializer encode the session map to JSON.
type JSONSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s JSONSerializer) Serialize(session *SessionImp) ([]byte, error) {
	m := make(map[string]interface{}, len(session.Values))
	for k, v := range session.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
			fmt.Printf("redistore.JSONSerializer.serialize() Error: %v", err)
			return nil, err
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s JSONSerializer) Deserialize(d []byte, session *SessionImp) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(d, &m)
	if err != nil {
		fmt.Printf("redistore.JSONSerializer.deserialize() Error: %v", err)
		return err
	}
	for k, v := range m {
		session.Values[k] = v
	}
	return nil
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(session *SessionImp) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(session.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize back to map[interface{}]interface{}
func (s GobSerializer) Deserialize(d []byte, session *SessionImp) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&session.Values)
}

// RediStore stores sessions in a redis backend.
type RediStore struct {
	Pool          *redis.Client
	Codecs        []securecookie.Codec
	Options       *Options // default configuration
	DefaultMaxAge int      // default Redis TTL for a MaxAge == 0 session
	maxLength     int
	keyPrefix     string
	serializer    SessionSerializer
}

// SetMaxLength sets RediStore.maxLength if the `l` argument is greater or equal 0
// maxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new RediStore is 4096. Redis allows for max.
// value sizes of up to 512MB (http://redis.io/topics/data-types)
// Default: 4096,
func (s *RediStore) SetMaxLength(l int) {
	if l >= 0 {
		s.maxLength = l
	}
}

// SetKeyPrefix set the prefix
func (s *RediStore) SetKeyPrefix(p string) {
	s.keyPrefix = p
}

// SetSerializer sets the serializer
func (s *RediStore) SetSerializer(ss SessionSerializer) {
	s.serializer = ss
}

// SetMaxAge restricts the maximum age, in seconds, of the session record
// both in database and a browser. This is to change session storage configuration.
// If you want just to remove session use your session `s` object and change it's
// `Options.MaxAge` to -1, as specified in
//    http://godoc.org/github.com/gorilla/sessions#Options
//
// Default is the one provided by this package value - `sessionExpire`.
// Set it to 0 for no restriction.
// Because we use `MaxAge` also in SecureCookie crypting algorithm you should
// use this function to change `MaxAge` value.
func (s *RediStore) SetMaxAge(v int) {
	var c *securecookie.SecureCookie
	var ok bool
	s.Options.MaxAge = v
	for i := range s.Codecs {
		if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
			c.MaxAge(v)
		} else {
			fmt.Printf("Can't change MaxAge on codec %v\n", s.Codecs[i])
		}
	}
}

// NewRediStore returns a new RediStore.
// size: maximum number of idle connections.
func NewRediStore(redis *redis.Client, keyPairs ...[]byte) (*RediStore, error) {
	return NewRediStoreWithPool(redis, keyPairs...)
}

//func dialWithDB(network, address, password, DB string) (redis.Conn, error) {
//	c, err := dial(network, address, password)
//	if err != nil {
//		return nil, err
//	}
//	if _, err := c.Do("SELECT", DB); err != nil {
//		c.Close()
//		return nil, err
//	}
//	return c, err
//}

// NewRediStoreWithDB - like NewRedisStore but accepts `DB` parameter to select
// redis DB instead of using the default one ("0")
//func NewRediStoreWithDB(size int, network, address, password, DB string, keyPairs ...[]byte) (*RediStore, error) {
//	return NewRediStoreWithPool(&redis.Pool{
//		MaxIdle:     size,
//		IdleTimeout: 240 * time.Second,
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
//			_, err := c.Do("PING")
//			return err
//		},
//		Dial: func() (redis.Conn, error) {
//			return dialWithDB(network, address, password, DB)
//		},
//	}, keyPairs...)
//}

// NewRediStoreWithPool instantiates a RediStore with a *redis.Pool passed in.
func NewRediStoreWithPool(redis *redis.Client, keyPairs ...[]byte) (*RediStore, error) {
	rs := &RediStore{
		Pool:   redis,
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &Options{
			Path:   "/",
			MaxAge: sessionExpire,
		},
		DefaultMaxAge: 60 * 20, // 20 minutes seems like a reasonable default
		maxLength:     4096,
		keyPrefix:     "session_",
		serializer:    GobSerializer{},
	}
	_, err := rs.ping()
	return rs, err
}

// ping does an internal ping against a server to check if it is alive.
func (s *RediStore) ping() (bool, error) {
	pong, err := s.Pool.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
		return false, err
	}
	return (pong == "PONG"), nil
}

// Close closes the underlying *redis.Pool
//func (s *RediStore) Close() error {
//	return s.Pool.Close()
//}

// Get returns a session for the given name after adding it to the registry.
//
// See gorilla/sessions FilesystemStore.Get().
func (s *RediStore) Get(c *kelly.Context, name string) (*SessionImp, error) {
	registry := GetRegistry(c)
	return registry.Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See gorilla/sessions FilesystemStore.New().
func (s *RediStore) New(c *kelly.Context, name string) (*SessionImp, error) {
	var err error
	session := NewSession(s, name)
	// make a copy
	options := *s.Options
	session.Options = &options
	session.IsNew = true
	if cookies, errCookie := c.Request().Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, cookies.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err := s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *RediStore) Save(c *kelly.Context, session *SessionImp) error {
	// Marked for deletion.
	if session.Options.MaxAge < 0 {
		if err := s.delete(session); err != nil {
			return err
		}

		http.SetCookie(c, NewCookie(session.name, "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.name, session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(c, NewCookie(session.name, encoded, session.Options))
	}
	return nil
}

// Save adds a single session to the response.
func (s *RediStore) Delete(c *kelly.Context, name string) error {
	session, error := s.Get(c, name)
	if error != nil {
		return error
	}
	session.Options.MaxAge = -1
	return s.Save(c, session)
}

// save stores the session in redis.
func (s *RediStore) save(session *SessionImp) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	age := session.Options.MaxAge * 1000 * 1000 * 1000
	if age == 0 {
		age = s.DefaultMaxAge
	}
	err = s.Pool.Set(s.keyPrefix+session.ID, b, time.Duration(age)).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *RediStore) load(session *SessionImp) (bool, error) {
	data, err := s.Pool.Get(s.keyPrefix + session.ID).Bytes()
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil // no data was associated with this key
	}
	return true, s.serializer.Deserialize(data, session)
}

// delete removes keys from redis if MaxAge<0
func (s *RediStore) delete(session *SessionImp) error {
	return s.Pool.Del(s.keyPrefix + session.ID).Err()
}
