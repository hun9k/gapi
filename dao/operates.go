package dao

import (
	"gorm.io/gorm"
)

// 新增
func Insert[M any](o Option, m M) error {
	return gorm.G[M](o.db).
		Create(o.ctx, &m)
}

// 删除
func Del[M any](o Option, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// query
	return gorm.G[M](o.db, clauses...).Delete(o.ctx)
}

// 修改
func Update[M any](o Option, m M, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// query
	return gorm.G[M](o.db, clauses...).Updates(o.ctx, m)
}

func Restore[M any](o Option, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// // unscoped
	// o.db = o.db.Unscoped()
	// // query
	// return gorm.G[M](o.db, clauses...).Update(o.ctx, "deleted_at", nil)
	result := o.db.Model(new(M)).WithContext(o.ctx).Unscoped().Clauses(clauses...).Update("deleted_at", nil)
	return int(result.RowsAffected), result.Error
}

// the rest clauses: ColList, Filter, Order, Limiter
func Select[M any](o Option, conds ...any) ([]M, error) {
	// clauses
	clauses := BuildClauses(conds...)

	// query
	return gorm.G[M](o.db, clauses...).Find(o.ctx)
}

// 统计
func Count[M any](o Option, col string, conds ...any) (int64, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// query
	return gorm.G[M](o.db, clauses...).Count(o.ctx, col)
}
