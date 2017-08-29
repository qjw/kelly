// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"github.com/urfave/negroni"
	"net/http"
)

type HandlerFunc func(c *Context)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(newContext(w, r, nil))
}

type Kelly interface {
	Router
	Run(addr ...string)
}

type kellyImp struct {
	*router
	n *negroni.Negroni
}

func (k *kellyImp) Run(addr ...string) {
	if k.n == nil {
		panic("invalid kelly")
	}
	k.router.doBeforeRun()
	k.n.Run(addr...)
}

func newImp(n *negroni.Negroni, handlers ...HandlerFunc) Kelly {
	rt := newRouterImp(handlers...)
	ky := &kellyImp{
		router: rt,
		n:      n,
	}
	ky.n = n

	n.UseHandler(rt)
	return ky
}

func EnableDebug(debug bool) {
	DebugFlag = debug
}

func NewClassic(handlers ...HandlerFunc) Kelly {
	return newImp(negroni.Classic(), handlers...)
}

func New(handlers ...HandlerFunc) Kelly {
	return newImp(negroni.New(negroni.NewRecovery()), handlers...)
}
