// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"fmt"
	"github.com/qjw/kelly"
	"io"
	"net/http"
	"os"
)

func InitUpload(r kelly.Router) {
	r.GET("/upload", func(c *kelly.Context) {
		data := `<form enctype="multipart/form-data" action="/upload" method="post">
  <input type="file" name="file1" />
  <input type="file" name="file2" />
  <input type="submit" value="upload" />
</form>`
		c.WriteHtml(http.StatusOK, data) // 返回html
	})

	r.POST("/upload", func(c *kelly.Context) {
		c.ParseMultipartForm()

		file, handler := c.MustGetFileVarible("file1")
		defer file.Close()
		f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		file2, handler2 := c.MustGetFileVarible("file2")
		defer file2.Close()
		f2, err := os.OpenFile("./"+handler2.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f2.Close()
		io.Copy(f2, file2)

		c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
			"code":        "/upload",
			"first name":  handler.Filename, // 获取form参数
			"second name": handler2.Filename,
		})
	})
}
