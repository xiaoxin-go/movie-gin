package register

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.POST("/register", register)
}