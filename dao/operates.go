package dao

import (
	"gorm.io/gorm"
)

// 新增
func InsertRow[M any](o Option, m M) error {
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

func DelRow[M any](o Option, id uint) (rowsAffected int, err error) {
	return gorm.G[M](o.db,
		ID(id).Where(),
	).Delete(o.ctx)
}

func DelRows[M any](o Option, ids []any) (rowsAffected int, err error) {
	if ids == nil {
		return 0, nil
	}
	return gorm.G[M](o.db,
		IDs(ids).Where(),
	).Delete(o.ctx)
}

// 修改
func Update[M any](o Option, m M, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// query
	return gorm.G[M](o.db, clauses...).Updates(o.ctx, m)
}
func UpdateRow[M any](o Option, m M, id uint, cols []string) (int, error) {
	return gorm.G[M](o.db,
		ID(id).Where(),
		Cols(cols).Select(),
	).Updates(o.ctx, m)
}

func UpdateRows[M any](o Option, m M, ids []any, cols []string) (int, error) {
	return gorm.G[M](o.db,
		IDs(ids).Where(),
		Cols(cols).Select(),
	).Updates(o.ctx, m)
}

func RestoreRow[M any](o Option, conds ...any) (int, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// // unscoped
	// o.db = o.db.Unscoped()
	// // query
	// return gorm.G[M](o.db, clauses...).Update(o.ctx, "deleted_at", nil)
	result := o.db.Model(new(M)).WithContext(o.ctx).Unscoped().Clauses(clauses...).Update("deleted_at", nil)
	return int(result.RowsAffected), result.Error
}

// 过滤器，排序，翻页查询 clauses: ColList, Filter, Order, Limiter
func Select[M any](o Option, conds ...any) ([]M, error) {
	// clauses
	clauses := BuildClauses(conds...)

	// query
	return gorm.G[M](o.db, clauses...).Find(o.ctx)
}

// 查单项
func SelectRow[M any](o Option, id uint) (M, error) {
	return gorm.G[M](o.db,
		ID(id).Where(),
	).First(o.ctx)
}

// 查列表
func SelectRows[M any](o Option, ids []any) ([]M, error) {
	// 默认ID倒序
	sort := Sorts{{"id", "DESC"}}

	return gorm.G[M](o.db,
		IDs(ids).Where(),
		sort.OrderBy(),
	).Find(o.ctx)
}

// 统计
func Count[M any](o Option, col string, conds ...any) (int64, error) {
	// clauses
	clauses := BuildClauses(conds...)
	// query
	return gorm.G[M](o.db, clauses...).Count(o.ctx, col)
}
