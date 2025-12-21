package base

import "github.com/hun9k/gapi/dao"

type GID struct {
	ID uint `uri:"id" binding:"required"`
}

type GIDs struct {
	IDs    []uint  `form:"id" binding:"required"`
	Option *string `form:"option" binding:""`
}

type ListQuery struct {
	Filter *string `form:"filter" binding:""`
	Sort   *string `form:"sort" binding:""`
	Range  *string `form:"range" binding:""`
}

type GQuery struct {
	dao.Filter `form:"filter" binding:""`
	dao.Sorts  `form:"order" binding:""`
	dao.Range  `form:"range" binding:""`
}

type Resp struct {
	Error   int    `json:"error"`
	Message any    `json:"message,omitempty"`
	Item    any    `json:"item,omitempty"`
	List    any    `json:"list,omitempty"`
	Total   *int64 `json:"total,omitempty"`
	Num     *int   `json:"num,omitempty"`
}

// // offset, limit
// type Limit [2]int

// func (l Limit) OffsetLimit() (int, int) {
// 	return l[0], l[1]
// }

// // []field, asc|desc
// type Sorts []Sort

// func (s Sorts) Order() string {
// 	orders := []string{}
// 	for _, sort := range s {
// 		orders = append(orders, sort.Order())
// 	}
// 	return strings.Join(orders, ", ")
// }

// // start, end
// type Range [2]int

// func (r Range) OffsetLimit() (int, int) {
// 	return r[0], r[1] - r[0] + 1
// }
