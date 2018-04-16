package store

import (
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/sessions/cache"
)

type serializer interface {
	Deserialize(d []byte, session *serverSession) error
	Serialize(session *serverSession) ([]byte, error)
}

type jsonSerializer struct {
}

// unsupported type: map[interface {}]interface {}
func (this jsonSerializer) Deserialize(d []byte, session *serverSession) error {
	return json.Unmarshal(d, &session.Values)
}

func (this jsonSerializer) Serialize(session *serverSession) ([]byte, error) {
	return json.Marshal(&session.Values)
}

///////////////////////////////////////////////////////////////////////////////////

type ServerStore struct {
	Codecs     []securecookie.Codec
	Cache      cache.Cache
	Serializer serializer
	MaxLength  int
}

func (this ServerStore) New(c *kelly.Context, options *Options, name string) (Session, error) {
	ss := newServerSession(&this, c, options, name)

	var err error
	if cookie, errCookie := c.Request().Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(
			name,
			cookie.Value,
			&ss.ID,
			this.Codecs...,
		)
		if err == nil {
			this.load(ss)
		}
	}
	return ss, nil
}

func (this ServerStore) Save(c *kelly.Context, s Session) error {
	ss := s.(*serverSession)
	if ss.Options.MaxAge != nil && *ss.Options.MaxAge < 0 {
		if err := this.del(ss); err != nil {
			log.Printf("server store del fail %v", err)
			return err
		}

		http.SetCookie(c, newCookie(ss.Name, "", ss.Options))
	} else {
		if ss.ID == "" {
			ss.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32),
			), "=")
			log.Printf("server store new session id %v", ss.ID)
		}
		if err := this.save(ss); err != nil {
			fmt.Println(err)
			return err
		}
		encoded, err := securecookie.EncodeMulti(
			ss.Name,
			ss.ID,
			this.Codecs...,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
		http.SetCookie(c, newCookie(ss.Name, encoded, ss.Options))
	}
	return nil
}

func (this ServerStore) Delete(c *kelly.Context, s Session) error {
	ss := s.(*serverSession)
	ss.Options.MaxAge = newInt(-1)
	ss.Written = true
	return this.Save(c, ss)
}

func (this ServerStore) save(ss *serverSession) error {
	b, err := this.Serializer.Serialize(ss)
	if err != nil {
		return err
	}
	if this.MaxLength != 0 && len(b) > this.MaxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	if ss.Options.MaxAge == nil || *ss.Options.MaxAge == 0 {
		err = this.Cache.Set(ss.Key(), b)
	} else {
		err = this.Cache.SetWithExpire(
			ss.Key(),
			b,
			time.Second*time.Duration(*ss.Options.MaxAge),
		)
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (this *ServerStore) load(session *serverSession) (bool, error) {
	data, err := this.Cache.Get(session.Key())
	if err != nil || len(data) < 1 {
		return false, err
	}
	return true, this.Serializer.Deserialize([]byte(data), session)
}

func (this *ServerStore) del(ss *serverSession) error {
	return this.Cache.Delete(ss.Key())
}

///////////////////////////////////////////////////////////////////////////////////

func NewServerStore(c cache.Cache, keyPairs ...[]byte) Store {
	cs := &ServerStore{
		Codecs:     securecookie.CodecsFromPairs(keyPairs...),
		Serializer: &jsonSerializer{},
		MaxLength:  4096,
		Cache:      c,
	}
	return cs
}
