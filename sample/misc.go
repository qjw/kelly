// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/middleware"
	"net/http"
)

func InitMisc(r kelly.Router) {
	api := r.Group("/misc")

	api.GET("/basic", middleware.Basic("king", "qiu"), func(c *kelly.Context) {
		c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
			"code": "/basic auth ok",
		})
	})

	api.GET("/basic_func",
		middleware.BasicFunc(func(user, pwd string) bool {
			return (user == "king" && pwd == "qiu") ||
				(user == "qiu" && pwd == "king")
		}),
		func(c *kelly.Context) {
			c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
				"code": "/basic auth ok",
			})
		})
}
