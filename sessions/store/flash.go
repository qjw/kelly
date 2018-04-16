package store

import (
	"fmt"

	"github.com/qjw/kelly"
)

const (
	maxAge          = 15       // flash session的生命周期
	defaultFlashKey = "_flash" // flash的session key
)

func FlashMiddleware(keyPairs ...[]byte) kelly.HandlerFunc {
	s := NewCookieStore(keyPairs...)
	return SessionMiddleware(s, &Options{
		MaxAge: &maxAge,
	}, defaultFlashKey)
}

// 添加flash消息
func AddFlash(c *kelly.Context, msg string) {
	s := GetSession(c, defaultFlashKey)
	v := s.Get(defaultFlashKey)
	if v == nil {
		s.Set(defaultFlashKey, []string{msg})
	} else {
		realVar := v.([]string)
		realVar = append(realVar, msg)
		s.Set(defaultFlashKey, realVar)
	}
	s.Save()
}

// 获取所有的flask，并且清空。
func Flashes(c *kelly.Context) []string {
	s := GetSession(c, defaultFlashKey)
	v := s.Get(defaultFlashKey)
	s.Clear()
	s.Save()
	if v == nil {
		return []string{}
	} else {
		realVar := v.([]string)
		return realVar
	}
}
