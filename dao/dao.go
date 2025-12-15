package dao

import (
	"context"

	"github.com/hun9k/gapi/base"
	"github.com/hun9k/gapi/db"
)

type dao struct {
	// cache key, assoc with db key
	key string

	ctx context.Context
	m   any
}

var daos = map[string]*dao{}

func Inst(ks ...string) *dao {
	key := db.DEFAULT_KEY
	if len(ks) > 0 {
		key = ks[0]
	}
	if _, ok := daos[key]; !ok {
		daos[key] = &dao{
			key: key,
		}
		daos[key].Reset()
	}
	return daos[key]
}

func (d *dao) Ctx(ctx context.Context) *dao {
	d.ctx = ctx
	return d
}
func (d *dao) M(m any) *dao {
	d.m = m
	return d
}

func (d *dao) Reset() *dao {
	d.ctx = context.Background()
	d.m = &base.Model{}
	return d
}
