package store

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/qjw/kelly"
)

const (
	maxAge         = (7 * 24 * 60 * 60) // 默认认证生命周期（秒）
	defaultAuthKey = "_session"         // 认证的session key
)

var (
	gm              sync.Mutex
	initFlag        bool         = false // 是否已经初始化
	goptions        *AuthOptions = nil   // 全局唯一的配置
	gCurrentUserKey string       = "current_user"
)

// 错误处理的回调函数
type HandlerFunc func(*kelly.Context, error)

type CastUser func(interface{}) (interface{}, error)

type AuthOptions struct {
	ErrorFunc    HandlerFunc
	User         interface{}
	CastUserFunc CastUser
	Store        Store
}

func defaultErrorFunc(c *kelly.Context, err error) {
	if err == nil {
		err = fmt.Errorf("%d", http.StatusUnauthorized)
	}
	c.WriteJson(http.StatusUnauthorized, kelly.H{
		"code":    http.StatusUnauthorized,
		"message": err.Error(),
	})
}

func defaultCastUser(user interface{}) (interface{}, error) {
	return user, nil
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
	gm.Lock()
	defer gm.Unlock()
	if initFlag {
		panic("init auth yet")
	}
	initFlag = true

	if options == nil || options.User == nil {
		panic("invalid options")
	}

	if options.ErrorFunc == nil {
		options.ErrorFunc = defaultErrorFunc
	}
	if options.CastUserFunc == nil {
		options.CastUserFunc = defaultCastUser
	}

	tp := checkUserType(options.User)
	//	gob.Register(options.User)
	goptions = options

	// 默认的session中间件
	sessionMiddleware := SessionMiddleware(
		goptions.Store,
		&Options{
			MaxAge: &maxAge,
		},
		defaultAuthKey,
	)

	return func(c *kelly.Context) {
		session := c.MustGet(AUTH_SESSION_NAME).(Session)
		value := session.Get(AUTH_SESSION_KEY)

		if value != nil {
			tp := checkUserType(value)
			if tp != auth_instance.userType {
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
	if !initFlag {
		panic("not init yet")
	}
	return func(c *kelly.Context) {
		if IsAuthenticated(c) {
			c.InvokeNext()
		} else {
			goptions.ErrorFunc(c, nil)
		}
	}
}

// 是否已经登录
func IsAuthenticated(c *kelly.Context) bool {
	user := c.Get(gCurrentUserKey)
	return user != nil
}

// 当前登录的用户
func LoggedUser(c *kelly.Context) interface{} {
	user := c.Get(gCurrentUserKey)
	return user
}

// 登录
func Login(c *kelly.Context, user interface{}) error {
	if !initFlag {
		panic("not init yet")
	}

	// 必须是选项制定的类型
	tp := checkUserType(user)
	if tp != goptions.userType {
		panic("invalid user type")
	}

	s := GetSession(c, defaultAuthKey)
	s.Set(defaultAuthKey, user)
	if err := s.Save(); err != nil {
		return err
	}

	// 更新c对象中的Value
	c.Set(gCurrentUserKey, user)
	return nil
}

// 注销
func Logout(c *kelly.Context) error {
	if !initFlag {
		panic("not init yet")
	}

	s := GetSession(c, defaultAuthKey)
	s.DeleteSelf()

	// 更新c对象中的Value
	c.Set(gCurrentUserKey, nil)
	return nil
}
