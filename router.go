// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"net/http"
)

type Router interface {
	GET(string, ...HandlerFunc) Router
	HEAD(string, ...HandlerFunc) Router
	OPTIONS(string, ...HandlerFunc) Router
	POST(string, ...HandlerFunc) Router
	PUT(string, ...HandlerFunc) Router
	PATCH(string, ...HandlerFunc) Router
	DELETE(string, ...HandlerFunc) Router

	ServeHTTP(http.ResponseWriter, *http.Request)

	// 新建子路由
	Group(string, ...HandlerFunc) Router

	// 动态插入中间件
	Use(...HandlerFunc) Router

	// 设置404处理句柄
	SetNotFoundHandle(HandlerFunc)
	// 设置405处理句柄
	SetMethodNotAllowedHandle(HandlerFunc)

	// 返回当前Router的绝对路径
	Path() string

	// 添加全局的 注解 函数。该router下面和子（孙）router下面的endpoint注册都会被触发
	GlobalAnnotation(handles ...AnnotationHandlerFunc) Router

	// 添加临时 注解 函数，只对使用返回的AnnotationRouter对象进行注册的endpoint有效
	Annotation(handles ...AnnotationHandlerFunc) AnnotationRouter

	// 用于支持设置context数据
	dataContext
}

func (rt *router) wrapParentHandle(n *negroni.Negroni) {
	if rt.parent != nil {
		rt.parent.wrapParentHandle(n)
		for _, v := range rt.parent.middlewares {
			n.UseFunc(wrapHandlerFunc(v))
		}
	}
}

func (rt *router) wrapHandle(handles ...HandlerFunc) httprouter.Handle {
	if DebugFlag && len(handles) == 0 {
		panic("invalid wrap handles")
	}

	tmpHandle := negroni.New()
	rt.wrapParentHandle(tmpHandle)

	for _, v := range rt.middlewares {
		tmpHandle.UseFunc(wrapHandlerFunc(v))
	}
	for _, v := range handles {
		tmpHandle.UseFunc(wrapHandlerFunc(v))
	}

	handle := func(c *Context) {
		tmpHandle.ServeHTTP(c, c.Request())
	}

	return func(wr http.ResponseWriter, r *http.Request, params httprouter.Params) {
		r = mapContextFilter(wr, r, params)
		handle(newContext(wr, r, nil))
	}
}

type endpoint struct {
	method             string
	path               string
	handles            []HandlerFunc
	endPointRegisterCB func()
}

func (this *endpoint) run() {
	if DebugFlag && this.endPointRegisterCB == nil {
		panic("invalid endpoint")
	}
	this.endPointRegisterCB()
	this.endPointRegisterCB = nil
}

type router struct {
	rt *httprouter.Router
	// 当前rouer路径
	path string
	// 绝对路径
	absolutePath string
	// 中间件
	middlewares []HandlerFunc
	// endpoint钩子函数
	epMiddlewares []AnnotationHandlerFunc
	// endpoint 次数
	endpointCnt int

	// 所有的子Group
	groups []*router

	// 父Group
	parent *router

	endpoints []*endpoint

	// 子group数量
	groupCnt int

	// 用于支持设置context数据
	dataContext

	// 被子类覆盖的方法，实现比较挫，再优化
	overiteInvokeAnnotation func(c *AnnotationContext)
}

func (rt *router) doBeforeRun() {
	for _, v := range rt.endpoints {
		v.run()
	}
	for _, v := range rt.groups {
		v.doBeforeRun()
	}
}

func (rt *router) Path() string {
	return rt.absolutePath
}

func (rt *router) SetNotFoundHandle(h HandlerFunc) {
	rt.rt.NotFound = h
}

func (rt *router) SetMethodNotAllowedHandle(h HandlerFunc) {
	rt.rt.MethodNotAllowed = h
}

func (rt *router) validatePath(path string) {
	if len(path) < 1 {
		panic(fmt.Errorf("invalid path %s", path))
	}
	if path == "/" {
		return
	}
	if path[0] != '/' || path[len(path)-1] == '/' {
		panic(fmt.Errorf("invalid path %s,must beginwith (NOT endwith) /", path))
	}
}

func (rt *router) validateParam(path string, handles ...HandlerFunc) {
	if len(handles) < 1 {
		panic(fmt.Errorf("must have one handle"))
	}
	rt.validatePath(path)
}

func (rt *router) methodImp(
	handle func(path string, handle httprouter.Handle),
	method string,
	path string,
	handles ...HandlerFunc) Router {

	rt.validateParam(path, handles...)
	f := rt.overiteInvokeAnnotation

	// 增加计数
	rt.endpointCnt += 1
	rt.endpoints = append(rt.endpoints, &endpoint{
		method:  method,
		path:    path,
		handles: handles,
		endPointRegisterCB: func() {
			handle(rt.absolutePath+path, rt.wrapHandle(handles...))
			f(&AnnotationContext{
				r:       rt,
				method:  method,
				path:    path,
				handles: handles,
			})
		},
	})
	return rt
}

func (rt *router) GET(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.GET, "GET", path, handles...)
}

func (rt *router) HEAD(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.HEAD, "HEAD", path, handles...)
}

func (rt *router) OPTIONS(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.OPTIONS, "OPTIONS", path, handles...)
}

func (rt *router) POST(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.POST, "POST", path, handles...)
}

func (rt *router) PUT(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.PUT, "PUT", path, handles...)
}

func (rt *router) PATCH(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.PATCH, "PATCH", path, handles...)
}

func (rt *router) DELETE(path string, handles ...HandlerFunc) Router {
	return rt.methodImp(rt.rt.DELETE, "DELETE", path, handles...)
}

func (rt *router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rt.rt.ServeHTTP(rw, r)
}

func (rt *router) Group(path string, handlers ...HandlerFunc) Router {
	rt.validatePath(path)
	newRt := &router{
		rt:           rt.rt,
		path:         path,
		absolutePath: rt.absolutePath + path,
		middlewares:  handlers,
		dataContext:  newMapContext(),
		parent:       rt,
	}
	newRt.overiteInvokeAnnotation = newRt.invokeAnnotation

	// 拷贝context
	ctx := rt.dataContext.(*mapContext)
	for k, v := range ctx.data {
		newRt.dataContext.Set(k, v)
	}

	// 增加计数
	rt.groupCnt += 1
	rt.groups = append(rt.groups, newRt)
	return newRt
}

func (rt *router) Use(handlers ...HandlerFunc) Router {
	if DebugFlag {
		if len(handlers) == 0 {
			panic("invalid handlers")
		}
		if len(rt.groups) != rt.groupCnt {
			panic("invalid router")
		}
	}
	for _, v := range handlers {
		rt.middlewares = append(rt.middlewares, v)
	}
	return rt
}

func (rt *router) GlobalAnnotation(handles ...AnnotationHandlerFunc) (r Router) {
	r = rt
	if DebugFlag && len(handles) == 0 {
		panic("EndPointHandlerFunc at list one")
		return
	}

	if len(rt.epMiddlewares) == 0 {
		rt.epMiddlewares = make([]AnnotationHandlerFunc, len(handles))
		copy(rt.epMiddlewares, handles)
	} else {
		for _, item := range handles {
			rt.epMiddlewares = append(rt.epMiddlewares, item)
		}
	}
	return
}

// 添加临时 注解 函数，只对使用返回的AnnotationRouter对象进行注册的endpoint有效
func (rt *router) Annotation(handles ...AnnotationHandlerFunc) AnnotationRouter {
	return newAnnotationRouter(rt, handles...)
}

func (rt *router) invokeParentAnnotation(c *AnnotationContext) {
	if rt.parent != nil {
		rt.parent.invokeParentAnnotation(c)
		for _, item := range rt.parent.epMiddlewares {
			item(c)
		}
	}
}

// 触发 注解 函数 不公开
func (rt *router) invokeAnnotation(c *AnnotationContext) {
	rt.invokeParentAnnotation(c)

	// 执行全局的ep 过滤器
	for _, item := range rt.epMiddlewares {
		item(c)
	}
}

func newRouterImp(handlers ...HandlerFunc) *router {
	httpRt := httprouter.New()
	rt := &router{
		rt:           httpRt,
		path:         "",
		absolutePath: "",
		dataContext:  newMapContext(),
	}
	rt.overiteInvokeAnnotation = rt.invokeAnnotation

	if len(handlers) > 0 {
		rt.middlewares = make([]HandlerFunc, len(handlers))
		copy(rt.middlewares, handlers)
	}

	return rt
}
