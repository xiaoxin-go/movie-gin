package main

import (
	"fmt"
	"log"
	"movie/model"
	"movie/utils"
	"time"
)

func main(){
	//insertActress()
	//insertFilm()
	//insertActressFilm()
	insertMovieImages()
	//insertLink()
	//insertImage()
}

func insertMovieImages(){
	films := make([]model.TFilm, 0)
	db := model.DB.Find(&films)
	if db.Error != nil{
		log.Fatal(db.Error.Error())
	}
	ch := make(chan string, 10)
	for _, film := range films{
		go func(film model.TFilm, ch chan string) {
			fmt.Printf("save image -> %s\n", film.Name)
			utils.SaveImage(film.Name, film.Image)
			ch <- film.Name
		}(film, ch)
	}
	for{
		select{
		case <-ch:
			//fmt.Println("end-----------", name)

		}

	}

	time.Sleep(10*time.Minute)
}

func insertActressFilm(){
	datas := []model.TActressFilm{
		{ActressId: 511, FilmId: 1984},
		{ActressId: 511, FilmId: 1985},
		{ActressId: 511, FilmId: 1986},
		{ActressId: 511, FilmId: 1987},
		{ActressId: 511, FilmId: 1988},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}

func insertImage(){
	datas := []model.TImage{
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}
func insertLink(){
	datas := []model.TLink{
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}

func insertFilm(){
	datas := []model.TFilm{
		{Name: "O型西瓜乳保姆", Title: "家里的保姆家务工作不好，但是经常用大大的胸部和臀部诱惑自己",  Length: "29", ReleaseDate: utils.StrToDate("2020-01-01")},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}

func insertActress(){
	datas := []model.TActress{
		{Name: "张三", Age: 31, Height: "160cm", Cup: "C"},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}
