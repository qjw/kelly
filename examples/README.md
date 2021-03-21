# 代码实例

## 安装依赖
``` bash
$ make mod
```

## 编译
``` bash
# 编译所有
$ make build
# 编译单个sample， 见下面的标题名称
$ make helloworld
```

> 下面是每个实例， 标题名称就是Makefile`编译名称`

# helloworld
``` bash
$ ./build/helloworld
```
浏览器打开<http://127.0.0.1:9999/>

# route
演示
+ 路由分组
+ 路由分层嵌套

``` bash
$ ./build/route
```

浏览器打开

+ <http://127.0.0.1:9999/>
+ <http://127.0.0.1:9999/api/v1>
+ <http://127.0.0.1:9999/api/v2>
+ <http://127.0.0.1:9999/api/v1/v3/>

# response
演示
+ 输出Json（格式化/紧凑格式）
+ html（原始格式/模板渲染）
+ 文本（支持format)
+ XMl
+ 重定向
+ 二进制（图片）

``` bash
$ ./build/response
```

浏览器打开

+ <http://127.0.0.1:9999/>
+ <http://127.0.0.1:9999/json>
+ <http://127.0.0.1:9999/xml>
+ <http://127.0.0.1:9999/str>
+ <http://127.0.0.1:9999/html>
+ <http://127.0.0.1:9999/image>
+ <http://127.0.0.1:9999/template>
+ <http://127.0.0.1:9999/redirect>