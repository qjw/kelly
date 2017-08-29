// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"net/http"
)

func InitParam(r kelly.Router) {
	r.GET("/path/:name", func(c *kelly.Context) {
		c.WriteIndentedJson(http.StatusOK, kelly.H{
			"code":  "/path",
			"path":  c.MustGetPathVarible("name"), // 获取path参数
			"query": c.GetDefaultQueryVarible("abc", "def"), // 获取query参数
		})
	})

	// -----------------------------------------------------------------
	r.GET("/form", func(c *kelly.Context) {
		data := `<form action="/form" method="post">
<p>First name: <input type="text" name="fname" /></p>
<p>Last name: <input type="text" name="lname" /></p>
<input type="submit" value="Submit" />
</form>`
		c.WriteHtml(http.StatusOK, data) // 返回html
	})
	r.POST("/form", func(c *kelly.Context) {
		c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
			"code":        "/form",
			"first name":  c.GetDefaultFormVarible("fname", "fname"), // 获取form参数
			"second name": c.GetDefaultFormVarible("lname", "lname"),
		})
	})
}
