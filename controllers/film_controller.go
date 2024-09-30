package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/libs"
	"movie/models"
	"movie/services"
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

func GetFilmDetail(ctx *gin.Context) {
	name := ctx.Query("name")
	film, e := services.GetFilmDetail(name)
	if e != nil {
		libs.HttpServerError(ctx, e.Error())
		return
	}
	libs.HttpSuccess(ctx, film, "OK")
}
