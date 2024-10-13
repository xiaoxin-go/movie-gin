package models

type TActressFilm struct {
	ActressId int `json:"t_actress_id"`
	FilmId    int `json:"t_film_id"`
}

func (TActressFilm) TableName() string {
	return "t_actress_film"
}
