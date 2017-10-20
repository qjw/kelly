   * [运行Sample](#运行sample)
   * [背景](#背景)
   * [Sample](#sample)
      * [参数](#参数)
         * [PATH变量](#path变量)
         * [Query变量](#query变量)
         * [Form变量](#form变量)
         * [获取Header](#获取header)
         * [获取Cookie](#获取cookie)
      * [文件上传](#文件上传)
   * [静态文件](#静态文件)
   * [输出](#输出)
      * [重定向](#重定向)
      * [设置Header](#设置header)
      * [设置Cookie](#设置cookie)
   * [Context数据](#context数据)
   * [中间件 middleware](#中间件-middleware)
      * [全局中间件](#全局中间件)
      * [动态添加中间件](#动态添加中间件)
      * [多级路由](#多级路由)
      * [空路由](#空路由)
      * [单个API中间件注入](#单个api中间件注入)
      * [其他Http方法](#其他http方法)
      * [处理404/405](#处理404405)
      * [重置Request](#重置request)
   * [数据绑定和校验](#数据绑定和校验)
      * [手动绑定](#手动绑定)
      * [校验规则](#校验规则)
         * [自定义错误输出](#自定义错误输出)
      * [自动绑定](#自动绑定)
         * [自动绑定实现](#自动绑定实现)
         * [单独校验](#单独校验)
   * [Session](#session)
      * [Flash](#flash)
      * [Session](#session-1)
         * [依赖](#依赖)
      * [认证授权](#认证授权)
   * [权限](#权限)
   * [内建中间件](#内建中间件)
      * [Csrf](#csrf)
      * [Cors](#cors)
   * [Annotation注解和Swagger](#annotation注解和swagger)
      * [打印所有的请求注册操作](#打印所有的请求注册操作)
      * [Swagger](#swagger)
         * [初始化](#初始化)
         * [Annotation 中间件](#annotation-中间件)
      * [router context](#router-context)
   * [验证码](#验证码)
      * [验证码使用](#验证码使用)
   * [模板](#模板)
      * [中间件](#中间件)
      * [Go内建模板](#go内建模板)
   * [二维码](#二维码)




# 运行Sample

> 为了避免干扰正式(测试)环境，建议先重置GOPATH环境变量

``` bash
# 重置GOPATH
king@king:~/tmp$ export GOPATH=/home/king/tmp/gopath
king@king:~/tmp$ echo $GOPATH
/home/king/tmp/gopath

# 安装
king@king:~/tmp$ go get github.com/qjw/kelly/sample

# 查看依赖
king@king:~/tmp/gopath/src$ find . -maxdepth 3 -type d | sed "/\.git/d"
.
./github.com
./github.com/dchest
./github.com/dchest/uniuri
./github.com/urfave
./github.com/urfave/negroni
./github.com/julienschmidt
./github.com/julienschmidt/httprouter
./github.com/go-playground
./github.com/go-playground/locales
./github.com/go-playground/universal-translator
./github.com/gorilla
./github.com/gorilla/securecookie
./github.com/qjw
./github.com/qjw/kelly
./gopkg.in
./gopkg.in/go-playground
./gopkg.in/go-playground/validator.v9
./gopkg.in/redis.v5
./gopkg.in/redis.v5/testdata
./gopkg.in/redis.v5/internal

king@king:~/tmp/gopath/src$ cd ../bin/

# 运行sample
king@king:~/tmp/gopath/bin$ ./sample
[negroni] listening on :9090
```

# 背景

作为web后端开发，标准的[net/http](https://golang.org/pkg/net/http/)非常高效灵活，足以适用非常多的场景，当然也有很多周边待补充，这就出现了各种web框架，甚至出现了替代默认的Http库的[valyala/fasthttp](https://github.com/valyala/fasthttp)

golang目前百花齐放，个人主要了解到的是两个项目

1. [beego: simple & powerful Go app framework](https://beego.me/)
1. [gin-gonic/gin](https://github.com/gin-gonic/gin)

beego没有实际用过，听说是大而全的项目，对开发者友好。不过由于了解甚少，草率评论并不合适，这里不作过多说明。

本着刨根问底的学习态度，最开始了解的是[martini](https://github.com/philsong/martini)，后查证效率偏低（大量用到[反射/reflect](https://golang.org/pkg/reflect/)），所以就进一步学习了[gin-gonic/gin](https://github.com/gin-gonic/gin)。

后者小巧灵活，学习成本低，并且提供了很多实用的补充，例如

1. 路由和中间件核心框架，路由基于[julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
1. gin.Context
1. binding
1. 校验，基于[go-playground/validator.v9](https://gopkg.in/go-playground/validator.v9)
1. Http Request工具函数，获取param/path/form/header/cookie等
1. Http Response工具函数，设置cookie，header，返回xml/json，返回template支持等
1. 内建的几个常用中间件

martini/gin都包含非常多的中间件，两者迁移非常容易，参考

1. <https://github.com/codegangsta/martini-contrib>
1. <https://github.com/gin-gonic/contrib>
1. <https://github.com/gin-contrib>

用久了，也发现gin也有一些问题

1. 依赖还是偏多（*虽然和很多库相比算较少的*），就写个hello world都下载半天依赖
1. 第三方middleware有的依赖gopkg.in的代码，另外一些依赖github.com的代码
1. gin.Context对Golang标准库[context](https://golang.org/pkg/context/)不友好
1. binding有一些问题，本人的优化版本在<https://github.com/qjw/go-gin-binding>
1. 虽然middleware很多，但选择性太多，质量参差不齐，不好选择，另外太多的第三方依赖不如将大部分常用的集成到一起来的方便。

经过多方对比考察，认为[urfave/negroni](https://github.com/urfave/negroni)作为路由/中间件基础框架非常合适，（*看看原型就知道他对[context](https://golang.org/pkg/context/)有多友好*）所以折腾就开始了。

``` golang
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

经过综合评估，决定自己弄个类似于gin的框架

原则是尽量踏着巨人的肩膀，避免一些通用组件重复造轮子，聚焦于优秀智慧的集成

> 本着尊重原作者的原则和对开源协议的尊重，我会尽量备注作者和出处，若有遗漏，请知会我<qiujinwu@gmail.com>

目前的主要工作包括

1. 基于[urfave/negroni](https://github.com/urfave/negroni)的核心框架
1. 基于[julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)的路由，并作优化以支持多级路由+路由middleware
1. binding基于<https://github.com/qjw/go-gin-binding>，后者代码来源于<https://github.com/gin-gonic/gin/tree/master/binding>
1. kelly.Context
1. 基于[go-playground/validator.v9](https://gopkg.in/go-playground/validator.v9)的校验，参考gin的代码
1. 常用的Request/Response工具函数，参考gin和其他一些框架，特别是<https://github.com/gin-gonic/gin/tree/master/render>
1. 复用[urfave/negroni](https://github.com/urfave/negroni)的recovery/log
1. 支持http 404/405的统一全局处理
1. 内建静态文件支持
1. 内建常用的middleware，含认证授权

# Sample
``` go
package main
import (
    "github.com/qjw/kelly"
    "net/http"
)

func main(){
    router := kelly.New()

    router.GET("/", func(c *kelly.Context) {
        c.WriteIndentedJson(http.StatusOK, kelly.H{
            "code":    "0",
        })
    })

    router.Run(":9090")
}
```
``` bash
king@king:~/tmp/gopath/src/sample$ go run main.go
[negroni] listening on :9090
```

## 参数

### PATH变量
``` go
// 根据key获取PATH变量值
GetPathVarible(string) (string, error)
// 根据key获取PATH变量值，若不存在，则panic
MustGetPathVarible(string) string
```

### Query变量
``` go
// 根据key获取QUERY变量值，可能包含多个（http://127.0.0.1:9090/path/abc?abc=bbb&abc=aaa）
GetMultiQueryVarible(string) ([]string, error)
// 根据key获取QUERY变量值，仅返回第一个
GetQueryVarible(string) (string, error)
// 根据key获取QUERY变量值，仅返回第一个,若不存在，则返回默认值
GetDefaultQueryVarible(string, string) string
// 根据key获取QUERY变量值，仅返回第一个,若不存在，则panic
MustGetQueryVarible(string) string
```

``` go
r.GET("/path/:name", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{
        "code":  "/path",
        "path":  c.MustGetPathVarible("name"), // 获取path参数
        "query": c.GetDefaultQueryVarible("abc", "def"), // 获取query参数
    })
})
```

### Form变量
``` go
// 根据key获取FORM变量值，可能get可能包含多个
GetMultiFormVarible(string) ([]string, error)
// 根据key获取FORM变量值，仅返回第一个
GetFormVarible(string) (string, error)
// 根据key获取FORM变量值，仅返回第一个,若不存在，则返回默认值
GetDefaultFormVarible(string, string) string
// 根据key获取FORM变量值，仅返回第一个,若不存在，则panic
MustGetFormVarible(string) string
```

``` go
r.GET("/form", func(c *kelly.Context) {
    data := `<form action="/form" method="post">
<p>First name: <input type="text" name="fname" /></p>
<p>Last name: <input type="text" name="lname" /></p>
<input type="submit" value="Submit" />
</form>`
    c.WriteHtml(http.StatusOK, data) // 返回html
})

r.POST("/form", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
        "code":        "/form",
        "first name":  c.GetDefaultFormVarible("fname", "fname"), // 获取form参数
        "second name": c.GetDefaultFormVarible("lname", "lname"),
    })
})
```

### 获取Header
``` go
// 根据key获取header值
GetHeader(string) (string, error)
// 根据key获取header值，若不存在，则返回默认值
GetDefaultHeader(string, string) string
// 根据key获取header值，若不存在，则panic
MustGetHeader(string) string
// Content-Type
ContentType() string
```

### 获取Cookie
``` go
// 根据key获取cookie值
GetCookie(string) (string, error)
// 根据key获取cookie值，若不存在，则返回默认值
GetDefaultCookie(string, string) string
// 根据key获取cookie值，若不存在，则panic
MustGetCookie(string) string
```

## 文件上传

``` go
// @ref http.Request.ParseMultipartForm
ParseMultipartForm() error
// 获取（上传的）文件信息
GetFileVarible(string) (multipart.File, *multipart.FileHeader, error)
MustGetFileVarible(string) (multipart.File, *multipart.FileHeader)
```

``` go
r.GET("/upload", func(c *kelly.Context) {
    data := `<form enctype="multipart/form-data" action="/upload" method="post">
<input type="file" name="file1" />
<input type="file" name="file2" />
<input type="submit" value="upload" />
</form>`
    c.WriteHtml(http.StatusOK, data) // 返回html
})

r.POST("/upload", func(c *kelly.Context) {
    c.ParseMultipartForm()

    file, handler := c.MustGetFileVarible("file1")
    defer file.Close()
    f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    io.Copy(f, file)

    file2, handler2 := c.MustGetFileVarible("file2")
    defer file2.Close()
    f2, err := os.OpenFile("./"+handler2.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f2.Close()
    io.Copy(f2, file2)

    c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
        "code":        "/upload",
        "first name":  handler.Filename, // 获取form参数
        "second name": handler2.Filename,
    })
})
```

# 静态文件
``` go
package main
import (
    "github.com/qjw/kelly"
    "net/http"
)

func main(){
    router := kelly.New()

    router.GET("/static/*path", kelly.Static(&kelly.StaticConfig{
        Dir:        http.Dir("/var/www/html"),
        Indexfiles: []string{"index.html"},
    }))

    router.GET("/static1/*path", kelly.Static(&kelly.StaticConfig{
        Dir:           http.Dir("/tmp"),
        EnableListDir: true,
    }))

    router.Run(":9090")
}
```

运行之后，可以访问<http://127.0.0.1:9090/static>和<http://127.0.0.1:9090/static1>

``` go
type StaticConfig struct {
    Dir           http.FileSystem
    // 是否支持枚举目录
    EnableListDir bool
    // 访问目录时，是否自动查找index
    Indexfiles    []string
}
```

参考

1. <https://github.com/urfave/negroni/blob/master/static.go>
1. <https://github.com/labstack/echo/blob/master/middleware/static.go>

# 输出
``` go
// 返回紧凑的json
WriteJson(int, interface{})
// 返回xml
WriteXml(int, interface{})
// 返回html
WriteHtml(int, string)
// 返回模板html
WriteTemplateHtml(int, *template.Template, interface{})
// 返回格式化的json
WriteIndentedJson(int, interface{})
// 返回文本
WriteString(int, string, ...interface{})
// 返回二进制数据
WriteData(int, string, []byte)
```
``` go
render.GET("/t", func() kelly.HandlerFunc {
    data := `<form action="#" method="get">
<p>First {{ .First }}: <input type="text" name="fname" /></p>
<p>Last {{ .Last }}: <input type="text" name="lname" /></p>
<input type="submit" value="Submit" />
</form>`

    // 通过闭包预先编译好
    t := template.Must(template.New("t1").Parse(data))
    return func(c *kelly.Context) {
        c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
            "First": "Qiu",
            "Last": "King",
        })
    }
}())

render.GET("/a", func(c *kelly.Context) {
    c.WriteString(http.StatusOK, "test %d %d", 123, 456) // 返回普通文本
})
```

## 重定向
``` go
// 返回重定向
Redirect(int, string)
```
``` go
c.Redirect(http.StatusFound, "/api/v1/flask_res")
```

## 设置Header
``` go
// 设置header
SetHeader(string, string)
```
``` go
func Version(ver string) kelly.HandlerFunc {
    return func(c *kelly.Context) {
        c.SetHeader("X-ACCOUNT-VERSION", ver)
        c.InvokeNext()
    }
}
```

## 设置Cookie
``` go
// 设置cookie
SetCookie(string, string, int, string, string, bool, bool)
```

# Context数据
由于Context在一个中间件链执行，为了方便传递数据，支持Context读写数据。*比如auth中间件就会保存current_user变量*

``` go
func Middleware(ver string) kelly.HandlerFunc {
    return func(c *kelly.Context) {
        // 设置context参数
        c.Set("v1", ver)

        // 调用下一个handle
        c.InvokeNext()
    }
}

router.GET("/", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{
        "code":    "/",
        "value":   c.MustGet("v1"),  // 获取context数据
    })
})
```
``` go
Set(interface{}, interface{}) dataContext
Get(interface{}) interface{}
MustGet(interface{}) interface{}
```

# 中间件 middleware

一个最简单的中间件实现

``` go
type HandlerFunc func(c *Context)
```
``` go
func Version(ver string) kelly.HandlerFunc {
    return func(c *kelly.Context) {
        c.SetHeader("X-ACCOUNT-VERSION", ver)
        c.InvokeNext()
    }
}
```

默认会中断中间件链的执行，若希望继续执行，需要手动调用【**c.InvokeNext()**】

``` go
func (c *Context) InvokeNext() {
    if c.next != nil {
        c.next.ServeHTTP(c, c.Request())
    }
}
```

## 全局中间件

全局中间件会被 **所有** 的请求执行

``` go
router := kelly.New(
    middleware.Version("v1"),
)
```

## 动态添加中间件
``` go
router.Use(
    middleware.Version("v1"),
    Middleware("v1", "v1", true),
)
```

## 多级路由

子路由也支持注入中间件，只影响它自己的请求，以及他的子路由

``` go
router := kelly.New(
    middleware.Version("v1"),
)

// 新建一个子router，并注入一个middleware
ar := r.Group(
    "/aaa",
    Middleware("v2", "v2", true),
)
ar.GET("/", func(c *kelly.Context) {
    c.WriteJson(http.StatusOK, kelly.H{ // 返回json（紧凑格式）
        "code": "/aaa",
    })
})

// 新建一个子router，并注入一个middleware
sar := ar.Group(
    "/bbb",
    Middleware("v3", "v3", true),
)
sar.GET("/", func(c *kelly.Context) {
    c.WriteXml(http.StatusOK, kelly.H{  // 返回XML
        "code": "/aaa/bbb",
    })
})
```

## 空路由

空路由指在创建子路由（Group）时，使用路径"/"的路由，返回的新路由和父路由使用相同的路径

空路由的好处是可以针对同一个url的不同请求，使用不同的中间件

> 假如/api/v1一部分不需要登录，剩下的则需要。正常情况下，需要对需要登录的请求每个都加入中间件进行认证，而空路由则可以只注册一次中间件，需要登录的请求都基于这个空路由来注册。

``` go
api.GET("/login",
    func(c *kelly.Context) {
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
```

## 单个API中间件注入

> 下面的代码，通过一个中间件作登录认证

``` go
api.GET("/",
    sessions.LoginRequired(),
    func(c *kelly.Context) {
        // 获取登录用户
        user := sessions.LoggedUser(c).(*User)
        c.WriteJson(http.StatusOK, kelly.H{
            "message": user.Name,
        })
    })
```

## 其他Http方法
``` go
GET(string, ...HandlerFunc)
HEAD(string, ...HandlerFunc)
OPTIONS(string, ...HandlerFunc)
POST(string, ...HandlerFunc)
PUT(string, ...HandlerFunc)
PATCH(string, ...HandlerFunc)
DELETE(string, ...HandlerFunc)
```

``` go
r.POST("/ok", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
        "code": "/csrf ok",
    })
})
```

## 处理404/405
``` go
router.SetNotFoundHandle(func(c *kelly.Context) {
    c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
})
router.SetMethodNotAllowedHandle(func(c *kelly.Context) {
    c.WriteString(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
})
```

```go
// 设置404处理句柄
SetNotFoundHandle(HandlerFunc)
// 设置405处理句柄
SetMethodNotAllowedHandle(HandlerFunc)
```

## 重置Request

例如使用[context](https://golang.org/pkg/context/)就需要替换默认的Request对象

``` go
func (c *Context) SetRequest(r *http.Request){
```

# 数据绑定和校验
## 手动绑定

> 会同时绑定和校验

``` go
type BindPathObj struct {
    A string `json:"aaa" binding:"required,max=32,min=6" error:"aerror"`
    B string `json:"bbb" binding:"required,max=32,min=6" error:"berror"`
    C string `json:"ccc" binding:"required,max=32,min=6" error:"cerror"`
}

api.GET("/path/:aaa/:bbb/:ccc", func(c *kelly.Context) {
    var obj BindPathObj
    if err, _ := c.BindPath(&obj); err == nil {
        c.WriteJson(http.StatusOK, obj)
    } else {
        c.WriteString(http.StatusOK, "param err")
    }
})

type BindJsonObj struct {
    Obj1 BindPathObj `json:"obj"`
    A    string      `json:"aaa" binding:"required,max=32,min=6" error:"aerror"`
    B    string      `json:"bbb" binding:"required,max=32,min=6" error:"berror"`
    C    string      `json:"ccc" binding:"required,max=32,min=6" error:"cerror"`
}

api.POST("/json", func(c *kelly.Context) {
    var obj BindJsonObj
    if err, _ := c.Bind(&obj); err == nil {
        c.WriteJson(http.StatusOK, obj)
    } else {
        c.WriteString(http.StatusOK, "param err")
    }
})
```

所有的绑定接口

``` go
// 绑定一个对象，根据Content-type自动判断类型
Bind(interface{}) (error, []string)
// 绑定json，从body取数据
BindJson(interface{}) (error, []string)
// 绑定xml，从body取数据
BindXml(interface{}) (error, []string)
// 绑定form，从body/query取数据
BindForm(interface{}) (error, []string)
// 绑定path变量
BindPath(interface{}) (error, []string)
```

## 校验规则

参考<https://godoc.org/gopkg.in/go-playground/validator.v9>

### 自定义错误输出

参考struct tag中的error

``` go
type BindJsonObj struct {
    Obj1 BindPathObj `json:"obj"`
    A    string      `json:"aaa" binding:"required" error:"aerror"`
    B    string      `json:"bbb" binding:"required" error:"berror"`
    C    string      `json:"ccc" binding:"required" error:"cerror"`
}
```

## 自动绑定

> 自动绑定可以减少非常多拖沓的重复代码

**注意，可以在bind时设置缺省参数**

``` go
api.GET("/path2/:aaa/:bbb/:ccc",
    kelly.BindPathMiddleware(&BindPathObj{
        AAA: "testa",
        BBB: "testb",
    }),
    func(c *kelly.Context) {
        c.WriteJson(http.StatusOK, c.GetBindPathParameter())
    })
api.POST("/form2",
    kelly.BindMiddleware(&BindPathObj{}),
    func(c *kelly.Context) {
        c.WriteJson(http.StatusOK, c.GetBindParameter())
    })
api.POST("/json2",
    kelly.BindMiddleware(&BindJsonObj{}),
    func(c *kelly.Context) {
        c.WriteJson(http.StatusOK, c.GetBindParameter())
    })
```
``` go
GetBindParameter() interface{}
GetBindJsonParameter() interface{}
GetBindXmlParameter() interface{}
GetBindFormParameter() interface{}
GetBindPathParameter() interface{}
```
### 自动绑定实现

> 本质上就是将原来每个接口都需要写的重复逻辑抽象到中间件，并且通过Context传递

``` go
func BindMiddleware(obj interface{}) HandlerFunc {
    return func(c *Context) {
        // 绑定对象
        err, msgs := c.Bind(obj)
        if err == nil {
            // 使用一个固定的key存储到request的Context中
            c.Set(contextBindKey, obj)
            // 继续
            c.InvokeNext()
        } else {
            handleValidateErr(c, err, msgs, obj)
        }
    }
}
```

### 单独校验

校验框架可以自动从http请求获取参数并校验，当然也可以单独对已经存在的struct进行校验

``` go
type Configuration struct {
    Port int    `json:"port" binding:"max=65536"`
    Host string `json:"host" binding:"ip4_addr"`
}

func F(){
    if err := kelly.Validate(conf); err != nil {
        panic(err)
    }
}
```

# Session

## Flash

flash用于在多个后端接口传递数据

> flask基于cookie的session，不依赖于redis等文件系统/数据库

``` go
sessions.InitFlash([]byte("abcdefghijklmn"))

api.GET("/flash", func(c *kelly.Context) {
    sessions.AddFlash(c, "hello world")
    c.Redirect(http.StatusFound, "/api/v1/flash_res")
})

api.GET("/flash_res", func(c *kelly.Context) {
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
```

## Session

对于简单的应用，可以使用基于cookie的session，即全部（加密）内容都通过cookie传输，其他的建议使用基于服务器的session，即正文存储在后端的存储/文件系统中，只将key通过cookie传输。

> 一些其他的方案，参考[json web token](https://jwt.io/)

``` go
// gopkg.in/redis.v5
// 初始化redis，返回一个store对象
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

store := initStore()

// 注入session的中间件，用于将session实例存入Context
api := r.Group("/api/v1",
    sessions.SessionMiddleware(store, sessions.AUTH_SESSION_NAME),
)

func(c *kelly.Context) {
    // 从Context获得session的实例
    session := c.MustGet(AUTH_SESSION_NAME).(Session)
    // 从session读取内容
    value := session.Get(AUTH_SESSION_KEY)
}
```
``` go
type Session interface {
    // Get returns the session value associated to the given key.
    Get(key interface{}) interface{}
    // Set sets the session value associated to the given key.
    Set(key interface{}, val interface{})
    // Delete removes the session value associated to the given key.
    Delete(key interface{})
    // Clear deletes all values in the session.
    Clear()
    // Options sets confuguration for a session.
    // Options(Options)
    // Save saves all sessions used during the current request.
    Save() error
}
```

### 依赖

1. <https://github.com/gorilla/sessions>
1. <https://github.com/martini-contrib/sessions>

最终修改的项目见<https://github.com/qjw/sessions>

## 认证授权

授权依赖[sessions](https://github.com/qjw/kelly/tree/master/sessions)

``` go
// 在注入session的中间件之后，注入
store := initStore()
api := r.Group("/api/v1",
    sessions.SessionMiddleware(store, sessions.AUTH_SESSION_NAME),
    sessions.AuthMiddleware(&sessions.AuthOptions{
        User: &User{},
    }),
)

// 通过一个中间件作登录权限认证
api.GET("/",
    sessions.LoginRequired(),
    func(c *kelly.Context) {
        // 获取登录用户
        user := sessions.LoggedUser(c).(*User)
        c.WriteJson(http.StatusOK, kelly.H{
            "message": user.Name,
        })
    })

// 登录
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

        // 重定向到首页
        c.Redirect(http.StatusFound, "/api/v1/")
    })

// 登出
api.GET("/logout",
    sessions.LoginRequired(),
    func(c *kelly.Context) {
        // 注销登录
        sessions.Logout(c)
        c.WriteJson(http.StatusFound, "/logout")
    })
```

# 权限

初始化，需要提供几个数据
1. 所有的权限id/名称
2. 一个通过user查询所有权限的函数

接下来使用中间件sessions.PermissionRequired实现自动权限判断

> 权限判断会自动检查，是否已登录

``` go

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
```


# 内建中间件
1. [Version](https://github.com/qjw/kelly/blob/master/middleware/version.go)
1. [NoCache](https://github.com/qjw/kelly/blob/master/middleware/no_cache.go)
1. [BasicAuth](https://github.com/qjw/kelly/blob/master/middleware/basic_auth.go)，来自<https://github.com/martini-contrib/auth>
1. [Throttle](https://github.com/qjw/kelly/blob/master/middleware/throttle.go)：请求频率限制，来自<https://github.com/martini-contrib/throttle>
1. [Cors](https://github.com/qjw/kelly/blob/master/middleware/cors.go),来自<https://github.com/gin-contrib/cors>
1. [Secure](https://github.com/qjw/kelly/blob/master/middleware/secure.go),来自<https://github.com/gin-contrib/secure>
1. [Csrf](https://github.com/qjw/kelly/blob/master/middleware/csrf.go),来自<https://github.com/tommy351/gin-csrf>
1. [Gzip](https://github.com/qjw/kelly/blob/master/middleware/gzip.go),来自<https://github.com/gin-contrib/gzip>，支持gzip/deflate

## Csrf
``` go
// 初始化
middleware.InitCsrf(middleware.CsrfConfig{
    Secret: []byte("fasdffasdfas"),
})

// 注入Middleware
api := r.Group("/csrf",
    middleware.Csrf(),
)

api.GET("/ok", func() kelly.HandlerFunc {
    // token放在一个hidden表单中自动带入
    data := `<form action="/csrf//ok" method="post">
<p>First {{ .First }}: <input type="text" name="fname" /></p>
<p><input type="hidden" name="_csrf" value="{{ .Token }}"> </p>
<input type="submit" value="Submit" />
</form>`

    // 通过闭包预先编译好
    t := template.Must(template.New("ok").Parse(data))
    return func(c *kelly.Context) {
        // 在前一个请求中，返回token给前端
        c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
            "First": "Qiu",
            "Token": middleware.GetCsrfToken(c),
        })
    }
}())

api.POST("/ok", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{ // 返回格式化的json
        "code": "/csrf ok",
    })
})
```


## Cors

> 由于实现的原因，未绑定的Path，中间件无法监听，而cors依赖于options方法，所以为了支持，需要手动添加options绑定

``` go
router := r.Group("/swagger", middleware.Cors(&middleware.CorsConfig{
    AllowAllOrigins: true,
    AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
    AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type"},
}))

// 绑定所有的options请求来支持中间件作跨域处理
router.OPTIONS("/*path", func(c *kelly.Context) {
    c.WriteString(http.StatusOK, "ok")
})
```

# Annotation注解和Swagger

很多语言都支持注解，例如Java和Python

1. <https://docs.oracle.com/javase/tutorial/java/annotations/index.html>
1. <http://www.infoq.com/cn/articles/cf-java-annotation>
1. <http://www.cnblogs.com/Jerry-Chou/archive/2012/05/23/python-decorator-explain.html>

Golang并不支持注解之类的语法，比较类似的是[Struct Tag](https://golang.org/pkg/reflect/#StructTag)

kelly实现了一个【**类似**】于注解的方法，在每次注册请求时被触发，回调的原型如下
``` go
type AnnotationHandlerFunc func(c *AnnotationContext)
```

AnnotationContext如下

``` go
// endpint Context，用于记录每个请求的信息
type AnnotationContext struct {
    // endpoint所属的路由对象，从这里可以获取他的Path
    r           Router
    // endpoint的Http方法，例如PUT
    method      string
    // endpoint的路径，例如 /aaa
    path        string
    // router的中间件链条
    middlewares []HandlerFunc
    // endpoint 自己的中间件链条，最后一个就是最终的Http请求处理函数
    handles     []HandlerFunc
}
```

**在程序真正监听端口提供服务之前，这条链条已经执行完毕**，所以并不影响运行性能

有两种方法添加注解

``` go
// 添加全局的 注解 函数。该router下面和子（孙）router下面的endpoint注册都会被触发
GlobalAnnotation(handles ...AnnotationHandlerFunc) Router

// 添加临时 注解 函数，只对使用返回的AnnotationRouter对象进行注册的endpoint有效
Annotation(handles ...AnnotationHandlerFunc) AnnotationRouter
```

GlobalAnnotation注入的函数会一直存在于router对象，以及它的子router。（务必在注册请求或者添加子router之前）

Annotation只对返回的新router对象有效（需要链式调用），相当于对单个router临时有效

``` go
router.Annotation(func(c *kelly.AnnotationContext) {
    log.Printf("have register %s%s %s", c.R().Path(), c.Path(), c.Method())
}).GET("/", func(c *kelly.Context) {
    log.Print(c.GetDefaultCookie("session", "ss"))
    log.Print(c.MustGet("v1"))
    c.Redirect(http.StatusFound, "/doc")
})
```

## 打印所有的请求注册操作

在**根Router**注入GlobalAnnotation即可

``` go
router := kelly.New()

// 增加全局的endpoint钩子
router.GlobalAnnotation(func(c *kelly.AnnotationContext) {
    handle := c.HandlerFunc()
    name := runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
    log.Printf("register [%7s|%2d|%2d]%s%s ---- %s",
        c.Method(), c.MiddlewareCnt(), c.HandleCnt(), c.R().Path(), c.Path(), name)
})
```

启动后的输出如下

``` bash
2017/09/19 20:44:30 register [    GET| 4| 1]/path/:name ---- main.InitParam.func1
2017/09/19 20:44:30 register [    GET| 4| 1]/form ---- main.InitParam.func2
2017/09/19 20:44:30 register [   POST| 4| 1]/form ---- main.InitParam.func3
2017/09/19 20:44:30 register [    GET| 5| 1]/aaa/ ---- main.InitGroupMiddleware.func1
```

也可以在**子router**，或者终端注册请求时(例如router.GET)添加

``` go
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
```

``` bash
2017/09/19 20:44:30 register [    GET| 4| 1]/ ---- main.main.func5
2017/09/19 20:44:30 have register / GET
2017/09/19 20:44:30 register [    GET| 4| 1]/health ---- main.main.func7
2017/09/19 20:44:30 have register /health GET
```

## Swagger

由于每次请求注册都会执行注入的回调AnnotationHandlerFunc，所以对于一些针对请求的业务相当有用，比如生成doc文档的swagger

### 初始化
``` go
// swagger
swagger.InitializeApiRoutes(router,
    &swagger.Config{
        BasePath:         "/api/v1",
        Title:            "Swagger测试工具",
        Description:      "Swagger测试工具",
        DocVersion:       "0.1",
        // swagger ui 用于显示文档，为了支持其他域名，需要后端开启cors
        SwaggerUiUrl:     "http://swagger.qiujinwu.com",
        // 文档访问的path，例如127.0.0.1:9090/doc
        SwaggerUrlPrefix: "doc",
        Debug:            true,
    }, // 默认可以直接通过struct生成文档，若依赖yaml文件，需要这个接口来loadyaml文件的内容
    func(key string) ([]byte, error) {
        // 自行修改路径，key是文件名
        return ioutil.ReadFile("/home/king/code/go/src/github.com/qjw/kelly/sample/swagger.yaml")
    },
)
```

### Annotation 中间件

基于yaml文件，留意**swagger.SwaggerFile**
``` go
router.Annotation(
    swagger.SwaggerFile("swagger.yaml:upload_material"),
).POST("/upload_material", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{
        "code": "0",
    })
})
```

基于struct对象，留意**swagger.Swagger**
``` go
router.Annotation(swagger.Swagger(&swagger.StructParam{
    ResponseData: &swagger.SuccessResp{},
    FormData:     &swaggerParam{},
    Summary:      "api1",
    Tags:         []string{"API接口"},
})).PATCH("/api1", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{
        "code": "0",
    })
})
```

## router context

在中间件链中，可以通过kelly.Context来读写数据实现中间件之间的数据传输。

在Annotation链中，也有类似的接口，具体保存在kelly.Router中

在swagger中，常见的场景是同一个router下面的api具有相同的Tag，所以我们可以在router层面写入一个全局的tag，然后每个API读取即可

留意 **swagger.SetGlobalParam**

``` go
// 增加中间件处理跨域问题
router := r.Group("/swagger", middleware.Cors(&middleware.CorsConfig{
    AllowAllOrigins: true,
    AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
    AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type"},
})).GlobalAnnotation(swagger.SetGlobalParam(&swagger.StructParam{
    Tags:         []string{"API接口"},
})).OPTIONS("/*path", func(c *kelly.Context) {
    c.WriteString(http.StatusOK, "ok")
})

router.Annotation(swagger.Swagger(&swagger.StructParam{
    ResponseData: &swagger.SuccessResp{},
    FormData:     &swaggerParam{},
    Summary:      "api1",
})).PATCH("/api1", func(c *kelly.Context) {
    c.WriteIndentedJson(http.StatusOK, kelly.H{
        "code": "0",
    })
})
```

# 验证码
``` go
// 生成验证码的ID，不是实际的验证码数字，参数是验证码的长度
func GenerateCaptchaID(len int) string
// request的url必须是 http(s)://host/path/{CaptchaID}.png
func ServerCaptcha(c *kelly.Context,width,height int)
// 验证验证码
func VerifyCaptcha(captchaID,captchaCode string) bool
```

> 新增依赖<https://github.com/dchest/captcha>

**github.com/dchest/captcha**内部使用了一个简化版的内建redis作为（类似于服务器session）验证码的容器，对外暴露一个id，这个ID用于

1. 请求验证码图片
2. 回传用于验证码校验

由于验证码并不是频繁调用，所以这种办法挺靠谱，分离id和实际的code最大的好处是，重新获取验证码，前端影响较小，因为通过内部容器做了一层映射，所以可以确保在验证码变化的情况下，验证码id保持一致。

一种简化的做法是去掉内建的映射容器，直接将加密过的验证码使用cookie或者其他方式传递和回传，保持无状态。

## 验证码使用

留意代码中的js，通过获取验证码图片时，补上【reload=1】参数即可实现更新，而无须修改html的所有验证码id

``` go
func InitCaptcha(r kelly.Router) {

    r.GET("/captcha", func() kelly.HandlerFunc {
        data := `<form action="/captcha" method="post">
<p>验证码: <input type="text" name="captchaStr" /></p>
<input type="hidden" name="captchaID" value="{{ .CaptchaID }}" /></p>
<img id="captcha-img" width="104" height="36" src="/captcha/image/{{ .CaptchaID }}.png" />
<input type="submit" value="提交" />
</form>
<script type="text/javascript" src="//cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
<script type="text/javascript">
$(function() {
    $("#captcha-img").click(function() {
      var captcha_url = $(this).attr("src").split("?")[0];
      captcha_url += "?reload=1&timestamp=" + new Date().getTime()
      $(this).attr("src",  captcha_url);
    });
})
</script>`

        // 通过闭包预先编译好
        t := template.Must(template.New("t1").Parse(data))
        return func(c *kelly.Context) {
            captchaID := toolkits.GenerateCaptchaID(4)
            c.WriteTemplateHtml(http.StatusOK, t, map[string]string{
                "CaptchaID": captchaID,
            })
        }
    }())

    r.GET("/captcha/image/:id", func(c *kelly.Context) {
        toolkits.ServerCaptcha(c,104,36)
    })

    r.POST("/captcha", func(c *kelly.Context) {
        if toolkits.VerifyCaptcha(
            c.MustGetFormVarible("captchaID"),
            c.MustGetFormVarible("captchaStr"),
        ){
            c.ResponseStatusOK()
        }else{
            c.ResponseStatusForbidden(nil)
        }
    })
}

```

# 模板
``` go
// 创建新的模板管理器
func NewTemplateManage(path string) TemplateManage

type TemplateManage interface {
    // 加载模板
    GetTemplate(string) (Template, error)
    MustGetTemplate(string) Template
}

type Template interface {
    // 渲染模板到kelly.Context
    Render(c *kelly.Context, context kelly.H) error
    MustRender(c *kelly.Context, context kelly.H)
}
```

> 新增依赖<https://github.com/flosch/pongo2>。语法参考后者[官网](https://www.florian-schlachter.de/?tag=pongo2)

``` go
func InitTemplate(r kelly.Router) {
    mng := toolkits.NewTemplateManage(ProjectRoot)
    r.GET("/template", func() kelly.HandlerFunc {
        temp := mng.MustGetTemplate("template/index.html")
        return func(c *kelly.Context) {
            temp.Render(c, kelly.H{
                "Body": "Kelly",
            })
        }
    }())
}
```

## 中间件
手动获取模板灵活性很高，不过可以用中间件简化逻辑

步骤

1. （全局）用TemplateManage接口初始化
2. 增加TemplateMiddleware中间件
3. 使用CurrentTemplate获取当前的template

``` go
toolkits.InitTemplateMiddleware(mng)
r.GET("/template2",
    toolkits.TemplateMiddleware("template/index.html"),
    func(c *kelly.Context) {
        toolkits.CurrentTemplate(c).Render(c, kelly.H{
            "Body": "Kelly",
        })
    })
```

## Go内建模板
只需要将初始化函数从NewTemplateManage 替换成**NewGoTemplateManage**即可

# 二维码

> 新增依赖<https://github.com/skip2/go-qrcode>。

``` go
func NewQRCode(content string, level int) (*Qrcode, error)
func (q *Qrcode) Image(size int) image.Image
func (q *Qrcode) Write(size int, out io.Writer) error
func (q *Qrcode) WriteFile(size int, filename string) error
func (q *Qrcode) WriteKelly(size int, c *kelly.Context) error
```

``` go
router.GET("/qrcode", func(c *kelly.Context) {
    qrcode,_ := toolkits.NewQRCode(c.MustGetQueryVarible("content"),toolkits.QrcodeMedium)
    qrcode.WriteKelly(400,c)
})
```
