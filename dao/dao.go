package dao

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"sync"

	"github.com/hun9k/gapi/cache"
	"github.com/hun9k/gapi/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// 执行选项
type Option struct {
	db       *gorm.DB
	ctx      context.Context
	unscoped bool // 是否忽略软删除
}

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

func (f Filter) ID() any {
	if id, ok := f["_id"]; ok {
		return id
	}
	return nil
}

func (f Filter) Where() clause.Expression {
	exprs := []clause.Expression{}
	keyword, kok := f["_keyword"]
	search, sok := f["_search"]
	// do id search
	if id := f.ID(); id != nil {
		exprs = append(exprs, clause.Expr{
			SQL:  "`id` = ?",
			Vars: []any{id},
		})
		delete(f, "_id")
		goto end
	}

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

end:
	if len(exprs) == 0 {
		return nil
	}
	return clause.Where{
		Exprs: exprs,
	}
}

func EmptyFilter() *Filter {
	return &Filter{}
}

// field, asc|desc
type Order []struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}

func (s Order) OrderBy() clause.Expression {
	var s1 = slices.Clone(s)
	columns := []clause.OrderByColumn{}
	for i, v := range s1 {
		if v.Field == "" {
			s = slices.Delete(s, i, i+1)
			continue
		}
		columns = append(columns, clause.OrderByColumn{
			Column: clause.Column{
				Name: v.Field,
			},
			Desc: v.Desc,
		})
	}
	if len(columns) == 0 {
		return nil
	}

	fmt.Println(s, s1)
	return clause.OrderBy{
		Columns: columns,
	}
}

func DefaultOrder() *Order {
	return &Order{
		{"id", true},
	}
}

type Pager map[string]int

const (
	PAGE_KEY = "page"
	SIZE_KEY = "size"
	SIZE_DFT = 12
	SIZE_MAX = 100
)

func (p Pager) Limit() clause.Expression {
	p.Clean()

	_, pOk := p[PAGE_KEY]
	size, sOk := p[SIZE_KEY]
	if !pOk || p[PAGE_KEY] <= 0 {
		p[PAGE_KEY] = 1
	}
	var limit *int
	if !sOk {
		limit = nil
	} else {
		if size <= 0 {
			size = SIZE_DFT
		} else if size > 100 {
			size = SIZE_MAX
		}
		p[SIZE_KEY] = size
		limit = &size
	}

	return clause.Limit{
		Limit:  limit,
		Offset: (p[PAGE_KEY] - 1) * p[SIZE_KEY],
	}
}
func (p Pager) Clean() {
	allKey := map[string]struct{}{PAGE_KEY: {}, SIZE_KEY: {}}
	for k := range p {
		if _, exists := allKey[k]; !exists {
			delete(p, k)
		}
	}
}

func NoLimitPager() *Pager {
	return &Pager{
		PAGE_KEY: 1,
	}
}

// 缓存字段
func CacheFields(cacher cache.Cacher, models ...any) {
	for _, model := range models {
		modelName, fields, err := ModelFields(model)
		if err != nil {
			continue
		}
		cacher.Set(modelName+":fields", fields, cache.NoExpiration)
	}
}

// ModelFields 解析模型的表字段列表
func ModelFields(model any) (string, []string, error) {
	// 解析模型为 schema
	sch, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return "", nil, err
	}

	// 提取列名（数据库表字段名）
	var fields []string
	for _, field := range sch.Fields {
		// 跳过关联字段
		if _, exists := sch.Relationships.Relations[field.Name]; exists {
			continue
		}
		fields = append(fields, field.DBName) // DBName 是数据库列名，Name 是结构体字段名
	}
	return sch.Name, fields, nil
}

func ModelName(model any) (string, error) {
	// 解析模型为 schema
	sch, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return "", err
	}

	return sch.Name, nil
}

// 模型迁移
func ModelMigrate(db *gorm.DB, models ...any) {
	db.AutoMigrate(models...)
}
