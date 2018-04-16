package store

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/qjw/kelly"
)

type CookieStore struct {
	Codecs []securecookie.Codec
}

func (this CookieStore) New(c *kelly.Context, options *Options, name string) (Session, error) {
	cs := newClientSession(&this, c, options, name)

	var err error
	if cookie, errCookie := c.Request().Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(
			name,
			cookie.Value,
			&cs.Values,
			this.Codecs...,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
	return cs, nil
}

func (this CookieStore) Save(c *kelly.Context, s Session) error {
	cs := s.(*clientSession)
	encoded, err := securecookie.EncodeMulti(cs.Name, cs.Values, this.Codecs...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	http.SetCookie(c, newCookie(cs.Name, encoded, cs.Options))
	return nil
}

func (this CookieStore) Delete(c *kelly.Context, s Session) error {
	cs := s.(*clientSession)
	cs.Options.MaxAge = newInt(-1)
	cs.Written = true
	return this.Save(c, cs)
}

///////////////////////////////////////////////////////////////////////////////////

func NewCookieStore(keyPairs ...[]byte) Store {
	cs := &CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
	}
	return cs
}
