// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/middleware"
	"html/template"
	"net/http"
)

func InitCsrf(r kelly.Router) {
	middleware.InitCsrf(middleware.CsrfConfig{
		Secret: []byte("fasdffasdfas"),
	})

	api := r.Group("/csrf",
		middleware.Csrf(),
	)

	api.GET("/ok", func() kelly.HandlerFunc {
		data := `<form action="/csrf//ok" method="post">
<p>First {{ .First }}: <input type="text" name="fname" /></p>
<p><input type="hidden" name="_csrf" value="{{ .Token }}"> </p>
<input type="submit" value="Submit" />
</form>`

		// 通过闭包预先编译好
		t := template.Must(template.New("ok").Parse(data))
		return func(c *kelly.Context) {
			c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
				"First": "Qiu",
				"Token": middleware.GetCsrfToken(c),
			})
		}
	}())

	api.POST("/ok", func(c *kelly.Context) {
		c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
			"code": "/csrf ok",
		})
	})

	api.GET("/fail", func() kelly.HandlerFunc {
		data := `<form action="/csrf/fail?_csrf={{ .Token }}" method="post">
<p>First {{ .First }}: <input type="text" name="fname" /></p>
<input type="submit" value="Submit" />
</form>`

		// 通过闭包预先编译好
		t := template.Must(template.New("fail").Parse(data))
		return func(c *kelly.Context) {
			c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
				"First": "Qiu",
				"Token": middleware.GetCsrfToken(c),
			})
		}
	}())

	api.POST("/fail", func(c *kelly.Context) {
		c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
			"code": "/csrf fail",
		})
	})

	api.GET("/ajax", func() kelly.HandlerFunc {
		data := `<html>
<head>
<script type="text/javascript" src="//cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
<script type="text/javascript">
$(function() {
	$(document).ajaxSend(function(elm, xhr, s){
		if (s.type == "POST") {
			xhr.setRequestHeader('x-csrf-token', {{ .Token }});
		}
	});

	$("#link").click(function(){
		$.ajax({
            type: "POST",
            url: "/csrf/fail",
            complete: function(msg){
            	alert(JSON.stringify(msg.responseText));
            },
        });
	})
})
</script>
</head>
<body>
  <!-- we will add our HTML content here -->
  <a href="#" id="link">Link</a>
</body>
</html>`

		// 通过闭包预先编译好
		t := template.Must(template.New("ajax").Parse(data))
		return func(c *kelly.Context) {
			c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
				"Token": middleware.GetCsrfToken(c),
			})
		}
	}())
}
