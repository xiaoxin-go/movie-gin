package services

import "movie/models"

func GetActressDetail(name string) (*models.TActress, error) {
	result := models.TActress{}
	if e := result.GetDetailByName(name); e != nil {
		return nil, e
	}
	return &result, nil
}
