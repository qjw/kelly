// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"github.com/qjw/kelly/binding"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"reflect"
)

type binder interface {
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

	GetBindParameter() interface{}
	GetBindJsonParameter() interface{}
	GetBindXmlParameter() interface{}
	GetBindFormParameter() interface{}
	GetBindPathParameter() interface{}
}

type binderImp struct {
	c *Context
}

func newBinder(c *Context) binder {
	return &binderImp{
		c: c,
	}
}

func (b *binderImp) Bind(obj interface{}) (error, []string) {
	bind := binding.Default(b.c.Request().Method, b.c.ContentType())
	return bindWith(b.c, obj, bind)
}

func (b *binderImp) BindJson(obj interface{}) (error, []string) {
	return bindWith(b.c, obj, binding.JSON)
}

func (b *binderImp) BindXml(obj interface{}) (error, []string) {
	return bindWith(b.c, obj, binding.XML)
}

func (b *binderImp) BindForm(obj interface{}) (error, []string) {
	return bindWith(b.c, obj, binding.Form)
}

func (b *binderImp) BindPath(obj interface{}) (error, []string) {
	return bindWith(b.c, obj, &pathBinding{
		r: b.c,
	})
}

func parseError(err error, obj interface{}) (error, []string) {
	if err == nil {
		return nil, nil
	} else {
		tips := make([]string, 0)
		real_err, ok := err.(validator.ValidationErrors)
		if ok {
			objt := reflect.TypeOf(obj).Elem()
			if objt.Kind() != reflect.Struct {
				return real_err, tips
			}

			for _, v := range real_err {
				elem, ok := objt.FieldByName(v.StructField())
				if !ok {
					continue
				}
				str, ok := elem.Tag.Lookup("error")
				if ok {
					// log.Printf("tag : %s\n",str)
					tips = append(tips, str)
				}
			}
			return real_err, tips
		} else {
			return err, tips
		}
	}
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func bindWith(c *Context, obj interface{}, b binding.Binding) (error, []string) {
	err := b.Bind(c.Request(), obj)
	return parseError(err, obj)
}

func handleValidateErr(c *Context, err error, msgs []string, obj interface{}) {
	c.WriteJson(http.StatusUnprocessableEntity, H{
		"code":  http.StatusUnprocessableEntity,
		"error": err.Error(),
		"msgs":  msgs,
		"obj":   obj,
	})
}

const (
	contextBindKey     = "_bind_key"
	contextBindJsonKey = "_bind_json_key"
	contextBindXmlKey  = "_bind_xml_key"
	contextBindFormKey = "_bind_form_key"
	contextBindPathKey = "_bind_path_key"
)

func (b *binderImp) GetBindParameter() interface{} {
	return b.c.MustGet(contextBindKey)
}

func (b *binderImp) GetBindJsonParameter() interface{} {
	return b.c.MustGet(contextBindJsonKey)
}

func (b *binderImp) GetBindXmlParameter() interface{} {
	return b.c.MustGet(contextBindXmlKey)
}

func (b *binderImp) GetBindFormParameter() interface{} {
	return b.c.MustGet(contextBindFormKey)
}

func (b *binderImp) GetBindPathParameter() interface{} {
	return b.c.MustGet(contextBindPathKey)
}

func BindMiddleware(objG func()interface{}) HandlerFunc {
	return func(c *Context) {
		obj := objG()
		err, msgs := c.Bind(obj)
		if err == nil {
			c.Set(contextBindKey, obj)
			c.InvokeNext()
		} else {
			handleValidateErr(c, err, msgs, obj)
		}
	}
}

func BindJsonMiddleware(objG func()interface{}) HandlerFunc {
	return func(c *Context) {
		obj := objG()
		err, msgs := c.BindJson(obj)
		if err == nil {
			c.Set(contextBindJsonKey, obj)
			c.InvokeNext()
		} else {
			handleValidateErr(c, err, msgs, obj)
		}
	}
}

func BindXmlMiddleware(objG func()interface{}) HandlerFunc {
	return func(c *Context) {
		obj := objG()
		err, msgs := c.BindXml(obj)
		if err == nil {
			c.Set(contextBindXmlKey, obj)
			c.InvokeNext()
		} else {
			handleValidateErr(c, err, msgs, obj)
		}
	}
}

func BindFormMiddleware(objG func()interface{}) HandlerFunc {
	return func(c *Context) {
		obj := objG()
		err, msgs := c.BindForm(obj)
		if err == nil {
			c.Set(contextBindFormKey, obj)
			c.InvokeNext()
		} else {
			handleValidateErr(c, err, msgs, obj)
		}
	}
}

func BindPathMiddleware(objG func()interface{}) HandlerFunc {
	return func(c *Context) {
		obj := objG()
		err, msgs := c.BindPath(obj)
		if err == nil {
			c.Set(contextBindPathKey, obj)
			c.InvokeNext()
		} else {
			handleValidateErr(c, err, msgs, obj)
		}
	}
}
