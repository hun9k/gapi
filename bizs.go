package gapi

type QueryStrId struct {
	ID uint `form:"id" json:"id,omitempty"`
}

type QueryStr struct {
	Sort   [2]string `form:"sort" json:"sort,omitempty"`   // [field, method]
	Range  [2]int    `form:"range" json:"range,omitempty"` // [start, end]
	Page   [2]int    `form:"page" json:"page,omitempty"`   // [page, limit]
	Filter struct {
		ID      []uint `form:"id" json:"id,omitempty"`
		Keyword string `form:"keyword" json:"keyword,omitempty"`
	} `form:"filter" json:"filter,omitempty"`
}

const (
	S_OK = iota
	S_ERR
)

const (
	M_OK  = "success"
	M_ERR = "some error"
)

type Resp struct {
	Status  int    `form:"status" json:"status,omitempty"`
	Message string `form:"message" json:"message,omitempty"`
	Data    any    `form:"data" json:"data,omitempty"`
}
