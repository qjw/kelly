package store

import (
	"github.com/qjw/kelly"
)

type clientSession struct {
	// session 名称
	Name string
	// 所有数据
	Values map[string]interface{}
	// 所属的store
	Store Store

	Context *kelly.Context

	Written bool
	Options *Options
}

func (this clientSession) Get(key string) interface{} {
	return this.Values[key]
}

func (this *clientSession) Set(key string, val interface{}) {
	this.Values[key] = val
	this.Written = true
}

func (this *clientSession) Delete(key string) {
	delete(this.Values, key)
	this.Written = true
}

func (this *clientSession) Clear() {
	for key := range this.Values {
		this.Delete(key)
	}
}

func (this *clientSession) Save() {
	this.Store.Save(this.Context, this)
	this.Written = false
}

func (this *clientSession) DeleteSelf() {
	this.Store.Delete(this.Context, this)
}

func newClientSession(store Store,
	context *kelly.Context,
	options *Options,
	name string,
) *clientSession {
	return &clientSession{
		Name:    name,
		Store:   store,
		Options: options,
		Context: context,
		Values:  make(map[string]interface{}),
		Written: false,
	}
}
