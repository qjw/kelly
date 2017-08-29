// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"net/http"
)

func InitStatic(r kelly.Router) {
	r.GET("/static/*path", kelly.Static(&kelly.StaticConfig{
		Dir:        http.Dir("/usr/share/nginx/html/"),
		Indexfiles: []string{"index.html"},
	}))

	r.GET("/static1/*path", kelly.Static(&kelly.StaticConfig{
		Dir:           http.Dir("/tmp"),
		EnableListDir: true,
	}))
}
