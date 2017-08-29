// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"reflect"
)

func normalize(values []string) []string {
	if values == nil {
		return nil
	}
	distinctMap := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.ToLower(value)
		if _, seen := distinctMap[value]; !seen {
			normalized = append(normalized, value)
			distinctMap[value] = true
		}
	}
	return normalized
}

type converter func(string) string

func convert(s []string, c converter) []string {
	var out []string
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}


// Check if an option is assigned
func isNonEmptyOption(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() != 0
	case reflect.Bool:
		return v.IsValid()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.Interface, reflect.Ptr, reflect.Func:
		return !v.IsNil()
	}
	return false
}

func setHttpCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(strconv.Itoa(code) + " - " + http.StatusText(code)))
}