// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

type StaticConfig struct {
	Dir           http.FileSystem
	EnableListDir bool
	Indexfiles    []string
}

// 根据一个目录生成一个HandlerFunc处理文件请求，在绑定Path时，必须使用下面的规则
// r.GET("/static/*path", kelly.Static(http.Dir("/tmp")))
// 在内部依赖于名称为path的路径变量
// 若将*path改成:path，将只能访问根目录的文件，无法嵌套
func Static(config *StaticConfig) HandlerFunc {
	staticTemp := `<pre>
{{ range $key, $value := . }}
	<a href="{{ $value.Url }}" style="color: {{ $value.Color }};">{{ $value.Name }}</a>
{{ end }}
</pre>`

	t := template.Must(template.New("staticTemp").Parse(staticTemp))
	if len(config.Indexfiles) > 0 {
		config.EnableListDir = false
		for _, v := range config.Indexfiles {
			if len(v) < 1 {
				panic(fmt.Errorf("invalid index file"))
			} else if strings.ContainsAny(v, "/") {
				panic(fmt.Errorf("invalid index file %s", v))
			}
		}
	}

	return func(c *Context) {
		r := c.Request()
		if r.Method != "GET" && r.Method != "HEAD" {
			c.WriteString(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		file := c.MustGetPathVarible("path")
		f, err := config.Dir.Open(file)
		if err != nil {
			c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		// 处理文件
		if fi.IsDir() {
			if len(config.Indexfiles) > 0 {
				serverIndex(config, file, c)
			}
			if config.EnableListDir {
				listDir(f, t, c)
			}
			return
		}

		http.ServeContent(c, r, file, fi.ModTime(), f)
	}
}

func serverIndex(config *StaticConfig, file string, c *Context) {
	var target = ""
	for _, v := range config.Indexfiles {
		file = path.Join(file, v)
		f, err := config.Dir.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil || fi.IsDir() {
			continue
		}

		target = v
		break
	}

	if len(target) == 0 {
		c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	} else {
		c.Redirect(http.StatusMovedPermanently, target)
	}
}

// 参考 https://github.com/labstack/echo/blob/master/middleware/static.go
func listDir(d http.File, t *template.Template, c *Context) {
	dirs, err := d.Readdir(-1)
	if err != nil {
		c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	data := []map[string]string{}
	for _, d := range dirs {
		name := d.Name()
		color := "#212121"
		if d.IsDir() {
			color = "#e91e63"
			name += "/"
		}

		data = append(data, map[string]string{
			"Name":  name,
			"Color": color,
			"Url":   name,
		})
	}

	c.WriteTemplateHtml(http.StatusOK, t, data)
}
