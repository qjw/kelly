package store

import (
	"fmt"

	"github.com/qjw/kelly"
)

const (
	authSessionName = "_session"
	sessionPrefix   = "ss."
)

type Session interface {
	// 获取session的value
	Get(key string) interface{}
	// 设置session
	Set(key string, val interface{})
	// 删除session记录
	Delete(key string)
	// 删除所有的session记录
	Clear()
	// 保存session
	Save()
	// 删除整个session
	DeleteSelf()
}

func GetSession(c *kelly.Context, name string) Session {
	// 从Context获得session的实例
	return c.MustGet(fmt.Sprintf("%s%s", sessionPrefix, name)).(Session)
}

func SessionMiddleware(store Store, options *Options, name string) kelly.HandlerFunc {
	newOptions := conbineOptions(options)
	return func(c *kelly.Context) {
		session, _ := store.New(c, newOptions, name)
		//		cs := session.(*clientSession)

		c.Set(fmt.Sprintf("%s%s", sessionPrefix, name), session)
		c.InvokeNext()

		// 自动保存
		//		fmt.Printf("fuck %v\n", cs.Written)
		//		if cs.Written {
		//			fmt.Printf("save\n")
		//			store.Save(c, session)
		//		}
	}
}
