package models

import (
	"fmt"
	"movie/database"
)

type TActress struct {
	BaseModel
	Name     string   `gorm:"size:50;unique" json:"name"`
	Height   string   `gorm:"size:50" json:"height"`
	Cup      string   `gorm:"size:50" json:"cup"`
	Birthday string   `gorm:"size:50" json:"birthday"`
	Films    []*TFilm `gorm:"many2many:t_actress_film;" json:"films"`
}

func (t *TActress) GetDetailByName(name string) error {
	if e := database.DB.Where("name = ?", name).Preload("Films").First(t).Error; e != nil {
		return fmt.Errorf("获取演员信息失败, name: %s, err: %w", name, e)
	}
	return nil
}
