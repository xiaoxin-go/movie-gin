package actress

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/actress", handler.List)
	e.GET("/actress/:id", handler.Get)
}
