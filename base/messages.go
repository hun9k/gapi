package base

import "github.com/hun9k/gapi/dao"

type GQuery struct {
	*dao.Filter `form:"filter" binding:""`
	*dao.Order  `form:"order" binding:""`
	*dao.Pager  `form:"pager" binding:""`
}

type Resp struct {
	Error    int               `json:"error"`
	Message  string            `json:"message,omitempty"`
	Messages map[string]string `json:"messages,omitempty"`
	Data     any               `json:"data,omitempty"`
}

type ListData struct {
	List   any        `json:"list"`
	Total  int64      `json:"total"`
	Pager  dao.Pager  `json:"pager"`
	Filter dao.Filter `json:"filter"`
	Order  dao.Order  `json:"order"`
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
