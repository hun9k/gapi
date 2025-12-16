package dao

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hun9k/gapi/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 新增
func Insert[M any](o Option, m *M) error {
	return gorm.G[M](db.Inst(o.dbKey)).
		Create(o.ctx, m)
}

// 删除
func Del[M any](o Option, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)

	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	// query
	return gorm.G[M](d, clauses...).Delete(o.ctx)
}

// 修改
func Update[M any](o Option, m M, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)

	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	// query
	return gorm.G[M](d, clauses...).Updates(o.ctx, m)
}

func Restore[M any](o Option, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)

	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	// query
	return gorm.G[M](d, clauses...).Update(o.ctx, "created_at", nil)
}

// the rest clauses: ColList, Filter, Order, Limiter
func Select[M any](o Option, conds ...any) ([]M, error) {
	// clauses
	clauses := BuildClauses(conds...)

	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	// query
	return gorm.G[M](d, clauses...).Find(o.ctx)
}

// 统计
func Count[M any](o Option, col string, conds ...any) (int64, error) {
	// clauses
	clauses := BuildClauses(conds...)

	d := db.Inst(o.dbKey)
	if o.unscoped {
		d = d.Unscoped()
	}
	// query
	return gorm.G[M](d, clauses...).Count(o.ctx, col)
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

func BuildClauses(conds ...any) []clause.Expression {
	clauses := []clause.Expression{}
	for _, c := range conds {
		switch v := c.(type) {
		case Selector:
			if v != nil && !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.Select())
			}
		case Wherer:
			if v != nil && !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.Where())
			}
		case OrderByer:
			if v != nil && !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.OrderBy())
			}
		case Limiter:
			if v != nil && !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.Limit())
			}
		}
	}
	return clauses
}

type Selector interface {
	Select() clause.Select
}

type Wherer interface {
	Where() clause.Where
}

type OrderByer interface {
	OrderBy() clause.OrderBy
}

type Limiter interface {
	Limit() clause.Limit
}

type Cols []string

func (s Cols) Select() clause.Select {
	sel := clause.Select{
		Columns: []clause.Column{},
	}
	for _, v := range s {
		sel.Columns = append(sel.Columns, clause.Column{
			Name: v,
		})
	}
	return sel
}

type Filter map[string]any

func (f Filter) Where() clause.Where {
	where := clause.Where{
		Exprs: []clause.Expression{},
	}
	// do keyword search
	keyword, kok := f["_keyword"]
	search, sok := f["_search"]
	if kok && sok && search != "" && keyword != "" {
		where.Exprs = append(where.Exprs, clause.Expr{
			SQL:  fmt.Sprintf("`%s` LIKE ?", search),
			Vars: []any{fmt.Sprintf(`%%%s%%`, keyword)},
		})
		delete(f, "_keyword")
		delete(f, "_search")
	}

	for k, v := range f {
		switch vv := v.(type) {
		case []any:
			if len(vv) > 0 {
				where.Exprs = append(where.Exprs, clause.Expr{
					SQL:  fmt.Sprintf("`%s` IN (?)", k),
					Vars: []any{v},
				})
			} else {
				where.Exprs = append(where.Exprs, clause.Expr{
					SQL: fmt.Sprintf("`%s` != `%s`", k, k),
				})
			}
		case any:
			if vvv, ok := vv.(bool); ok {
				if vvv {
					where.Exprs = append(where.Exprs, clause.Expr{
						SQL: fmt.Sprintf("`%s` = 1", k),
					})
				} else {
					where.Exprs = append(where.Exprs, clause.Expr{
						SQL: fmt.Sprintf("`%s` = 0", k),
					})
				}
				continue
			} else {
				where.Exprs = append(where.Exprs, clause.Expr{
					SQL:  fmt.Sprintf("`%s` = ?", k),
					Vars: []any{v},
				})
			}
		case nil:
			where.Exprs = append(where.Exprs, clause.Expr{
				SQL: fmt.Sprintf("`%s` IS NULL", k),
			})
		}
	}
	return where
}

// field, asc|desc
type Order []struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}

func (s Order) OrderBy() clause.OrderBy {
	orderBy := clause.OrderBy{
		Columns: []clause.OrderByColumn{},
	}
	for _, v := range s {
		if v.Field == "" {
			continue
		}
		orderBy.Columns = append(orderBy.Columns, clause.OrderByColumn{
			Column: clause.Column{
				Name: v.Field,
			},
			Desc: v.Desc,
		})
	}
	return orderBy
}

// Pager, size
type Pager [2]int

func (p Pager) Limit() clause.Limit {
	page, size := 1, 12
	if p[0] > 0 {
		page = p[0]
	}
	if p[1] > 0 {
		size = p[1]
	}
	if size > 100 {
		size = 100
	}
	return clause.Limit{
		Limit:  &size,
		Offset: (page - 1) * size,
	}
}
