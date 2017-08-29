// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/xml"
	"net/http"
)

var xmlContentType = "application/xml; charset=utf-8"

func WriteXml(w http.ResponseWriter, code int, data interface{}) error {
	writeContentType(w, xmlContentType)
	w.WriteHeader(code)

	return xml.NewEncoder(w).Encode(data)
}
