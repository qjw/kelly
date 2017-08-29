# 基于[gin binding](https://github.com/gin-gonic/gin/tree/master/binding)的自动绑定库

在原来基础上做了如下修改
1. 将原来的校验库从<http://gopkg.in/go-playground/validator.v8>替换成<http://gopkg.in/go-playground/validator.v9>
2. 增加gin PATH变量的binding
3. 将所有binding的字段都统一到tag的json标签，去除之前的form标签依赖
4. 增加可选变量绑定到指针的支持（若不存在，指针为空）
5. Validate函数暴露出来，用于只需要校验的场景

校验规则，请参考 **<http://godoc.org/gopkg.in/go-playground/validator.v9>**

使用前，运行命令 **./dep.sh** 来下载工具 **govendor** 以及使用govendor同步所有依赖的库

# 用法
``` go
func Test(c *gin.Context) {
    var form CreateMenuParam
    err, _ := binding.Bind(c, form);
}
```

若无法**自动判定**使用的bind后端（form/json等），可以使用
``` go
// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
func BindJSON(c * gin.Context,obj interface{}) (error,[]string) {
	return BindWith(c, obj, JSON)
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func BindWith(c * gin.Context,obj interface{}, b Binding) (error,[]string) {
	err := b.Bind(c, obj)
	return parseError(err,obj)
}
```

若已经完成了绑定，**只需要校验**，可以使用
``` go
func Validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
```