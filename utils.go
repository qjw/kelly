// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import (
	"encoding/json"
	"github.com/qjw/kelly/binding"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
)

func wrapHandlerFunc(f HandlerFunc) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		f(newContext(rw, r, next))
	}
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func JsonConfToStruct(path string, obj interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, obj); err != nil {
		return err
	}
	return nil
}

func Validate(obj interface{}) error {
	return binding.Validate(obj)
}

func Version() string {
	return version
}
