// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Context struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	r *http.Request

	// 下一个处理逻辑，用于middleware
	next http.HandlerFunc

	// 用于支持设置context数据
	dataContext
	// render
	renderOp
	// request
	request
	// binder
	binder
}

func (c *Context) Request() *http.Request {
	return c.r
}

func (c *Context) SetBody(body []byte) {
	c.r.Body.Close()
	c.r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

// 重置request，用于支持context库补充额外value
func (c *Context) SetRequest(r *http.Request) {
	c.r = r
	// 无须更新dataContext，因为后者不依赖于context库，而是map实现
}

func (c *Context) InvokeNext() {
	if c.next != nil {
		c.next.ServeHTTP(c, c.Request())
	} else {
		panic("invalid invoke next")
	}
}

func newContext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) *Context {
	c := &Context{
		ResponseWriter: w,
		Hijacker:       w.(http.Hijacker),
		Flusher:        w.(http.Flusher),
		r:              r,
		next:           next,
		dataContext:    newMapHttpContext(r),
		request:        newRequest(r),
	}
	c.renderOp = newRender(c)
	c.binder = newBinder(c)
	return c
}
