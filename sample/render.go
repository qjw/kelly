// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"html/template"
	"net/http"
)

func InitRender(r kelly.Router) {
	render := r.Group("/render")

	render.GET("/t", func() kelly.HandlerFunc {
		data := `<form action="#" method="get">
<p>First {{ .First }}: <input type="text" name="fname" /></p>
<p>Last {{ .Last }}: <input type="text" name="lname" /></p>
<input type="submit" value="Submit" />
</form>`

		// 通过闭包预先编译好
		t := template.Must(template.New("t1").Parse(data))
		return func(c *kelly.Context) {
			c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
				"First": "Qiu",
				"Last": "King",
			})
		}
	}())
}
