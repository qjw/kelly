// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

// endpint Context，用于记录每个请求的信息
type AnnotationContext struct {
	r       Router
	method  string
	path    string
	handles []HandlerFunc
}

func (e AnnotationContext) R() Router {
	return e.r
}

func (e AnnotationContext) Method() string {
	return e.method
}

func (e AnnotationContext) Path() string {
	return e.path
}

func (e AnnotationContext) HandleCnt() int {
	return len(e.handles)
}

func (e AnnotationContext) Handle(index int) HandlerFunc {
	if index < 0 || index >= e.HandleCnt() {
		panic("invalid index")
	}
	return e.handles[index]
}

func (e AnnotationContext) HandlerFunc() HandlerFunc {
	if len(e.handles) < 1 {
		panic("invalid endPoint context")
	}
	return e.handles[len(e.handles)-1]
}

//--------------------------------------------------------------------
type AnnotationHandlerFunc func(c *AnnotationContext)

type AnnotationRouter interface {
	Router
}

//--------------------------------------------------------------------

func newAnnotationRouter(r *router, handles ...AnnotationHandlerFunc) AnnotationRouter {
	return &annotationRouter{
		router:      r,
		middlewares: handles,
	}
}

type annotationRouter struct {
	*router
	// endpoint钩子函数
	middlewares []AnnotationHandlerFunc
}

// 覆盖默认的
func (r *annotationRouter) Annotation(handles ...AnnotationHandlerFunc) AnnotationRouter {
	for _, v := range handles {
		r.middlewares = append(r.middlewares, v)
	}
	return r
}

func (r *annotationRouter) doMethod(
	f func(path string, handles ...HandlerFunc) Router,
	path string,
	handles ...HandlerFunc,
) Router {
	old := r.router.overiteInvokeAnnotation
	r.router.overiteInvokeAnnotation = r.invokeAnnotation
	defer func() {
		r.router.overiteInvokeAnnotation = old
	}()
	return f(path, handles...)
}

func (rt *annotationRouter) GET(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.GET, path, handles...)
}

func (rt *annotationRouter) HEAD(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.HEAD, path, handles...)
}

func (rt *annotationRouter) OPTIONS(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.OPTIONS, path, handles...)
}

func (rt *annotationRouter) POST(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.POST, path, handles...)
}

func (rt *annotationRouter) PUT(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.PUT, path, handles...)
}

func (rt *annotationRouter) PATCH(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.PATCH, path, handles...)
}

func (rt *annotationRouter) DELETE(path string, handles ...HandlerFunc) Router {
	return rt.doMethod(rt.router.DELETE, path, handles...)
}

// 覆盖默认的
func (rt *annotationRouter) invokeAnnotation(c *AnnotationContext) {
	// 执行全局的ep 过滤器
	rt.router.invokeAnnotation(c)

	// 执行本地的
	for _, item := range rt.middlewares {
		item(c)
	}
}
