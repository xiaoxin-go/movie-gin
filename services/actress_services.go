package services

import (
	"fmt"
	"movie/database"
	"movie/models"
)

func GetActressDetail(name string, page, pageSize int) (*models.TActress, int64, error) {
	result := models.TActress{}
	if e := result.FirstByName(name); e != nil {
		return nil, 0, e
	}
	var count int64
	db := database.DB.Model(&models.TActressFilm{}).Where("t_actress_id = ?", result.Id).Count(&count)
	if db.Error != nil {
		return nil, 0, fmt.Errorf("获取演员电影ID失败, err: %w", db.Error)
	}
	if count > 0 {
		filmIds := make([]int, 0)
		if e := db.Offset((page-1)*pageSize).Limit(pageSize).Pluck("t_film_id", &filmIds).Error; e != nil {
			return nil, 0, fmt.Errorf("获取电影信息失败, err: %w", e)
		}
		films, e := new(models.TFilm).FindByIds(filmIds)
		if e != nil {
			return nil, 0, e
		}
		result.Films = films
	}
	return &result, count, nil
}
