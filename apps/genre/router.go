package genre

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/genre", handler.List)
	e.GET("/genre/:id", handler.Get)
	e.GET("/genre/:id/film", handler.Film)
}
