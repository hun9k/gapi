package dao

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hun9k/gapi/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 执行选项
type Option struct {
	db       *gorm.DB
	ctx      context.Context
	unscoped Unscoped // 是否忽略软删除
}
type Unscoped bool

func MkOpt(options ...any) Option {
	op := Option{
		db:       db.Inst(db.DEFAULT_KEY),
		ctx:      context.Background(),
		unscoped: false,
	}

	for _, opt := range options {
		switch v := opt.(type) {
		case *gorm.DB:
			op.db = v
		case context.Context:
			op.ctx = v
		case Unscoped:
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
			if !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.Select())
			}
		case Wherer:
			if !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.Where())
			}
		case OrderByer:
			if !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.OrderBy())
			}
		case Limiter:
			if !reflect.ValueOf(v).IsNil() {
				clauses = append(clauses, v.Limit())
			}
		}
	}
	return clauses
}

type IDer interface {
	Where() clause.Expression
}

type ID uint

func (id ID) Where() clause.Expression {
	return clause.Expr{
		SQL:  "`id` = ?",
		Vars: []any{id},
	}
}

type IDser interface {
	Where() clause.Expression
}

type IDs []any

func (ids IDs) Where() clause.Expression {
	if len(ids) == 0 {
		return nil
	}

	return clause.Expr{
		SQL:  "`id` IN (?)",
		Vars: []any{ids},
	}
}

type Selector interface {
	Select() clause.Expression
}

type Wherer interface {
	Where() clause.Expression
}

type OrderByer interface {
	OrderBy() clause.Expression
}

type Limiter interface {
	Limit() clause.Expression
}

type Cols []string

func (s Cols) Select() clause.Expression {
	if len(s) == 0 {
		return nil
	}

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

func (f Filter) Where() clause.Expression {
	if len(f) == 0 {
		return nil
	}

	exprs := []clause.Expression{}
	keyword, kok := f["_keyword"]
	search, sok := f["_search"]

	// do keyword search
	if kok && sok && search != "" && keyword != "" {
		exprs = append(exprs, clause.Expr{
			SQL:  fmt.Sprintf("(`%s` LIKE ?)", search),
			Vars: []any{fmt.Sprintf(`%%%s%%`, keyword)},
		})
		delete(f, "_keyword")
		delete(f, "_search")
	}

	// do other search
	for k, v := range f {
		switch vv := v.(type) {
		case []any:
			if len(vv) > 0 {
				exprs = append(exprs, clause.Expr{
					SQL:  fmt.Sprintf("`%s` IN (?)", k),
					Vars: []any{v},
				})
			} else {
				exprs = append(exprs, clause.Expr{
					SQL: fmt.Sprintf("`%s` != `%s`", k, k),
				})
			}
		case any:
			if vvv, ok := vv.(bool); ok {
				if vvv {
					exprs = append(exprs, clause.Expr{
						SQL: fmt.Sprintf("`%s` = 1", k),
					})
				} else {
					exprs = append(exprs, clause.Expr{
						SQL: fmt.Sprintf("`%s` = 0", k),
					})
				}
				continue
			} else {
				exprs = append(exprs, clause.Expr{
					SQL:  fmt.Sprintf("`%s` = ?", k),
					Vars: []any{v},
				})
			}
		case nil:
			exprs = append(exprs, clause.Expr{
				SQL: fmt.Sprintf("`%s` IS NULL", k),
			})
		}
	}

	if len(exprs) == 0 {
		return nil
	}
	return clause.Where{
		Exprs: exprs,
	}
}

func CheckFilter(f Filter) Filter {
	if len(f) == 0 {
		return Filter{}
	}
	return f
}

func FilterIDs(f Filter) []any {
	if f == nil {
		return nil
	}

	v, exists := f["id"]
	if !exists {
		return nil
	}

	return v.([]any)
}

func (o Option) DB() *gorm.DB {
	if o.unscoped {
		return o.db.Unscoped()
	}
	return o.db
}

func (o Option) Ctx() context.Context {
	return o.ctx
}

func (o Option) Unscoped() bool {
	return o.Unscoped()
}

// [[field, ASC|DESC], ...]
type Sorts []Sort
type Sort []string

func (s Sorts) OrderBy() clause.Expression {
	if len(s) == 0 {
		return nil
	}

	var columns []clause.OrderByColumn

	for _, v := range s {
		columns = append(columns, clause.OrderByColumn{
			Column: clause.Column{
				Name: v[0],
			},
			Desc: strings.ToUpper(v[1]) == "DESC",
		})
	}

	return clause.OrderBy{
		Columns: columns,
	}
}

func CheckSort(s Sorts) Sorts {
	if len(s) == 0 {
		return Sorts{}
	}

	var s1 Sorts
	for _, v := range s {
		if len(v) == 0 || v[0] == "" {
			continue
		}
		if v[1] == "" {
			v[1] = "ASC"
		}
		s1 = append(s1, v)
	}

	return s1
}

// start - end
type Range []int

func (r Range) Limit() clause.Expression {
	if len(r) != 2 || r[0] < 0 || r[1] < 0 || r[0] > r[1] {
		return nil
	}

	l := r[1] - r[0] + 1
	return clause.Limit{
		Offset: r[0],
		Limit:  &l,
	}
}

func CheckRange(r Range, total int64) Range {
	if r == nil || len(r) != 2 {
		return nil
	}
	if r[1] > int(total)-1 {
		r[1] = int(total) - 1
	}
	if r[0] < 0 || r[0] > r[1] {
		return nil
	}
	return r
}

// 模型迁移
func ModelMigrate(db *gorm.DB, models ...any) {
	db.AutoMigrate(models...)
}
