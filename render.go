// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"encoding/xml"
	"github.com/qjw/kelly/render"
	"html/template"
	"net/http"
	"net/url"
)

type renderOp interface {
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
	// 返回重定向
	Redirect(int, string)
	// 设置header
	SetHeader(string, string)
	// 设置cookie
	SetCookie(string, string, int, string, string, bool, bool)

	Abort(int, string)
	ResponseStatusOK()
	ResponseStatusBadRequest(error)
	ResponseStatusUnauthorized(error)
	ResponseStatusForbidden(error)
	ResponseStatusNotFound(error)
	ResponseStatusInternalServerError(error)
}

type renderImp struct {
	http.ResponseWriter
	c *Context
}

func (r *renderImp) SetCookie(
	name string,
	value string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(r, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (r *renderImp) SetHeader(key, value string) {
	if len(value) == 0 {
		r.Header().Del(key)
	} else {
		r.Header().Set(key, value)
	}
}

func (r *renderImp) WriteJson(code int, obj interface{}) {
	if err := render.WriteJSON(r, code, obj); err != nil {
		panic(err)
	}
}

func (r *renderImp) WriteIndentedJson(code int, obj interface{}) {
	if err := render.WriteIndentedJSON(r, code, obj); err != nil {
		panic(err)
	}
}

func (r *renderImp) WriteHtml(code int, data string) {
	if err := render.WriteHtml(r, code, data); err != nil {
		panic(err)
	}
}

func (r *renderImp) WriteTemplateHtml(code int, temp *template.Template, data interface{}) {
	if err := render.WriteTemplateHtml(r, code, temp, data); err != nil {
		panic(err)
	}
}

func (r *renderImp) WriteXml(code int, obj interface{}) {
	if err := render.WriteXml(r, code, obj); err != nil {
		panic(err)
	}
}

func (r *renderImp) WriteString(code int, format string, values ...interface{}) {
	if err := render.WriteString(r, code, format, values); err != nil {
		panic(err)
	}
}

func (r *renderImp) Redirect(code int, location string) {
	if err := render.Redirect(r, code, r.c.Request(), location); err != nil {
		panic(err)
	}
}

func (r *renderImp) WriteData(code int, contentType string, data []byte) {
	render.WriteData(r, code, contentType, data)
}

func (r *renderImp) Abort(code int, msg string) {
	if len(msg) == 0 {
		msg = http.StatusText(code)
	}
	r.WriteJson(code, H{
		"code":    code,
		"message": msg,
	})
}

func (r *renderImp) ResponseStatusOK() {
	r.Abort(http.StatusOK, "")
}
func (r *renderImp) ResponseStatusBadRequest(err error) {
	if err != nil {
		r.Abort(http.StatusBadRequest, err.Error())
	} else {
		r.Abort(http.StatusBadRequest, "")
	}
}
func (r *renderImp) ResponseStatusUnauthorized(err error) {
	if err != nil {
		r.Abort(http.StatusUnauthorized, err.Error())
	} else {
		r.Abort(http.StatusUnauthorized, "")
	}
}
func (r *renderImp) ResponseStatusForbidden(err error) {
	if err != nil {
		r.Abort(http.StatusForbidden, err.Error())
	} else {
		r.Abort(http.StatusForbidden, "")
	}
}
func (r *renderImp) ResponseStatusNotFound(err error) {
	if err != nil {
		r.Abort(http.StatusNotFound, err.Error())
	} else {
		r.Abort(http.StatusNotFound, "")
	}
}
func (r *renderImp) ResponseStatusInternalServerError(err error) {
	if err != nil {
		r.Abort(http.StatusInternalServerError, err.Error())
	} else {
		r.Abort(http.StatusInternalServerError, "")
	}
}

func newRender(c *Context) renderOp {
	return &renderImp{
		ResponseWriter: c,
		c:              c,
	}
}

//-------------------------------------------------------------------------------------------
// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
type H map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}
