package base

type GetReq struct {
	ID uint `uri:"id" form:"id" json:"id" binding:"required"`
}

type GetRowsReq struct {
	// dao.CommonFilter `form:"filter" json:"filter"`
	// dao.Sort         `form:"sort" json:"sort"`
	// dao.Range        `form:"range" json:"range"`
}

type Resp struct {
	Code     int               `json:"status"`
	Message  string            `json:"message,omitempty"`
	Messages map[string]string `json:"messages,omitempty"`
	Data     any               `form:"data" json:"data,omitempty"`
}
