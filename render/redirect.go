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
	"net/http"
)

func Redirect(w http.ResponseWriter,code int,r *http.Request,location string) error {
	if (code < http.StatusMultipleChoices || code > http.StatusPermanentRedirect) && code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", code))
	}
	http.Redirect(w, r, location, code)
	return nil
}
