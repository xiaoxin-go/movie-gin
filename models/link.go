package models

import (
	"fmt"
	"gorm.io/gorm"
	"movie/database"
	"time"
)

type TLink struct {
	BaseModel
	FilmId    int       `json:"film_id"`
	Name      string    `gorm:"size:50" json:"name"`
	Magnet    string    `gorm:"size:2000" json:"magnet"`
	Size      string    `gorm:"size:50" json:"size"`
	ShareDate time.Time `json:"share_date"`
}

func (t *TLink) BulkCreate(tx *gorm.DB, links []*TLink) error {
	if tx == nil {
		tx = database.DB
	}
	if e := tx.Create(&links).Error; e != nil {
		return fmt.Errorf("create links error:%w", e)
	}
	return nil
}

func (t *TLink) DeleteByFilmId(tx *gorm.DB, filmId int) error {
	if tx == nil {
		tx = database.DB
	}
	if e := tx.Where("film_id = ?", filmId).Delete(t).Error; e != nil {
		return fmt.Errorf("delete links error:%w", e)
	}
	return nil
}
