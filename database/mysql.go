package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"movie/config"
	"time"
)

var DB *gorm.DB

func ConnectDatabase(conf config.Mysql) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name) // 修改为你的数据库 DSN
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 设置连接池参数
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)                  // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 设置最大打开连接数
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 设置空闲连接最大生命周期
	sqlDB.SetConnMaxLifetime(1 * time.Hour)    // 设置连接最大生命周期

	log.Println("Database connected")
}
