// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import "net/http"

type ServeHTTP func(http.ResponseWriter, *http.Request) (*http.Request, bool)
