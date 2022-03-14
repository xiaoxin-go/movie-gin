package link

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/film/link", handler.List)
}
