package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/dao"
	"github.com/hun9k/gapi/log"
	"gorm.io/gorm"
)

func Get[M any](ctx *gin.Context) {
	// bind request
	req := GQuery{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Info("get bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// id search
	if req.Filter != nil && req.Filter.ID() != nil { // id search
		list, err := dao.Select[M](
			dao.MkOpt(), req.Filter,
		)
		if err != nil {
			log.Info("select error", "path", ctx.Request.URL.Path, "error", err)
			ctx.JSON(http.StatusInternalServerError, Resp{
				Error:   2,
				Message: err.Error(),
			})
		}
		if len(list) == 1 { // id exists
			ctx.JSON(http.StatusOK, Resp{
				Error: 0,
				Data:  list[0],
			})
			return
		} else { // id not exists
			err := gorm.ErrRecordNotFound
			log.Info("selecct error", "path", ctx.Request.URL.Path, "error", err)
			ctx.JSON(http.StatusNotFound, Resp{
				Error:   2,
				Message: err.Error(),
			})
			return
		}
	}

	// get list
	if req.Filter == nil {
		req.Filter = dao.EmptyFilter()
	}
	if req.Pager == nil {
		req.Pager = dao.NoLimitPager()
	}
	if req.Order == nil {
		req.Order = dao.DefaultOrder()
	}
	list, err := dao.Select[M](
		dao.MkOpt(),
		req.Filter, req.Order, req.Pager,
	)
	if err != nil {
		log.Info("select error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// list search
	total, err := dao.Count[M](
		dao.MkOpt(), "*",
		req.Filter,
	)
	if err != nil {
		log.Info("count error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, Resp{
		Error: 0,
		Data: ListData{
			List:   list,
			Total:  total,
			Order:  *req.Order,
			Pager:  *req.Pager,
			Filter: *req.Filter,
		},
	})
}

func Post[M any](ctx *gin.Context) {
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
	if err := dao.Insert(dao.MkOpt(), item); err != nil {
		log.Info("insert error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, Resp{
		Error: 0,
		Data:  item,
	})
}

func Put[M any](ctx *gin.Context) {
	// bind query
	req := GQuery{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// bind body
	item := new(M)
	if err := ctx.ShouldBind(&item); err != nil {
		log.Info("bind body error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// update rows
	num, err := dao.Update(dao.MkOpt(), item, req.Filter, req.Order, req.Pager)
	if err != nil {
		log.Info("update error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, Resp{
		Message: "data is num of updated rows",
		Data:    num,
	})
}

func Delete[M any](ctx *gin.Context) {
	// bind request
	req := GQuery{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// delete rows
	num, err := dao.Del[M](dao.MkOpt(), req.Filter, req.Order, req.Pager)
	if err != nil {
		log.Info("delete error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, Resp{
			Error:   2,
			Message: err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, Resp{
		Message: "data is num of deleted rows",
		Data:    num,
	})
}

func Restore[M any](ctx *gin.Context) {
	// bind query
	req := GQuery{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Info("bind query error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, Resp{
			Error:   1,
			Message: err.Error(),
		})
		return
	}

	// restore rows
	num, err := dao.Restore[M](dao.MkOpt(), req.Filter)
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
		Data:    num,
	})
}
