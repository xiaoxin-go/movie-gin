package model

import (
    "time"
)

type TActress struct{
    Id int `gorm:"primary_key"json:"id"`
    Name string `gorm:"column:name" json:"name"`
    Age int `gorm:"column:age" json:"age"`
    Height string `gorm:"column:height" json:"height"`
    Cup string `gorm:"column:cup" json:"cup"`
    Birthday string `gorm:"column:birthday" json:"birthday"`
    Url string `gorm:"column:url" json:"url"`
    Image string `gorm:"image" json:"image"`
}

type TFilm struct{
    Id int `gorm:"primary_key" json:"id"`
    Name string `gorm:"column:name" json:"name"`
    Title string `gorm:"column:title" json:"title"`
    ReleaseDate time.Time `gorm:"column:release_date" json:"release_date"`
    Length string `gorm:"column:length" json:"length"`
    Image string `gorm:"column:image" json:"image"`
    Genres []TGenre `gorm:"many2many:t_genre_film;foreignKey:id;joinForeignKey:FilmId;References:id;JoinReferences:GenreId" json:"genres"`
    Actresses []TActress `gorm:"many2many:t_actress_film;foreignKey:id;joinForeignKey:FilmId;References:id;JoinReferences:ActressId" json:"actresses"`
}

type TActressFilm struct{
    Id int `gorm:"primary_key" json:"id"`
    ActressId int `gorm:"column:actress_id" json:"actress_id"`
    FilmId int `gorm:"column:film_id" json:"film_id"`
}

type TLink struct{
    Id int `gorm:"primary_key" json:"id"`
    FilmId int `gorm:"column:film_id" json:"film_id"`
    Magnet string `gorm:"column:magnet" json:"magnet"`
    Name string `gorm:"column:name" json:"name"`
    Size string `gorm:"column:size" json:"size"`
    ShareDate time.Time `gorm:"column:share_date" json:"share_date"`
}

type TGenre struct{
    Id int `gorm:"primary_key"json:"id"`
    Name string `gorm:"column:name" json:"name"`
    Films []TFilm `gorm:"many2many:t_genre_film;foreignKey:id;joinForeignKey:GenreId;References:id;JoinReferences:FilmId"`
}

type TGenreFilm struct{
    Id int `gorm:"primary_key" json:"id"`
    GenreId int `gorm:"genre_id"`
    FilmId int `gorm:"column:film_id" json:"film_id"`
}


type TImage struct{
    Id int `gorm:"primary_key" json:"id"`
    FilmId int `gorm:"column:film_id" json:"film_id"`
    Name string `gorm:"column:name" json:"name"`
    Url string `gorm:"column:url" json:"url"`
    SimpleUrl string `gorm:"column:simple_url" json:"simple_url"`
}


