package tmpls

var Resource_bizs = `package {{.Resource}}

import (
	"github.com/hun9k/gapi"
)

// 创建资源
type PostBody struct {
	{{.ResourceBody}}
}

// 删除一个资源
type DeleteIdQuery struct {
	gapi.QueryStrId
}

// 删除多个资源
type DeleteQuery struct {
	gapi.QueryStr
}

// 恢复一个资源
type RestoreIdQuery struct {
	gapi.QueryStrId
}

// 恢复多个资源
type RestoreQuery struct {
	gapi.QueryStr
}

// 更新一个资源
type PutIdQuery struct {
	gapi.QueryStrId
}
type PutIdBody struct {
	{{.ResourceBody}}
}

// 更新多个资源
type PutQuery struct {
	gapi.QueryStr
}
type PutBody struct {
	{{.ResourceBody}}
}

// 更新一个资源的部分字段
type PatchIdQuery struct {
	gapi.QueryStrId
}
type PatchIdBody struct {
	{{.ResourcePatchBody}}
}

// 更新多个资源的部分字段
type PatchQuery struct {
	gapi.QueryStr
}
type PatchBody struct {
}

// 获取一个资源
type GetIdQuery struct {
	gapi.QueryStrId
}

// 获取多个资源
type GetQuery struct {
	gapi.QueryStr
}

`
