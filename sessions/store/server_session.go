package store

import (
	"fmt"

	"github.com/qjw/kelly"
)

type serverSession struct {
	*clientSession
	// session id
	ID string
}

func (this serverSession) Key() string {
	if len(this.Options.KeyPrefix) > 0 {
		return fmt.Sprintf("%s.%s", this.Options.KeyPrefix, this.ID)
	} else {
		return this.ID
	}
}

func (this *serverSession) Save() {
	this.Store.Save(this.Context, this)
	this.Written = false
}

func (this *serverSession) DeleteSelf() {
	this.Store.Delete(this.Context, this)
}

func newServerSession(store Store,
	context *kelly.Context,
	options *Options,
	name string,
) *serverSession {
	cs := newClientSession(store, context, options, name)
	return &serverSession{
		clientSession: cs,
	}
}
