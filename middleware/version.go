// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"github.com/qjw/kelly"
)

func Version(ver string) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		c.SetHeader("X-ACCOUNT-VERSION", ver)
		c.InvokeNext()
	}
}
