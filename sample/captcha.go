// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/toolkits"
	"html/template"
	"net/http"
)

func InitCaptcha(r kelly.Router) {

	r.GET("/captcha", func() kelly.HandlerFunc {
		data := `<form action="/captcha" method="post">
<p>验证码: <input type="text" name="captchaStr" /></p>
<input type="hidden" name="captchaID" value="{{ .CaptchaID }}" /></p>
<img id="captcha-img" width="104" height="36" src="/captcha/image/{{ .CaptchaID }}.png" />
<input type="submit" value="提交" />
</form>
<script type="text/javascript" src="//cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
<script type="text/javascript">
$(function() {
	$("#captcha-img").click(function() {
	  var captcha_url = $(this).attr("src").split("?")[0];
	  captcha_url += "?reload=1&timestamp=" + new Date().getTime()
	  $(this).attr("src",  captcha_url);
	});
})
</script>`

		// 通过闭包预先编译好
		t := template.Must(template.New("t1").Parse(data))
		return func(c *kelly.Context) {
			captchaID := toolkits.GenerateCaptchaID(4)
			c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
				"CaptchaID": captchaID,
			})
		}
	}())

	r.GET("/captcha/image/:id", func(c *kelly.Context) {
		toolkits.ServerCaptcha(c, 104, 36)
	})

	r.POST("/captcha", func(c *kelly.Context) {
		if toolkits.VerifyCaptcha(
			c.MustGetFormVarible("captchaID"),
			c.MustGetFormVarible("captchaStr"),
		) {
			c.ResponseStatusOK()
		} else {
			c.ResponseStatusForbidden(nil)
		}
	})
}
