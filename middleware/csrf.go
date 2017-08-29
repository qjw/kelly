// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// 参考 https://github.com/tommy351/gin-csrf

package middleware

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"github.com/dchest/uniuri"
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/sessions"
	"io"
	"net/http"
)

const (
	csrfSalt   = "csrfSalt"
	csrfToken  = "csrfToken"
)

var (
	defaultCookieKey                = "_crsf"
	defaultStore     sessions.Store = nil
	goptions *CsrfConfig = nil
)

var defaultIgnoreMethods = []string{"GET", "HEAD", "OPTIONS"}

var defaultErrorFunc = func(c *kelly.Context) {
	c.WriteJson(http.StatusForbidden, kelly.H{
		"code":    http.StatusForbidden,
		"message": "CSRF token mismatch",
	})
}

var defaultTokenGetter = func(c *kelly.Context) string {
	r := c.Request()

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	}

	return ""
}

// CsrfConfig stores configurations for a CSRF middleware.
type CsrfConfig struct {
	Secret        []byte
	IgnoreMethods []string
	ErrorFunc     kelly.HandlerFunc
	TokenGetter   func(c *kelly.Context) string
}

func tokenize(secret []byte, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt+"-"+string(secret))
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash
}

func inArray(arr []string, value string) bool {
	inarr := false

	for _, v := range arr {
		if v == value {
			inarr = true
			break
		}
	}

	return inarr
}

func InitCsrf(options CsrfConfig){
	if defaultStore == nil {
		store := sessions.NewCookieStore(options.Secret)
		defaultStore = store
		gob.Register([]interface{}{})
	}else{
		panic("init csrf yet")
	}

	if options.IgnoreMethods == nil {
		options.IgnoreMethods = defaultIgnoreMethods
	}

	if options.ErrorFunc == nil {
		options.ErrorFunc = defaultErrorFunc
	}

	if options.TokenGetter == nil {
		options.TokenGetter = defaultTokenGetter
	}

	goptions = &options
}

// Middleware validates CSRF token.
func Csrf() kelly.HandlerFunc {
	if defaultStore == nil{
		panic("not init csrf yet")
	}

	return func(c *kelly.Context) {
		session, _ := defaultStore.Get(c, defaultCookieKey)

		// 是否在忽略的http方法中
		if inArray(goptions.IgnoreMethods, c.Request().Method) {
			c.InvokeNext()
			return
		}

		// 找到cookie中的盐
		var salt string
		if s, ok := session.Get(csrfSalt).(string); !ok || len(s) == 0 {
			goptions.ErrorFunc(c)
			return
		} else {
			salt = s
		}
		session.Delete(csrfSalt)
		session.Save()

		token := goptions.TokenGetter(c)

		// 通过标注通道生成的token和cookie中盐生成的token不一致
		if len(token) < 1 || tokenize(goptions.Secret, salt) != token {
			goptions.ErrorFunc(c)
		} else {
			c.InvokeNext()
		}
	}
}

// GetToken returns a CSRF token.
func GetCsrfToken(c *kelly.Context) string {
	if defaultStore == nil{
		panic("not init csrf yet")
	}

	session, _ := defaultStore.Get(c, defaultCookieKey)
	secret := goptions.Secret

	// 如果已经设置了token，就直接返回
	if t := c.Get(csrfToken); t != nil {
		return t.(string)
	}

	// 生成token
	salt := uniuri.New()
	token := tokenize(secret, salt)

	// 将salt保存到cookie中
	session.Set(csrfSalt, salt)
	session.Save()

	// 设置context，下次复用
	c.Set(csrfToken, token)

	return token
}
