// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"net/http"
)

func Middleware(title, ver string, ok bool) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		if ok {
			// 设置header
			c.SetHeader(title,ver)
			// 设置context参数
			c.Set(title, ver)

			// 调用下一个handle
			c.InvokeNext()
		} else {
			// 不调用@InvokeNext 意味着中断执行流程
			c.WriteIndentedJson(http.StatusForbidden, kelly.H{
				"code": http.StatusForbidden,
			})
		}
	}
}

func InitGroupMiddleware(r kelly.Router) {
	// 新建一个子router，并注入一个middleware
	ar := r.Group(
		"/aaa",
		Middleware("v2", "v2", true),
	)
	ar.GET("/", func(c *kelly.Context) {
		c.WriteJson(http.StatusOK, kelly.H{ // 返回json（紧凑格式）
			"code": "/aaa",
		})
	})
	ar.GET("/a/b/c/d", func(c *kelly.Context) {
		c.Redirect(http.StatusMovedPermanently, "/aaa/bbb/a") // 重定向
	})
	ar.GET("/b/*path", func(c *kelly.Context) {
		c.WriteJson(http.StatusOK, kelly.H{ // 返回json（紧凑格式）
			"code": c.MustGetPathVarible("path"),
		})
	})
	ar.GET("/c/:path/d", func(c *kelly.Context) {
		c.WriteJson(http.StatusOK, kelly.H{ // 返回json（紧凑格式）
			"code": c.MustGetPathVarible("path"),
		})
	})
	ar.GET("/d/:path/:path2", func(c *kelly.Context) {
		c.WriteJson(http.StatusOK, kelly.H{ // 返回json（紧凑格式）
			"code": c.MustGetPathVarible("path"),
			"code2": c.MustGetPathVarible("path2"),
		})
	})

	// 新建一个子router，并注入一个middleware
	sar := ar.Group(
		"/bbb",
		Middleware("v3", "v3", true),
	)
	sar.GET("/", func(c *kelly.Context) {
		c.WriteXml(http.StatusOK, kelly.H{  // 返回XML
			"code": "/aaa/bbb",
		})
	})
	sar.GET("/a", func(c *kelly.Context) {
		c.WriteString(http.StatusOK, "test %d %d", 123, 456) // 返回普通文本
	})
}
