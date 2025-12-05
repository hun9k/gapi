package biz

const (
	S_OK = iota
	S_ERR
)

const (
	M_OK  = "success"
	M_ERR = "some error"
)

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

type Resp struct {
	Code     int               `json:"status"`
	Message  string            `json:"message,omitempty"`
	Messages map[string]string `json:"messages,omitempty"`
	Data     any               `form:"data" json:"data,omitempty"`
}
