package routes

import (
	"github.com/gin-gonic/gin"
	"movie/controllers"
)

func FilmRoutes(r *gin.RouterGroup) {
	RegisterRestRoutes(r, "film", controllers.NewFilmController())
	r.GET("film/detail", controllers.GetFilmDetail)
}
