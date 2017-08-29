# 适配[gin-gonic/gin](https://github.com/gin-gonic/gin)的session管理

参考以下项目，因为改动非常大，所以并非基于某一个clone的

1. <https://github.com/gorilla/sessions>
1. <https://github.com/martini-contrib/sessions>


gorilla/sessions依赖于<https://github.com/gorilla/context>，后者内部依赖一个加锁的map，不是很中意。在1.7之后，内建了[context](https://golang.org/pkg/context/)模块，可以在一定程度上优化**gorilla/context**的问题。

之所以**一定程度**是因为context库并不会改变现有的http.Request，而是返回一个新的对象，这导致一个很严重的问题，除非直接修改传入的http.Request对象，否则就无法链式的调用下去，参见如下代码

``` go
// 注意传入的next
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        userContext:=context.WithValue(context.Background(),"user","张三")
        ageContext:=context.WithValue(userContext,"age",18)
        // 这里必须递归调用
        next.ServeHTTP(rw, r.WithContext(ageContext))
    })
}
```

在上面的代码中，通过context包可以在http.Request对象上附加信息，但是由于会生成新的http.Request对象，所以链式调用，后续的handle并不会读取到新添加的数据，在很多场景无法使用或者会导致代码很难看。

此问题参考

1. <http://www.flysnow.org/2017/07/29/go-classic-libs-gorilla-context.html>
1. <https://stackoverflow.com/questions/40199880/how-to-use-golang-1-7-context-with-http-request-for-authentication?rq=1>

martini-contrib/sessions也直接依赖gorilla/sessions，所以也需要优化

---

之所以基于martini-contrib/sessions来修改，是因为

1. 存在redis的store，因个人喜好，替换了一个新redis库[redis.v5](https://gopkg.in/redis.v5)
1. 适配了gin.Context对象

由于**gin.Context自带context**，可以直接附加数据，所以完全可以绕开<https://github.com/gorilla/context>和<https://golang.org/pkg/context/>

另外**增加了Store的delete**方法，用于删除整个cookie，而不是cookie里面的某个key。

# Session

目前支持三个session后端

1. cookie，session的内容全部序列化到cookie中返回到浏览器，Flash使用此方式
2. file，session的内容存在**本地文件**中，session的id通过cookie返回到浏览器
3. redis，session的内容存在**redis数据库**中，session的id通过cookie返回到浏览器

很少直接使用session

``` go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

func main() {
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		// 设置session。每个session都包含若干key/value对
		session, _ := store.Get(c, "session_test")
		session.Set("key", "value")
		// 保存
		store.Save(c, session)
		// 或者 保存所有的session
		// sessions.Save(c)

		c.Redirect(http.StatusFound, "/pong")
	})

	r.GET("/pong", func(c *gin.Context) {
		// 获取session的值
		session, _ := store.Get(c, "session_test")
		value := session.Get("key")
		if value != nil {
			c.JSON(200, gin.H{
				"message": value.(string),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "",
			})
		}
	})

	r.GET("/middle",
		sessions.GinSessionMiddleware(store,"session_test"),
		func(c *gin.Context) {
			// 使用中间件，自动设置session到gin.Context中，避免大量的全局变量传递
			session := c.MustGet("session").(sessions.Session)
			value := session.Get("key")
			if value != nil {
				c.JSON(200, gin.H{
					"message": value.(string),
				})
			} else {
				c.JSON(200, gin.H{
					"message": "",
				})
			}
		})
	r.Run("0.0.0.0:9090")
}
```

# Flask

由于 **[gorilla/securecookie](https://github.com/gorilla/securecookie)** 需要一个初始密钥进行加密，所以初始化有个密钥的参数

``` go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"net/http"
)

func main() {
	r := gin.Default()
	sessions.InitFlash([]byte("abcdefghijklmn"))

	r.GET("/ping", func(c *gin.Context) {
		sessions.AddFlash(c, "hello world")
		c.Redirect(http.StatusFound, "/pong")
	})

	r.GET("/pong", func(c *gin.Context) {
		msgs := sessions.Flashes(c)
		if len(msgs) > 0 {
			c.JSON(200, gin.H{
				"message": msgs[0].(string),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "",
			})
		}
	})
	r.Run("0.0.0.0:9090")
}
```

输入<http://127.0.0.1:9090/ping> 自动跳转到<http://127.0.0.1:9090/pong>，并且显示ping设置的"hello world"

# 认证

实际情况中，只有少量接口不需要授权，所以实行**白名单**的方式会比较简单，但由于缺乏好的机制，这里还是传统的黑名单方式，即需要授权的接口自行添加中间件进行权限检查

> 所谓白名单就是做一个全局过滤（默认全部都需要授权），其中保存一个列表，在列表中请求的放开。

目前有两个难点

1. 缺乏好的机制标示某个请求，一般使用请求url，问题是存在PATH变量的情况很麻烦
2. 没有一种机制能够在http handle处自动注入，因为无法获取当前handle的消息，要不就统一编码白名单

``` go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

type User struct {
	Id   int
	Name string
}

func initStore() sessions.Store{
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
	r := gin.Default()
	r.Use(sessions.GinSessionMiddleware(store, sessions.AUTH_SESSION_NAME))
	r.Use(sessions.GinAuthMiddleware(&sessions.AuthOptions{
		User:&User{},
	}))

	r.GET("/index",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/login",
		func(c *gin.Context) {
			// 是否已经登录
			if sessions.IsAuthenticated(c){
				c.Redirect(http.StatusFound, "/index")
				return
			}
			// 登录授权
			sessions.Login(c,&User{
				Id:1,
				Name:"king",
			})
			c.Redirect(http.StatusFound, "/index")
		})
	r.GET("/logout",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 注销登录
			sessions.Logout(c)
			c.JSON(http.StatusFound, "/logout")
		})
	r.Run("0.0.0.0:9090")
}

```

# 授权

考虑到不同的系统，权限实现有所区别(通常都使用角色来归类权限)，这里做了一个简单的抽象
``` go
// 获取用户的所有权限
type UsePermissionGetter func(interface{}) (map[int]bool, error)

// 获取所有的权限
type AllPermisionsGetter func() (map[string]int, error)
```

权限通过string/int的map存储，在sessions.PermissionRequired("perm3")来作权限控制时，事先转换未权限的ID

在中间件处理中，调用另外一个接口获取当前用户的所有权限，并且作cache，当同一请求后续的操作中可以直接使用。

```
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

type User2 struct {
	Id   int
	Name string
}

func initStore2() sessions.Store {
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
	store := initStore2()
	r := gin.Default()
	r.Use(sessions.GinSessionMiddleware(store, sessions.AUTH_SESSION_NAME))
	r.Use(sessions.GinAuthMiddleware(&sessions.AuthOptions{
		User: &User2{},
	}))
	sessions.InitPermission(&sessions.PermissionOptions{
		UserPermissionGetter: func(user interface{}) (map[int]bool, error) {
			ruser := user.(*User2)
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

	r.GET("/index",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/perm1",
		sessions.PermissionRequired("perm1"),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/perm2",
		sessions.PermissionRequired("perm2"),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/perm3",
		sessions.PermissionRequired("perm3"),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/login",
		func(c *gin.Context) {
			// 是否已经登录
			if sessions.IsAuthenticated(c) {
				c.Redirect(http.StatusFound, "/index")
				return
			}

			// 登录授权
			sessions.Login(c, &User2{
				Id:   1,
				Name: c.DefaultQuery("name", "p1"),
			})
			c.Redirect(http.StatusFound, "/index")
		})
	r.GET("/logout",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 注销登录
			sessions.Logout(c)
			c.JSON(http.StatusFound, "/logout")
		})
	r.Run("0.0.0.0:9090")
}
```

## todo

1. 多权限的and/or支持
2. 权限层级关系