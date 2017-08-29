// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/gin-gonic/gin
package render

import (
	"encoding/json"
	"net/http"
)

// bug https://github.com/golang/go/issues/14914

const jsonContentType = "application/json; charset=utf-8"

func WriteIndentedJSON(w http.ResponseWriter, code int, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
	return nil
}

func WriteJSON(w http.ResponseWriter, code int, obj interface{}) error {
	writeContentType(w, jsonContentType)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(obj)
}
