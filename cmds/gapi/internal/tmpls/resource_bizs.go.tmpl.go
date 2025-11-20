package tmpls

var Resource_bizs = `package {{.Resource}}

import (
	"github.com/hun9k/gapi"
)

// 创建一个资源
type ReqPost struct {
}
type RespPost struct {
	gapi.Resp
}

// 删除一个资源
type ReqDeleteId struct {
	gapi.ReqId
}
type RespDeleteId struct {
	gapi.Resp
}

// 删除多个资源
type ReqDelete struct {
	gapi.Req
}
type RespDelete struct {
	gapi.Resp
}

// 恢复一个资源
type ReqPatchIdRestore struct {
	gapi.ReqId
}
type RespPatchIdRestore struct {
	gapi.Resp
}

// 恢复多个资源
type ReqPatchRestore struct {
	gapi.Req
}
type RespPatchRestore struct {
	gapi.Resp
}

// 更新一个资源
type ReqPutId struct {
	gapi.ReqId
}
type RespPutId struct {
	gapi.Resp
}

// 更新多个资源
type ReqPut struct {
	gapi.Req
}
type RespPut struct {
	gapi.Resp
}

// 更新一个资源的部分字段
type ReqPatchId struct {
	gapi.ReqId
}
type RespPatchId struct {
	gapi.Resp
}

// 更新多个资源的部分字段
type ReqPatch struct {
	gapi.Req
}
type RespPatch struct {
	gapi.Resp
}

// 获取一个资源
type ReqGetId struct {
	gapi.ReqId
}
type RespGetId struct {
	gapi.Resp
}

// 获取多个资源
type ReqGet struct {
	gapi.Req
}
type RespGet struct {
	gapi.Resp
}

`
