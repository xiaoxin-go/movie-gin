package models

type TActressFilm struct {
	BaseModel
	ActressId int `json:"actress_id"`
	FilmId    int `json:"film_id"`
}
