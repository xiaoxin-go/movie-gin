package main

import (
	"log"
	"movie/model"
)

func clearActress(){
	actresses := make([]model.TActress, 0)
	db := model.DB.Where("image = ?", "").Find(&actresses)
	if db.Error != nil{
		log.Fatalf("get actresses error => %s\n", db.Error.Error())
	}
	for _, actress := range actresses{
		var count int64
		db = model.DB.Model(&model.TActressFilm{}).Where("actress_id = ?", actress.Id).Count(&count)
		if db.Error != nil{
			continue
		}
		if count < 2{
			log.Println("del actress => ", actress.Name)
			model.DB.Delete(&actress)
		}
	}
}
func clearGenres(){
	genres := make([]model.TGenre, 0)
	db := model.DB.Find(&genres)
	if db.Error != nil{
		log.Fatalf("get actresses error => %s\n", db.Error.Error())
	}
	for _, genre := range genres{
		var count int64
		db = model.DB.Model(&model.TGenreFilm{}).Where("genre_id = ?", genre.Id).Count(&count)
		if db.Error != nil{
			continue
		}
		if count < 1{
			log.Println("del actress => ", genre.Name)
			model.DB.Delete(&genre)
		}
	}
}

func main(){
	clearGenres()
}
