// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"context"
	"fmt"
	"net/http"
)

// 包装http.Request的context操作
type httpContext struct {
	r *http.Request
}

func (c *httpContext) Set(key, value interface{}) dataContext {
	c.r = contextSet(c.r, key, value)
	return c
}

func (c *httpContext) Get(key interface{}) interface{} {
	return contextGet(c.r, key)
}

func (c *httpContext) MustGet(key interface{}) interface{} {
	return contextMustGet(c.r, key)
}

func newHttpContext(r *http.Request) dataContext {
	c := &httpContext{
		r: r,
	}
	return c
}

func contextSet(r *http.Request, key, value interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

func contextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func contextMustGet(r *http.Request, key interface{}) interface{} {
	v := r.Context().Value(key)
	if v == nil {
		panic(fmt.Errorf("get context value fail by '%v'", key))
	}
	return v
}
