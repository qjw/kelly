// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

func WriteData(w http.ResponseWriter, code int, contentType string, data []byte) error {
	if len(contentType) > 0 {
		writeContentType(w, contentType)
	}
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
