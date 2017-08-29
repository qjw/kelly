// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"
	"io"
	"html/template"
)

const htmlContentType = "text/html; charset=utf-8"

func WriteHtml(w http.ResponseWriter, code int, data string) error {
	writeContentType(w, htmlContentType)
	w.WriteHeader(code)

	if _,err := io.WriteString(w, data);err != nil{
		return err
	}
	return nil
}

func WriteTemplateHtml(w http.ResponseWriter, code int, temp *template.Template, data interface{}) error {
	writeContentType(w, htmlContentType)
	w.WriteHeader(code)

	return temp.Execute(w, data)
}