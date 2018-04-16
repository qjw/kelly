package main

import (
	"fmt"
	"time"

	"github.com/qjw/kelly/sessions/cache"
	"github.com/qjw/kelly/sessions/cache/redigo"
)

func redigoTest() {
	c, err := redigo.NewCache("redis://:pwd@localhost:6379/8")
	if err != nil {
		panic(err)
	}
	c.Set("fuck", "you")
	fmt.Println(c.Get("fuck"))
	fmt.Println(c.Exisits("fuck") == nil)
	c.Delete("fuck")
	fmt.Println(c.Exisits("fuck") == nil)
	c.SetWithExpire("fuck", "shit", 1*time.Second)
	fmt.Println(c.Exisits("fuck") == nil)
	time.Sleep(2 * time.Second)
	fmt.Println(c.Exisits("fuck") == nil)
}

func fileTest() {
	c, err := cache.NewCache("/tmp/cache")
	if err != nil {
		panic(err)
	}
	c.Set("fuck", "you")
	fmt.Println(c.Get("fuck"))
	//	c.Delete("fuck")
	//	c.SetWithExpire("fuck", "shit", 50*time.Second)
}

func jsonTest() {
	c, err := cache.NewJsonCache("/tmp/cache")
	if err != nil {
		panic(err)
	}
	c.Set("fuck", "you")
	fmt.Println(c.Get("fuck"))
	fmt.Println(c.Exisits("fuck") == nil)
	c.Delete("fuck")
	fmt.Println(c.Exisits("fuck") == nil)
	c.SetWithExpire("fuck", "shit", 1*time.Second)
	fmt.Println(c.Exisits("fuck") == nil)
	time.Sleep(2 * time.Second)
	fmt.Println(c.Exisits("fuck") == nil)
}

func main() {
	//	redigoTest()
	//	fileTest()
	jsonTest()
}
