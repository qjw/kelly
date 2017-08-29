// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package toolkits

import (
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/qjw/kelly"
	gotemplate "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

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

var (
	gManage TemplateManage = nil
)

const (
	templateContextKey = "_template"
)

func InitTemplateMiddleware(manage TemplateManage) {
	if gManage != nil {
		panic("init yet")
	}
	gManage = manage
}

// 通过中间件预编译模板，并设置到context
func TemplateMiddleware(path string) kelly.HandlerFunc {
	if gManage == nil {
		panic("not init yet")
	}
	temp := gManage.MustGetTemplate(path)
	return func(c *kelly.Context) {
		c.Set(templateContextKey, temp)
		c.InvokeNext()
	}
}

// 获得当前context的模板
func CurrentTemplate(c *kelly.Context) Template {
	return c.MustGet(templateContextKey).(Template)
}

//----------------------------------------------------------------------------------------------------------------------

type template struct {
	temp *pongo2.Template
}

func (this *template) Render(c *kelly.Context, context kelly.H) error {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	return this.temp.ExecuteWriter(pongo2.Context(context), c)
}

func (this *template) MustRender(c *kelly.Context, context kelly.H) {
	if err := this.Render(c, context); err != nil {
		panic(err)
	}
}

type goTemplate struct {
	temp *gotemplate.Template
}

func (this *goTemplate) Render(c *kelly.Context, context kelly.H) error {
	c.WriteTemplateHtml(http.StatusOK, this.temp, context)
	return nil
}

func (this *goTemplate) MustRender(c *kelly.Context, context kelly.H) {
	c.WriteTemplateHtml(http.StatusOK, this.temp, context)
}

//----------------------------------------------------------------------------------------------------------------------

type templateManage struct {
	tset *pongo2.TemplateSet
}

func (this *templateManage) GetTemplate(temp string) (Template, error) {
	if tpl, err := this.tset.FromFile(temp); err == nil {
		return &template{temp: tpl}, nil
	} else {
		return nil, err
	}
}

func (this *templateManage) MustGetTemplate(temp string) Template {
	if tpl, err := this.tset.FromFile(temp); err == nil {
		return &template{temp: tpl}
	} else {
		panic(err)
	}
}

type goTemplateManage struct {
	root string
}

func (this *goTemplateManage) GetTemplate(temp string) (Template, error) {
	path := this.root + temp
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	} else if fi.IsDir() {
		return nil, fmt.Errorf("'%s' is dir", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t, err := gotemplate.New("ok").Parse(string(data))
	if err != nil {
		return nil, err
	}
	return &goTemplate{
		temp: t,
	}, nil

}

func (this *goTemplateManage) MustGetTemplate(temp string) Template {
	t, err := this.GetTemplate(temp)
	if err != nil {
		panic(err)
	}
	return t
}

//----------------------------------------------------------------------------------------------------------------------

// 创建新的模板管理器
func NewTemplateManage(path string) TemplateManage {
	pongo2_loader := pongo2.MustNewLocalFileSystemLoader(path)
	pongo2_set := pongo2.NewSet("default", pongo2_loader)
	return &templateManage{
		tset: pongo2_set,
	}
}

// 创建新的模板管理器
func NewGoTemplateManage(path string) TemplateManage {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if !fi.IsDir() {
		panic(fmt.Errorf("'%s' not a dir", path))
	}

	abpath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return &goTemplateManage{
		root: abpath + "/",
	}
}
