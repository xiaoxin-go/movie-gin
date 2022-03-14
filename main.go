package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	config "movie/conf"
	"movie/libs"
	"movie/routers"
)


func main() {
	flag.Parse()
	r := gin.New()
	gin.Default()
	// 模板渲染
	//r.LoadHTMLGlob("dist/index.html")
	//r.StaticFS("/static", http.Dir("dist/static"))
	// 日志
	err := libs.InitLogger(&config.Config.Log)
	if err != nil{
		panic(err)
	}

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册中间件
	r.Use(libs.GinLogger(), libs.GinRecovery(true))
	r.Use(libs.Cors())

	// 初始化路由配置
	routers.Init(r)

	if err := r.Run(fmt.Sprintf(":%s", config.Config.Port)); err != nil {
		panic(err.Error())
	}
}
