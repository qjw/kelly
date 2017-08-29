// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package toolkits

import (
	"github.com/dchest/captcha"
	"github.com/qjw/kelly"
)

// 生成二维码的ID，不是实际的验证码数字，参数是验证码的长度
func GenerateCaptchaID(len int) string {
	return captcha.NewLen(4)
}

// request的url必须是 http(s)://host/path/{CaptchaID}.png
func ServerCaptcha(c *kelly.Context, width, height int) {
	handler := captcha.Server(width, height)
	handler.ServeHTTP(c, c.Request())
}

// 验证验证码
func VerifyCaptcha(captchaID, captchaCode string) bool {
	return captcha.VerifyString(captchaID, captchaCode)
}
