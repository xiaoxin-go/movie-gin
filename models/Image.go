package models

import (
	"fmt"
	"gorm.io/gorm"
	"movie/database"
)

type TImage struct {
	BaseModel
	FilmId int    `json:"film_id"`
	Name   string `gorm:"size:20" json:"name"`
}

func (t *TImage) BulkCreate(tx *gorm.DB, images []*TImage) (err error) {
	if len(images) == 0 {
		return
	}
	if tx == nil {
		tx = database.DB
	}
	if e := tx.Create(&images).Error; e != nil {
		return fmt.Errorf("插入图片失败, err: %w", e)
	}
	return nil
}

func (t *TImage) FindByFilmId(filmId int) ([]*TImage, error) {
	images := make([]*TImage, 0)
	if e := database.DB.Where("film_id = ?", filmId).Find(&images).Error; e != nil {
		return nil, fmt.Errorf("获取图片信息失败, err: %w", e)
	}
	return images, nil
}
