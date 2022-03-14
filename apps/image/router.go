package image

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/film/image", handler.List)
}
