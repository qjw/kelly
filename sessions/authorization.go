// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package sessions

import (
	"encoding/gob"
	"log"
	"net/http"
	"reflect"
	"github.com/qjw/kelly"
)

const (
	AUTH_SESSION_NAME = "session"
	AUTH_SESSION_KEY  = "_user_"
)

type authInstance struct {
	authOptions *AuthOptions
	userType    reflect.Type
}

var (
	auth_instance *authInstance = nil
)

type CastUser func(interface{}) interface{}

// Options stores configurations for a CSRF middleware.
type AuthOptions struct {
	ErrorFunc    kelly.HandlerFunc
	User         interface{}
	CastUserFunc CastUser
}

func defaultErrorFunc(c *kelly.Context) {
	c.WriteJson(http.StatusUnauthorized, kelly.H{
		"code":    http.StatusUnauthorized,
		"message": http.StatusText(http.StatusUnauthorized),
	})
}

func defaultCastUser(user interface{}) interface{} {
	return user
}

func checkUserType(user interface{}) reflect.Type {
	t := reflect.TypeOf(user)
	if t.Kind() != reflect.Ptr {
		panic("must be pointer")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("must be struct")
	}
	return t
}

// 自动注入登录的user信息
func AuthMiddleware(options *AuthOptions) kelly.HandlerFunc {
	if auth_instance != nil {
		log.Panic("init auth yet")
	}
	if options == nil || options.User == nil {
		log.Panic("invalid options")
	}

	if options.ErrorFunc == nil {
		options.ErrorFunc = defaultErrorFunc
	}
	if options.CastUserFunc == nil {
		options.CastUserFunc = defaultCastUser
	}

	tp := checkUserType(options.User)
	gob.Register(options.User)
	auth_instance = &authInstance{
		authOptions: options,
		userType:    tp,
	}

	return func(c *kelly.Context) {
		session := c.MustGet(AUTH_SESSION_NAME).(Session)
		value := session.Get(AUTH_SESSION_KEY)

		if value != nil{
			tp := checkUserType(value)
			if tp != auth_instance.userType{
				log.Printf("invalid user type,skip")
				value = nil
			}
		}

		value = auth_instance.authOptions.CastUserFunc(value)
		c.Set(AUTH_SESSION_KEY, value)
		c.InvokeNext()
	}
}

// 必须要登录的中间件检查
func LoginRequired() kelly.HandlerFunc {
	if auth_instance == nil {
		panic("not init yet")
	}
	return func(c *kelly.Context) {
		if IsAuthenticated(c) {
			c.InvokeNext()
		} else {
			auth_instance.authOptions.ErrorFunc(c)
		}
	}
}

// 是否已经登录
func IsAuthenticated(c *kelly.Context) bool {
	user := c.MustGet(AUTH_SESSION_KEY)
	return user != nil
}

// 当前登录的用户
func LoggedUser(c *kelly.Context) interface{}{
	user := c.MustGet(AUTH_SESSION_KEY)
	return user
}

// 登录
func Login(c *kelly.Context, user interface{}) error {
	if auth_instance == nil {
		panic("not init yet")
	}

	tp := checkUserType(user)
	if tp != auth_instance.userType{
		panic("invalid user type")
	}

	session := c.MustGet(AUTH_SESSION_NAME).(Session)
	session.Set(AUTH_SESSION_KEY, user)
	if err := session.Save(); err != nil {
		return err
	}

	// 更新c对象中的Value
	c.Set(AUTH_SESSION_KEY, user)
	return nil
}

// 注销
func Logout(c *kelly.Context) error {
	session := c.MustGet(AUTH_SESSION_NAME).(Session)
	session.Delete(AUTH_SESSION_KEY)
	if err := session.Save(); err != nil {
		return err
	}

	// 更新c对象中的Value
	c.Set(AUTH_SESSION_KEY, nil)
	return nil
}
