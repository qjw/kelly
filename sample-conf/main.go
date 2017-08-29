// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import "github.com/qjw/kelly"

func main() {
	router := kelly.InitByConf("/home/king/code/go/src/github.com/qjw/kelly/sample-conf/conf.json")

	router.GET("/", func(c *kelly.Context) { c.ResponseStatusOK() })

	router.StartRun()
}
