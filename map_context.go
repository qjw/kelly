// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type contextMap map[interface{}]interface{}

const (
	contextKey  = "context__"
	pathParamID = "pathVarible_"
)

type mapHttpContext struct {
	r *http.Request
}

func (c *mapHttpContext) Set(key, value interface{}) dataContext {
	datas := contextMustGet(c.r, contextKey).(contextMap)
	datas[key] = value
	return c
}

func (c mapHttpContext) Get(key interface{}) interface{} {
	datas := contextMustGet(c.r, contextKey).(contextMap)
	if data, ok := datas[key]; ok {
		return data
	} else {
		return nil
	}
}

func (c mapHttpContext) MustGet(key interface{}) interface{} {
	datas := contextMustGet(c.r, contextKey).(contextMap)
	if data, ok := datas[key]; ok {
		return data
	} else {
		panic(fmt.Errorf("can not get context value by '%v'", key))
	}
}

func newMapHttpContext(r *http.Request) dataContext {
	c := &mapHttpContext{
		r: r,
	}
	return c
}

type mapContext struct {
	data contextMap
}

func (c *mapContext) Set(key, value interface{}) dataContext {
	c.data[key] = value
	return c
}

func (c mapContext) Get(key interface{}) interface{} {
	if data, ok := c.data[key]; ok {
		return data
	} else {
		return nil
	}
}

func (c mapContext) MustGet(key interface{}) interface{} {
	if data, ok := c.data[key]; ok {
		return data
	} else {
		panic(fmt.Errorf("can not get context value by '%v'", key))
	}
}

func newMapContext() dataContext {
	c := &mapContext{}
	c.data = make(contextMap)
	return c
}

func mapContextFilter(_ http.ResponseWriter, r *http.Request, params httprouter.Params) *http.Request {
	contextMap := contextMap{
		pathParamID: params,
	}
	return contextSet(r, contextKey, contextMap)
}

func getPathParams(r *http.Request) httprouter.Params {
	datas := contextMustGet(r, contextKey).(contextMap)
	return datas[pathParamID].(httprouter.Params)
}
