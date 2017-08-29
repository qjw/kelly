// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"github.com/qjw/kelly"
	"net/http"
	"time"
)

func NoCache() kelly.HandlerFunc {
	return func(c *kelly.Context) {
		c.SetHeader("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.SetHeader("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.SetHeader("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.InvokeNext()
	}
}
