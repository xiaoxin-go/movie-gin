package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/libs"
	"movie/models"
	"movie/services"
	"strconv"
)

type ActressController struct {
	libs.Controller
}

func NewActressController() libs.Restfuller {
	controller := &ActressController{}
	controller.ModelFunc = func() libs.Instance {
		return new(models.TActress)
	}
	controller.ListFunc = func() any {
		return new([]*models.TActress)
	}
	return controller
}

func GetActressDetail(ctx *gin.Context) {
	name := ctx.Query("name")
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 20
	}
	Actress, count, e := services.GetActressDetail(name, page, size)
	if e != nil {
		libs.HttpServerError(ctx, e.Error())
		return
	}
	libs.HttpListSuccess(ctx, Actress, count)
}
