// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package kelly

import "fmt"

type Configuration struct {
	Port int    `json:"port" binding:"max=65536"`
	Host string `json:"host" binding:"ip4_addr"`
}

type ConfKelly interface {
	Kelly
	StartRun()
}

type confKelly struct {
	Kelly
	conf *Configuration
}

func (this confKelly) StartRun() {
	this.Run(fmt.Sprintf("%s:%d", this.conf.Host, this.conf.Port))
}

func InitByConf(path string) ConfKelly {
	conf := &Configuration{
		Port: 9090,
		Host: "127.0.0.1",
	}
	if err := JsonConfToStruct(path, conf); err != nil {
		panic(err)
	}

	if err := Validate(conf); err != nil {
		panic(err)
	}

	router := New()
	confKelly := &confKelly{
		Kelly: router,
		conf:  conf,
	}
	return confKelly
}
