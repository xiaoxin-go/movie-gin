package actress

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/actress", handler.List)
	e.GET("/actress/:id", handler.Get)
	e.GET("/actress/:id/film", handler.Film)
	e.DELETE("/actress/:id", handler.Delete)
	e.POST("/actress/:id/films", handler.Films)
}
