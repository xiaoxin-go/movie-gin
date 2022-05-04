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
	e.POST("/film/:id/link", handler.AddLink)
	e.POST("/film/:id/image", handler.AddImage)
	e.PUT("/film/:id/image", handler.SaveImage)
	e.POST("/film/:id/collect", handler.Collect)
	e.POST("/film/:id/uncollect", handler.UnCollect)
	e.GET("/film/:id/iscollect", handler.IsCollect)
	e.GET("/film/:id/isplayer", handler.IsPlayer)
	e.GET("/film/player", handler.Player)
}
