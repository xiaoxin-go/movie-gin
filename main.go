package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"movie/config"
	"movie/database"
	"movie/migrate"
	"movie/routes"
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

	// 启动服务
	log.Println("启用服务--->")
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("监听%s\n", addr)
	if e := r.Run(addr); e != nil {
		log.Fatalln(e)
	}
}
