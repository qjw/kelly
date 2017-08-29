// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"io"
	"net/http"
)

const plainContentType = "text/plain; charset=utf-8"

func WriteString(w http.ResponseWriter, code int, format string, data []interface{}) error{
	writeContentType(w, plainContentType)
	w.WriteHeader(code)

	if len(data) > 0 {
		if _,err := fmt.Fprintf(w, format, data...);err != nil{
			return err
		}
	} else {
		if _,err := io.WriteString(w, format);err != nil{
			return err
		}
	}
	return nil
}
