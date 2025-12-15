package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/hun9k/gapi/db"
	"gorm.io/gorm"
)

// 新增
func Insert[M any](o Option, m M) error {
	return gorm.G[M](db.Inst(o.dbKey)).
		Create(o.ctx, &m)
}

// 删除
func DelById[M any](o Option, id uint) (int, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	return gorm.G[M](d).
		Where(id).
		Delete(o.ctx)
}
func DelByIds[M any](o Option, ids []uint) (int, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	return gorm.G[M](d).
		Where(ids).
		Delete(o.ctx)
}
func DelByQuery[M any](o Option, query string, args ...any) (int, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	return gorm.G[M](d).
		Where(query, args...).
		Delete(o.ctx)
}

// 恢复
func RecoverById[M any](o Option, id uint) (int, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}

	return gorm.G[M](d).
		Where(id).
		Update(o.ctx, "deleted_at", nil)
}
func RecoverByIds[M any](o Option, ids []uint) (int, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}

	return gorm.G[M](d).
		Where(ids).
		Update(o.ctx, "deleted_at", nil)
}

// 修改
func UpdateById[M any](o Option, m M, cols []string, id uint) (int, error) {
	return gorm.G[M](db.Inst(o.dbKey)).
		Where(id).
		Updates(o.ctx, m)
}
func UpdateByIds[M any](o Option, m M, cols []string, ids []uint) (int, error) {
	return gorm.G[M](db.Inst(o.dbKey)).
		Where(ids).
		Updates(o.ctx, m)
}
func UpdateByQuery[M any](o Option, m M, cols []string, query string, args ...any) (int, error) {
	return gorm.G[M](db.Inst(o.dbKey)).
		Where(query, args...).
		Updates(o.ctx, m)
}

// 查询
func Select[M any](o Option, id uint) (M, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}

	return gorm.G[M](d).
		Where(id).
		First(o.ctx)
}

// the rest clauses: Filter, Order, Limiter
func SelectRows[M any](o Option, clauses ...any) ([]M, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}

	// query init
	q := gorm.G[M](d).Where(nil)

	// select clauses
	for _, c := range clauses {
		switch v := c.(type) {
		case ColList:
			q = q.Select(v.Fields())
		case Filter:
			// where
			query, args := v.QueryArgs()
			q = q.Where(query, args...)
		case Order:
			// order
			q = q.Order(v.Order())
		case Limiter:
			offset, limit := v.OffsetLimit()
			q.Offset(offset).Limit(limit)
		}
	}

	return q.Find(o.ctx)
}

// 统计
func Count[M any](o Option, col string, f Filter) (int64, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}

	query, args := f.QueryArgs()
	return gorm.G[M](d).
		Where(query, args...).
		Count(o.ctx, col)
}
func CountAll[M any](o Option, col string) (int64, error) {
	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}

	return gorm.G[M](d).
		Count(o.ctx, col)
}

// 执行选项
type Option struct {
	dbKey    string
	ctx      context.Context
	unscoped bool // 是否忽略软删除
}

func MkOpt(options ...any) Option {
	op := Option{
		dbKey:    db.DEFAULT_KEY,
		ctx:      context.Background(),
		unscoped: false,
	}

	for _, opt := range options {
		switch v := opt.(type) {
		case string:
			op.dbKey = v
		case context.Context:
			op.ctx = v
		case bool:
			op.unscoped = v
		}

	}
	return op
}

type ColList interface {
	Fields() string
}

type Cols []string

func (s Cols) Fields() string {
	return strings.Join(s, ", ")
}

type Filter interface {
	QueryArgs() (string, []any)
}

type CommonFilter map[string]any

func (f CommonFilter) QueryArgs() (string, []any) {
	qs := []string{}
	args := []any{}
	for k, v := range f {
		switch v.(type) {
		case string, int, int64, float64:
			qs = append(qs, fmt.Sprintf("`%s` = ?", k))
			args = append(args, v)
		case bool:
			if v.(bool) {
				qs = append(qs, "`%s` = 1")
			} else {
				qs = append(qs, "`%s` = 0")
			}
		case []uint, []int, []string, []float64:
			qs = append(qs, fmt.Sprintf("`%s` IN (?)", k))
			args = append(args, v)
		}
	}
	query := strings.Join(qs, " AND ")
	return query, args
}

type Order interface {
	Order() string
}

// field, asc|desc
type Sort [2]string

func (s Sort) Order() string {
	return fmt.Sprintf("`%s` %s", s[0], s[1])
}

// []field, asc|desc
type Sorts []Sort

func (s Sorts) Order() string {
	orders := []string{}
	for _, sort := range s {
		orders = append(orders, sort.Order())
	}
	return strings.Join(orders, ", ")
}

type Limiter interface {
	OffsetLimit() (int, int)
}

// start, end
type Range [2]int

func (r Range) OffsetLimit() (int, int) {
	return r[0], r[1] - r[0] + 1
}

// page, pagesize
type page [2]int

func (p page) OffsetLimit() (int, int) {
	return (p[0] - 1) * p[1], p[1]
}

// offset, limit
type Limit [2]int

func (l Limit) OffsetLimit() (int, int) {
	return l[0], l[1]
}
