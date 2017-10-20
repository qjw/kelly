// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/sessions"
	"net/http"
)

type User struct {
	Id   int
	Name string
}

func InitApiV1(r kelly.Router, store sessions.Store) {

	sessions.InitFlash([]byte("abcdefghijklmn"))

	sessions.InitPermission(&sessions.PermissionOptions{
		UserPermissionGetter: func(user interface{}) (map[int]bool, error) {
			ruser := user.(*User)
			if ruser.Name == "p1" {
				return map[int]bool{
					1: true,
				}, nil
			} else if ruser.Name == "p2" {
				return map[int]bool{
					1: true,
					2: true,
				}, nil
			} else {
				return map[int]bool{}, nil
			}
		},
		AllPermisionsGetter: func() (map[string]int, error) {
			return map[string]int{
				"perm1": 1,
				"perm2": 2,
				"perm3": 3,
			}, nil
		},
	})

	api := r.Group("/api/v1",
		sessions.SessionMiddleware(store, sessions.AUTH_SESSION_NAME),
		sessions.AuthMiddleware(&sessions.AuthOptions{
			User: &User{},
		}),
	)

	api.GET("/flask", func(c *kelly.Context) {
		sessions.AddFlash(c, "hello world")
		c.Redirect(http.StatusFound, "/api/v1/flask_res")
	})

	api.GET("/flask_res", func(c *kelly.Context) {
		msgs := sessions.Flashes(c)
		if len(msgs) > 0 {
			c.WriteJson(http.StatusOK, kelly.H{
				"message": msgs[0].(string),
			})
		} else {
			c.WriteJson(http.StatusOK, kelly.H{
				"message": "",
			})
		}
	})

	api.GET("/",
		sessions.LoginRequired(),
		func(c *kelly.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User)
			c.WriteJson(http.StatusOK, kelly.H{
				"message": user.Name,
			})
		})
	api.GET("/p1",
		sessions.PermissionRequired("perm1"),
		func(c *kelly.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User)
			c.WriteJson(http.StatusOK, kelly.H{
				"perm": "p1",
				"user": user.Name,
			})
		})
	api.GET("/p2",
		sessions.PermissionRequired("perm2"),
		func(c *kelly.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User)
			c.WriteJson(http.StatusOK, kelly.H{
				"perm": "p2",
				"user": user.Name,
			})
		})
	api.GET("/p3",
		sessions.PermissionRequired("perm3"),
		func(c *kelly.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User)
			c.WriteJson(http.StatusOK, kelly.H{
				"perm": "p3",
				"user": user.Name,
			})
		})

	api.GET("/login",
		func(c *kelly.Context) {
			// 是否已经登录
			if sessions.IsAuthenticated(c) {
				c.Redirect(http.StatusFound, "/api/v1/")
				return
			}

			// 登录授权
			sessions.Login(c, &User{
				Id:   1,
				Name: c.GetDefaultQueryVarible("name", "p1"),
			})
			c.Redirect(http.StatusFound, "/api/v1/")
		})

	api2 := api.Group("/",sessions.LoginRequired())
	api2.GET("/logout",
		func(c *kelly.Context) {
			// 注销登录
			sessions.Logout(c)
			c.WriteJson(http.StatusFound, "/logout")
		})
}
