package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"movie/config"
	"movie/database"
	"movie/migrate"
	"movie/routes"
	"net/http"
	"os"
)

func main() {
	log.Println("启动服务--->")
	r := gin.Default()

	// 加载配置
	log.Println("加载配置--->")
	config.LoadConfig("config.json")

	// 连接数据库
	log.Println("连接数据库--->")
	database.ConnectDatabase(config.AppConfig.Database)

	// 注册表
	log.Println("初始化表结构--->")
	migrate.AutoMigrate()

	// 设置路由
	log.Println("注册路由--->")
	routes.SetupRoutes(r)

	if f, err := os.Stat("dist/"); err == nil && f.IsDir() {
		// 模板渲染
		r.LoadHTMLGlob("dist/index.html")
		// 静态目录
		r.StaticFS("/css", http.Dir("dist/css"))
		r.StaticFS("/js", http.Dir("dist/js"))
		r.StaticFS("/img", http.Dir("dist/img"))
	}

	// 启动服务
	log.Println("启用服务--->")
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("监听%s\n", addr)
	if e := r.Run(addr); e != nil {
		log.Fatalln(e)
	}
}
