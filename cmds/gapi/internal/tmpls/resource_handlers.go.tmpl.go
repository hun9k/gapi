package tmpls

var Resource_handlers = `package {{.Resource}}

import "github.com/gin-gonic/gin"

// 创建一个资源
func Post(ctx *gin.Context) {

}

// 删除一个资源
func DeleteId(ctx *gin.Context) {

}

// 删除多个资源
func Delete(ctx *gin.Context) {

}

// 恢复一个资源
func PatchIdRestore(ctx *gin.Context) {

}

// 恢复多个资源
func PatchRestore(ctx *gin.Context) {

}

// 更新一个资源
func PutId(ctx *gin.Context) {

}

// 更新多个资源
func Put(ctx *gin.Context) {

}

// 更新一个资源的部分字段
func PatchId(ctx *gin.Context) {

}

// 更新多个资源的部分字段
func Patch(ctx *gin.Context) {

}

// 获取一个资源
func GetId(ctx *gin.Context) {

}

// 获取多个资源
func Get(ctx *gin.Context) {

}

`
