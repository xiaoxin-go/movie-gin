package services

import "movie/models"

func GetFilmDetail(name string) (*models.TFilm, error) {
	film := models.TFilm{}
	if e := film.GetDetailByName(name); e != nil {
		return nil, e
	}
	return &film, nil
}
