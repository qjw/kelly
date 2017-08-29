// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package sessions

import (
	"encoding/base32"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"github.com/qjw/kelly"
)

var fileMutex sync.RWMutex

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(path string, keyPairs ...[]byte) *FilesystemStore {
	if path == "" {
		path = os.TempDir()
	}
	fs := &FilesystemStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		path: path,
	}

	fs.MaxAge(fs.Options.MaxAge)
	return fs
}

// FilesystemStore stores sessions in the filesystem.
//
// It also serves as a reference for custom stores.
//
// This store is still experimental and not well tested. Feedback is welcome.
type FilesystemStore struct {
	Codecs  []securecookie.Codec
	Options *Options // default configuration
	path    string
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new FilesystemStore is 4096.
func (s *FilesystemStore) MaxLength(l int) {
	for _, c := range s.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *FilesystemStore) Get(c *kelly.Context, name string) (*SessionImp, error) {
	registry := GetRegistry(c)
	return registry.Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *FilesystemStore) New(c *kelly.Context, name string) (*SessionImp, error) {
	session := NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if cookies, errCookie := c.Request().Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, cookies.Value, &session.ID, s.Codecs...)
		if err == nil {
			err = s.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save adds a single session to the response.
//
// If the Options.MaxAge of the session is <= 0 then the session file will be
// deleted from the store path. With this process it enforces the properly
// session cookie handling so no need to trust in the cookie management in the
// web browser.
func (s *FilesystemStore) Save(c *kelly.Context, session *SessionImp) error {
	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		if err := s.erase(session); err != nil {
			return err
		}
		http.SetCookie(c, NewCookie(session.name, "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}
	if err := s.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.name, session.ID,
		s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(c, NewCookie(session.name, encoded, session.Options))
	return nil
}

// Save adds a single session to the response.
func (s *FilesystemStore) Delete(c *kelly.Context, name string) error {
	session, error := s.Get(c, name)
	if error != nil {
		return error
	}
	session.Options.MaxAge = -1
	return s.Save(c, session)
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *FilesystemStore) MaxAge(age int) {
	s.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// save writes encoded session.Values to a file.
func (s *FilesystemStore) save(session *SessionImp) error {
	encoded, err := securecookie.EncodeMulti(session.name, session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return ioutil.WriteFile(filename, []byte(encoded), 0600)
}

// load reads a file and decodes its content into session.Values.
func (s *FilesystemStore) load(session *SessionImp) error {
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.RLock()
	defer fileMutex.RUnlock()
	fdata, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = securecookie.DecodeMulti(session.name, string(fdata),
		&session.Values, s.Codecs...); err != nil {
		return err
	}
	return nil
}

// delete session file
func (s *FilesystemStore) erase(session *SessionImp) error {
	filename := filepath.Join(s.path, "session_"+session.ID)

	fileMutex.RLock()
	defer fileMutex.RUnlock()

	err := os.Remove(filename)
	return err
}
