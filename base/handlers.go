package base

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/dao"
	"github.com/hun9k/gapi/log"
)

// 创建
func Create[M any](ctx *gin.Context) {
	// bind request
	item := new(M)
	if err := ctx.ShouldBind(&item); err != nil {
		log.Info("get bind body error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// insert row
	if err := dao.InsertRow(dao.MkOpt(), item); err != nil {
		log.Info("insert error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, item)
}

// 删除单条记录，利用ID
func Delete[M any](ctx *gin.Context) {
	// bind uri
	req := IDUri{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Info("bind ID error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// delete item
	num, err := dao.DelRow[M](dao.MkOpt(), req.ID)
	if err != nil {
		log.Info("delete error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	item := map[string]uint{}
	if num == 1 {
		item = map[string]uint{"id": req.ID}
	}
	ctx.JSON(http.StatusOK, item)
}

// 删除列表，利用ID列表
func DeleteMany[M any](ctx *gin.Context) {
	// bind request
	req := Cond{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// delete rows
	ids := dao.FilterIDs(req.Filter)
	if _, err := dao.DelRows[M](dao.MkOpt(), ids); err != nil {
		log.Info("delete error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	list := make([]map[string]any, len(ids))
	for i, id := range ids {
		list[i] = map[string]any{"id": id}
	}
	ctx.JSON(http.StatusOK, list)
}

// 更新单条记录，利用ID
func Update[M any](ctx *gin.Context, model M, cols []string) {
	// bind ID
	req := IDUri{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Info("bind uri error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			List: []uint{1, 2, 3},
		})
	}

	// update row
	if _, err := dao.UpdateRow(dao.MkOpt(), model, req.ID, cols); err != nil {
		log.Info("update error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, model)
}

// 更新多条记录，利用ID列表
func UpdateMany[M any](ctx *gin.Context, model M, cols []string) {
	// bind query
	req := Cond{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// update rows
	ids := dao.FilterIDs(req.Filter)
	if _, err := dao.UpdateRows(dao.MkOpt(), model, ids, cols); err != nil {
		log.Info("update error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	list := make([]map[string]any, len(ids))
	for i, id := range ids {
		list[i] = map[string]any{"id": id}
	}
	ctx.JSON(http.StatusOK, list)
}

func Restore[M any](ctx *gin.Context) {
	// bind query
	req := Cond{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// restore rows
	num, err := dao.RestoreRow[M](dao.MkOpt(), req.Filter)
	if err != nil {
		log.Info("restore error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, Resp{
		Message: "data is num of restored rows",
		Item:    num,
	})
}

// 获取单条记录，利用ID
func GetOne[M any](ctx *gin.Context) {
	// bind uri
	req := IDUri{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Info("bind ID error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// get item
	item, err := dao.SelectRow[M](
		dao.MkOpt(), req.ID,
	)
	if err != nil {
		log.Info("select error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
	}

	// response
	ctx.JSON(http.StatusOK, item)
}

func Get[M any](ctx *gin.Context) {
	if ids, ok := ctx.GetQueryArray("id"); ok && len(ids) > 0 {
		// id exists in query
		GetMany[M](ctx)
	} else {
		// id isn't exists
		GetList[M](ctx)
	}
}

// 获取列表，利用ID列表，无翻页，ID倒序
func GetMany[M any](ctx *gin.Context) {
	// bind request
	req := Cond{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// select by IDs
	list, err := dao.SelectRows[M](
		dao.MkOpt(), dao.FilterIDs(req.Filter),
	)
	if err != nil {
		log.Info("select error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
	}

	// response
	ctx.JSON(http.StatusOK, list)
}

func ShouldBindListQuery(ctx *gin.Context) (*Cond, error) {
	req := ListQuery{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, err
	}
	// marshal filter，sort, range
	filter := dao.Filter{}
	if err := json.Unmarshal([]byte(*req.Filter), &filter); err != nil {
		return nil, err
	}
	sort := dao.Sort{}
	if err := json.Unmarshal([]byte(*req.Sort), &sort); err != nil {
		return nil, err
	}
	rangee := dao.Range{}
	if err := json.Unmarshal([]byte(*req.Range), &rangee); err != nil {
		return nil, err
	}

	return &Cond{
		Filter: filter,
		Sorts:  dao.Sorts{sort},
		Range:  rangee,
	}, nil
}
func GetList[M any](ctx *gin.Context) {
	// bind request
	req, err := ShouldBindListQuery(ctx)
	if err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// get total
	total, err := dao.Count[M](
		dao.MkOpt(), "*", req.Filter,
	)
	if err != nil {
		log.Info("count error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// get list
	list, err := dao.Select[M](
		dao.MkOpt(),
		dao.CheckFilter(req.Filter),
		dao.CheckSort(req.Sorts),
		dao.CheckRange(req.Range, total),
	)
	if err != nil {
		log.Info("select error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.Header("Content-Range", fmt.Sprintf("items %d-%d/%d", req.Range[0], req.Range[1], total))
	ctx.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
