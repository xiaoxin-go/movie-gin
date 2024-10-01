package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/libs"
	"movie/models"
	"movie/services"
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
	Actress, e := services.GetActressDetail(name)
	if e != nil {
		libs.HttpServerError(ctx, e.Error())
		return
	}
	libs.HttpSuccess(ctx, Actress, "OK")
}
