package film

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/film", handler.List)
	e.GET("/film/:id", handler.Get)
}
