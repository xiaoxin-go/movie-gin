package models

import (
	"fmt"
	"gorm.io/gorm"
	"movie/database"
	"time"
)

type TFilm struct {
	BaseModel
	Name        string      `gorm:"size:20;unique" json:"name"`
	Title       string      `gorm:"size:2000" json:"title"`
	ReleaseDate time.Time   `json:"release_date"`
	Length      string      `gorm:"size:20" json:"length"`
	Actresses   []*TActress `gorm:"many2many:t_actress_film;" json:"actresses"`
	Images      []*TImage   `gorm:"foreignKey:film_id" json:"images"`
	Links       []*TLink    `gorm:"foreignKey:film_id" json:"links"`
}

func (t *TFilm) FirstByName(name string) error {
	if e := database.DB.Where("name = ?", name).First(t).Error; e != nil {
		return fmt.Errorf("获取电影信息失败, name: %s, err: %w", name, e)
	}
	return nil
}
func (t *TFilm) FindByIds(ids []int) ([]*TFilm, error) {
	result := make([]*TFilm, 0)
	if e := database.DB.Where("id in ?", ids).Find(&result).Error; e != nil {
		return nil, fmt.Errorf("获取电影信息失败, err: %w", e)
	}
	return result, nil
}
func (t *TFilm) GetDetailByName(name string) error {
	if e := database.DB.Where("name = ?", name).Preload("Actresses").Preload("Images").Preload("Links").First(t).Error; e != nil {
		return fmt.Errorf("获取电影信息失败, name: %s, err: %w", name, e)
	}
	return nil
}
func (t *TFilm) Create(tx *gorm.DB) error {
	if tx == nil {
		tx = database.DB
	}
	if e := tx.Create(t).Error; e != nil {
		return fmt.Errorf("创建电影失败, err: %w", e)
	}
	return nil
}
