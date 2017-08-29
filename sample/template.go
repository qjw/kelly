// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/toolkits"
	"github.com/qjw/kelly/middleware"
)

func InitTemplate(r kelly.Router) {
	mng := toolkits.NewTemplateManage(ProjectRoot)
	// mng := toolkits.NewGoTemplateManage(ProjectRoot)
	r.GET("/template", func() kelly.HandlerFunc {
		temp := mng.MustGetTemplate("template/index.html")
		// temp := mng.MustGetTemplate("template/gotemplate.html")

		return func(c *kelly.Context) {
			temp.Render(c, kelly.H{
				"Body": "Kelly",
			})
		}
	}())

	toolkits.InitTemplateMiddleware(mng)
	r.GET("/template2",
		middleware.Gzip(middleware.BestSpeed, middleware.GzipMethod),
		toolkits.TemplateMiddleware("template/index.html"),
		// toolkits.TemplateMiddleware("template/gotemplate.html"),

		func(c *kelly.Context) {
			toolkits.CurrentTemplate(c).Render(c, kelly.H{
				"Body": "Kelly",
			})
		})
}
