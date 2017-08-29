// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/xml"
	"net/http"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(r *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return Validate(obj)
}
