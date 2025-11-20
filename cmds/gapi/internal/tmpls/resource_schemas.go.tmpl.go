package tmpls

var Resource_schemas = `package schemas

import (
	"github.com/hun9k/gapi"
)

type {{.Schema.Name}} struct {
	// 自定义字段

	// 嵌入基础字段，ID，CreatedAt，UpdatedAt, DeletedAt
	gapi.Model
}

`
