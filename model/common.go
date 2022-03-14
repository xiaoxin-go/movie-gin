package model

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
	"time"
)

type Model struct{
	Id int `gorm:"primary_key"json:"id"`
	CreateTime time.Time	`gorm:"column:create_time"json:"create_time"`
	UpdateTime time.Time	`gorm:"column:update_time"json:"update_time"`
}

func (m *Model) BeforeCreate(tx *gorm.DB)(err error){
	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	return
}
func (m *Model) BeforeUpdate(tx *gorm.DB)(err error){
	tx.Statement.SetColumn("UpdateTime", time.Now())
	return
}