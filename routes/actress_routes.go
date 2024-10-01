package routes

import (
	"github.com/gin-gonic/gin"
	"movie/controllers"
)

func ActressRoutes(r *gin.RouterGroup) {
	RegisterRestRoutes(r, "actress", controllers.NewActressController())
	r.GET("actress/detail", controllers.GetActressDetail)
}
