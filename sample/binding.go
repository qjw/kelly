// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"net/http"
)

type BindPathObj struct {
	A string `json:"aaa,omitempty" binding:"required,max=32,min=1" error:"aerror"`
	B string `json:"bbb,omitempty" binding:"required,max=32,min=1" error:"berror"`
	C string `json:"ccc,omitempty" binding:"required,max=32,min=1" error:"cerror"`
}

type BindJsonObj struct {
	Obj1 BindPathObj `json:"obj,omitempty"`
	A    string      `json:"aaa,omitempty" binding:"required,max=32,min=1" error:"aerror"`
	B    string      `json:"bbb,omitempty" binding:"required,max=32,min=1" error:"berror"`
	C    string      `json:"ccc,omitempty" binding:"required,max=32,min=1" error:"cerror"`
}

func InitBinding(r kelly.Router) {
	api := r.Group("/bind")

	api.GET("/path/:aaa/:bbb/:ccc", func(c *kelly.Context) {
		var obj BindPathObj
		if err, _ := c.BindPath(&obj); err == nil {
			c.WriteJson(http.StatusOK, obj)
		} else {
			c.WriteString(http.StatusOK, "param err")
		}
	})
	api.GET("/path2/:aaa/:bbb/:ccc",
		kelly.BindPathMiddleware(&BindPathObj{}),
		func(c *kelly.Context) {
			c.WriteJson(http.StatusOK, c.GetBindPathParameter())
		})

	api.GET("/query", func(c *kelly.Context) {
		var obj BindPathObj
		if err, _ := c.Bind(&obj); err == nil {
			c.WriteJson(http.StatusOK, obj)
		} else {
			c.WriteString(http.StatusOK, "param err")
		}
	})
	api.GET("/query2",
		kelly.BindMiddleware(&BindPathObj{
			A:"dft a",
			B:"dft b",
		}),
		func(c *kelly.Context) {
			c.WriteJson(http.StatusOK, c.GetBindParameter())
		})

	api.GET("/form", func(c *kelly.Context) {
		data := `<html>
<head>
<script type="text/javascript" src="//cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
<script type="text/javascript">
$(function() {
	$("#link").click(function(){
		$.ajax({
            type: "POST",
            url: "/bind/json",
            contentType: "application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify({
            	aaa: "123",
            	bbb: "456",
            	ccc: "789",
            	obj: {
					aaa: "123",
					bbb: "456",
					ccc: "789",
            	},
            }),
            complete: function(msg){
            	alert(JSON.stringify(msg.responseText));
            },
        });
	})
})
</script>
</head>
<body>
	<form action="/bind/form" method="post">
	<p>AAA: <input type="text" name="aaa" /></p>
	<p>BBB: <input type="text" name="bbb" /></p>
	<p>CCC: <input type="text" name="ccc" /></p>
	<input type="submit" value="Submit" />
	</form>
	<br/>
	<!-- we will add our HTML content here -->
	<a href="#" id="link">Json按钮</a>
</body>
</html>`

		c.WriteHtml(http.StatusOK, data)
	})

	api.POST("/form", func(c *kelly.Context) {
		var obj BindPathObj
		if err, _ := c.Bind(&obj); err == nil {
			c.WriteJson(http.StatusOK, obj)
		} else {
			c.WriteString(http.StatusOK, "param err")
		}
	})

	api.POST("/json", func(c *kelly.Context) {
		var obj BindJsonObj
		if err, _ := c.Bind(&obj); err == nil {
			c.WriteJson(http.StatusOK, obj)
		} else {
			c.WriteString(http.StatusOK, "param err")
		}
	})

	api.GET("/form2", func(c *kelly.Context) {
		data := `<html>
<head>
<script type="text/javascript" src="//cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
<script type="text/javascript">
$(function() {
	$("#link").click(function(){
		$.ajax({
            type: "POST",
            url: "/bind/json2",
            contentType: "application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify({
            	aaa: "123",
            	bbb: "456",
            	ccc: "789",
            	obj: {
					aaa: "123",
					bbb: "456",
					ccc: "789",
            	},
            }),
            complete: function(msg){
            	alert(JSON.stringify(msg.responseText));
            },
        });
	})
})
</script>
</head>
<body>
	<form action="/bind/form2" method="post">
	<p>AAA: <input type="text" name="aaa" /></p>
	<p>BBB: <input type="text" name="bbb" /></p>
	<p>CCC: <input type="text" name="ccc" /></p>
	<input type="submit" value="Submit" />
	</form>
	<br/>
	<!-- we will add our HTML content here -->
	<a href="#" id="link">Json按钮</a>
</body>
</html>`

		c.WriteHtml(http.StatusOK, data)
	})

	api.POST("/form2",
		kelly.BindMiddleware(&BindPathObj{}),
		func(c *kelly.Context) {
			c.WriteJson(http.StatusOK, c.GetBindParameter())
		})
	api.POST("/json2",
		kelly.BindMiddleware(&BindJsonObj{}),
		func(c *kelly.Context) {
			c.WriteJson(http.StatusOK, c.GetBindParameter())
		})
}
