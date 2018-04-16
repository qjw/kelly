package main

import (
	"fmt"
	"net/http"

	"github.com/qjw/kelly"
	"github.com/qjw/kelly/sessions/store"
)

func main() {
	router := kelly.New(store.FlashMiddleware([][]byte{[]byte("abcdefg")}...))

	router.GET("/", func(c *kelly.Context) {
		store.AddFlash(c, "hello world")
		store.AddFlash(c, "hello world2")
		store.AddFlash(c, "hello world3")
		c.Redirect(http.StatusFound, "/res")
	})

	router.GET("/res", func(c *kelly.Context) {
		msgs := store.Flashes(c)
		if len(msgs) > 0 {
			c.WriteJson(http.StatusOK, msgs)
		} else {
			c.WriteJson(http.StatusOK, kelly.H{
				"message": "",
			})
		}
	})

	router.Run(fmt.Sprintf("0.0.0.0:%d", 8888))
}
