// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kelly

import (
	"errors"
	"fmt"
	"github.com/qjw/kelly/binding"
	"net/http"
	"reflect"
	"strings"
)

type pathBinding struct {
	r request
}

func (pathBinding) Name() string {
	return "path"
}

func (obj *pathBinding) Bind(r *http.Request, data interface{}) error {
	value := reflect.ValueOf(data)
	if err := obj.read(r, value); err != nil {
		return err
	}
	return binding.Validate(data)
}

func (obj pathBinding) read(r *http.Request, val reflect.Value) error {
	// t := reflect.ValueOf(data).Type()
	typ := val.Type()
	switch typ.Kind() {
	case reflect.Struct:
		//typ := t.Elem()
		//val := v.Elem()
		count := typ.NumField()
		for i := 0; i < count; i++ {
			typeField := typ.Field(i)
			structField := val.Field(i)

			// 只能是基本类型
			fmt.Println(typeField.Type.Kind())
			if _, ok := kindMapping[typeField.Type.Kind()]; !ok {
				return errors.New("path object invalid field type")
			}

			if !structField.CanSet() {
				continue
			}

			tag := typeField.Tag.Get("json")
			name := parseTag(tag)
			if name == "" {
				name = typeField.Name
			}
			if name == "-" {
				continue
			}

			value, err := obj.r.GetPathVarible(name)
			if err != nil {
				return err
			}
			if err := binding.SetWithProperType(typeField.Type, value, structField); err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		return obj.read(r, val.Elem())
	default:
		return errors.New("path object invalid type")
	}
}

var kindMapping = map[reflect.Kind]string{
	reflect.Bool:    "boolean",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Uint:    "integer",
	reflect.Uint8:   "integer",
	reflect.Uint16:  "integer",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.String:  "string",
}

func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
