package main

import (
	"fmt"

	"github.com/qjw/kelly"
	"github.com/qjw/kelly/sessions/cache"
	"github.com/qjw/kelly/sessions/cache/redigo"
	"github.com/qjw/kelly/sessions/store"
)

const SESSION_NAME = "my_session1"
const SESSION_NAME2 = "my_session2"

func set(r kelly.Router, cookie string) {
	r.GET("/", func(c *kelly.Context) {
		s := store.GetSession(c, cookie)
		value := s.Get("key1")
		value2 := s.Get("key2")
		fmt.Printf("%v %v\n", value, value2)
		c.Abort(200, "ok")
	})

	r.GET("/set", func(c *kelly.Context) {
		s := store.GetSession(c, cookie)
		s.Set("key1", cookie+"aaa")
		s.Set("key2", cookie+"bbb")
		s.Save()
		c.Abort(200, "ok")
	})

	r.GET("/del", func(c *kelly.Context) {
		s := store.GetSession(c, cookie)
		s.Delete("key1")
		s.Save()
		c.Abort(200, "ok")
	})

	r.GET("/remove", func(c *kelly.Context) {
		s := store.GetSession(c, cookie)
		s.DeleteSelf()
		c.Abort(200, "ok")
	})
}

func main() {
	//	s := store.NewCookieStore([][]byte{[]byte("abcdefg")}...)
	c, _ := cache.NewJsonCache("/tmp/ccc")
	c, _ = redigo.NewCache("redis://server:6379/8")
	s := store.NewServerStore(c, [][]byte{[]byte("abcdefg")}...)

	router := kelly.New(store.SessionMiddleware(s, &store.Options{
		KeyPrefix: "tt2",
	}, SESSION_NAME))
	set(router, SESSION_NAME)
	maxAge := 10
	r := router.Group(
		"/api",
		store.SessionMiddleware(s, &store.Options{
			Path:      "/api",
			MaxAge:    &maxAge,
			KeyPrefix: "tt",
		}, SESSION_NAME2),
	)
	set(r, SESSION_NAME2)
	router.Run(fmt.Sprintf("0.0.0.0:%d", 8888))
}
