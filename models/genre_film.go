package models

type TGenreFilm struct {
	BaseModel
	GenreId int `json:"genre_id"`
	FilmId  int `json:"film_id"`
}
