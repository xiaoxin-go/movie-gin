package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/libs"
	"movie/models"
	"movie/services"
	"strings"
)

type FilmController struct {
	libs.Controller
}

func NewFilmController() libs.Restfuller {
	controller := &FilmController{}
	controller.ModelFunc = func() libs.Instance {
		return new(models.TFilm)
	}
	controller.ListFunc = func() any {
		return new([]*models.TFilm)
	}
	return controller
}

func (*FilmController) Create(ctx *gin.Context) {
	params := &models.TFilm{}
	if e := ctx.ShouldBindJSON(&params); e != nil {
		libs.HttpParamsError(ctx, e.Error())
		return
	}
	params.Name = strings.TrimSpace(params.Name)
	if params.Name == "" {
		libs.HttpParamsError(ctx, "电影名不能为空")
		return
	}
	if e := services.CreateFilm(params.Name); e != nil {
		libs.HttpServerError(ctx, e.Error())
		return
	}
	libs.HttpSuccess(ctx, nil, "创建成功")
}

func GetFilmDetail(ctx *gin.Context) {
	name := ctx.Query("name")
	film, e := services.GetFilmDetail(name)
	if e != nil {
		libs.HttpServerError(ctx, e.Error())
		return
	}
	libs.HttpSuccess(ctx, film, "OK")
}
