package film

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/film", handler.List)
	e.POST("/film", handler.Post)
	e.GET("/film/:id", handler.Get)
	e.GET("/film/:id/like", handler.Like)
	e.DELETE("/film/:id", handler.Delete)
	e.POST("/film/:id/cover", handler.Cover)
}
