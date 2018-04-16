package store

import (
	"github.com/qjw/kelly"
)

type Store interface {
	New(c *kelly.Context, options *Options, name string) (Session, error)
	Save(c *kelly.Context, s Session) error
	Delete(c *kelly.Context, s Session) error
}
