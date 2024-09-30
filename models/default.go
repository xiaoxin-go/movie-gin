package models

import "time"

type BaseModel struct {
	Id        int       `gorm:"primary key" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (b *BaseModel) GetId() int {
	return b.Id
}
