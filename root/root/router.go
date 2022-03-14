package root

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine){
	e.GET("/publickey", handler.GetPublicKey)
	e.GET("/get_user", handler.GetUser)
	e.GET("/get_menu", handler.GetMenu)
	e.GET("/", handler.Index)
}