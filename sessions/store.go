// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package sessions

import (
	"github.com/qjw/kelly"
)

type Store interface {
	// Get should return a cached session.
	Get(c *kelly.Context, name string) (*SessionImp, error)

	// New should create and return a new session.
	//
	// Note that New should never return a nil session, even in the case of
	// an error if using the Registry infrastructure to cache the session.
	New(c *kelly.Context, name string) (*SessionImp, error)

	// Save should persist session to the underlying store implementation.
	Save(c *kelly.Context, s *SessionImp) error

	// 删除cookie
	Delete(c *kelly.Context, name string) error
}
