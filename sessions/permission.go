// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package sessions

import (
	"log"
	"net/http"
	"sync"
	"github.com/qjw/kelly"
)

const (
	PERMISSION_SESSION_KEY = "_permission_"
)

// 全局对象
type permissionInstance struct {
	permissionOptions *PermissionOptions
	// 懒加载的锁
	mu                sync.Mutex
	// 懒加载
	permissions       *map[string]int
}

// 全局对象实例
var (
	permission_instance *permissionInstance = nil
)

// 获取用户的所有权限
type UsePermissionGetter func(interface{}) (map[int]bool, error)

// 获取所有的权限
type AllPermisionsGetter func() (map[string]int, error)

// 选项，对外
type PermissionOptions struct {
	ErrorFunc            kelly.HandlerFunc
	UserPermissionGetter UsePermissionGetter
	AllPermisionsGetter  AllPermisionsGetter
}

// 默认错误处理函数
func defaultPermErrorFunc(c *kelly.Context) {
	c.WriteJson(http.StatusForbidden, kelly.H{
		"code":    http.StatusForbidden,
		"message": http.StatusText(http.StatusForbidden),
	})
}

// 初始化
func InitPermission(options *PermissionOptions) error {
	if permission_instance != nil {
		log.Panic("init permission yet")
	}
	if options == nil || options.UserPermissionGetter == nil || options.AllPermisionsGetter == nil {
		log.Panic("invalid options")
	}

	if options.ErrorFunc == nil {
		options.ErrorFunc = defaultPermErrorFunc
	}

	permission_instance = &permissionInstance{
		permissionOptions: options,
	}
	return nil
}

func getAllPermission(){
	permission_instance.mu.Lock()
	defer permission_instance.mu.Unlock()

	if permission_instance.permissions != nil{
		return
	}

	permissions,err := permission_instance.permissionOptions.AllPermisionsGetter()
	if err != nil{
		log.Printf("get all permission err %s",err.Error())
		return
	}
	permission_instance.permissions = &permissions
}

// 必须要登录的中间件检查
func PermissionRequired(perm string) kelly.HandlerFunc {
	if permission_instance == nil {
		panic("not init yet")
	}

	if permission_instance.permissions == nil{
		getAllPermission()
		if permission_instance.permissions == nil{
			panic("can NOT get all permissions")
		}
	}

	permId,ok := (*permission_instance.permissions)[perm]
	if !ok{
		log.Fatalf("invalid perm name %s", perm)
	}

	return func(c *kelly.Context) {
		user := LoggedUser(c)
		if user == nil {
			permission_instance.permissionOptions.ErrorFunc(c)
			return
		}

		session := c.MustGet(AUTH_SESSION_NAME).(Session)
		value := session.Get(PERMISSION_SESSION_KEY)
		if value != nil {
			permissions, ok := value.(map[int]bool)
			if !ok {
				panic("invalid permissions")
			}

			_, ok = permissions[permId]
			if !ok {
				permission_instance.permissionOptions.ErrorFunc(c)
				return
			}
		} else {
			permissions, err := permission_instance.permissionOptions.UserPermissionGetter(user)
			if err != nil {
				log.Printf("get permission err %s", err.Error())
				permission_instance.permissionOptions.ErrorFunc(c)
				return
			}

			_, ok := permissions[permId]
			if !ok {
				permission_instance.permissionOptions.ErrorFunc(c)
				return
			}
		}
		c.InvokeNext()
	}
}
