// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/middleware"
	"github.com/qjw/kelly/middleware/swagger"
	"github.com/qjw/kelly/sessions"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

const (
	ProjectRoot = "/home/king/code/go/src/github.com/qjw/kelly/sample"
)

func initStore() sessions.Store {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       3,
	})
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatal("failed to connect redis")
	}

	store, err := sessions.NewRediStore(redisClient, []byte("abcdefg"))
	if err != nil {
		log.Print(err)
	}
	return store
}

func main() {
	store := initStore()

	router := kelly.New(
		middleware.NoCache(),
		middleware.Secure(&middleware.SecureConfig{
			AllowedHosts: []string{"127.0.0.1:9090"},
		}),
	)
	kelly.EnableDebug(true)

	// 增加全局的endpoint钩子
	router.GlobalAnnotation(func(c *kelly.AnnotationContext) {
		handle := c.HandlerFunc()
		name := runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
		log.Printf("register [%7s|%2d]%s%s ---- %s",
			c.Method(), c.HandleCnt(), c.R().Path(), c.Path(), name)
	})

	router.SetNotFoundHandle(func(c *kelly.Context) {
		c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
	router.SetMethodNotAllowedHandle(func(c *kelly.Context) {
		c.WriteString(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})

	// swagger
	swagger.InitializeApiRoutes(router,
		&swagger.Config{
			BasePath:         "/swagger",
			Title:            "Swagger测试工具",
			Description:      "Swagger测试工具",
			DocVersion:       "0.1",
			SwaggerUiUrl:     "http://swagger.qiujinwu.com",
			SwaggerUrlPrefix: "doc",
			Debug:            true,
		},
		func(key string) ([]byte, error) {
			// 自行修改路径，key是文件名
			return ioutil.ReadFile(ProjectRoot + "/swagger.yaml")
		},
	)

	InitParam(router)
	InitGroupMiddleware(router)
	InitRender(router)
	InitStatic(router)
	InitApiV1(router, store)
	InitCsrf(router)
	InitMisc(router)
	InitBinding(router)
	InitSwagger(router)
	InitUpload(router)
	InitCaptcha(router)
	InitTemplate(router)

	router.Annotation(func(c *kelly.AnnotationContext) {
		log.Printf("have register %s%s %s", c.R().Path(), c.Path(), c.Method())
	}).GET("/", func(c *kelly.Context) {
		log.Print(c.GetDefaultCookie("session", "ss"))
		log.Print(c.MustGet("v1"))
		c.Redirect(http.StatusFound, "/doc")
	})

	router.Annotation(func(c *kelly.AnnotationContext) {
		log.Printf("have register %s%s %s", c.R().Path(), c.Path(), c.Method())
	}).GET("/health", func(c *kelly.Context) {
		c.WriteString(http.StatusOK, "ok")
	})

	router.Use(
		middleware.Version("v1"),
		Middleware("v1", "v1", true),
	)
	router.Run(":9090")
}
