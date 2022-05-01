package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	config "movie/conf"
	"movie/model"
	"net/http"
	url2 "net/url"
	"os"
	"time"
)

func InsertFilmData(data FilmData, addActress bool){
	log.Println("insert film data......................")
	log.Println(data.Name)
	film := model.TFilm{
		Name: data.Name,
		Title: data.Title,
		Image: data.ImageUrl,
		ReleaseDate: data.ReleaseDate,
		Length: data.Length,
	}
	if len(data.Actresses) > 3{
		return
	}
	film = insertFilm(film)
	log.Println("filmId==========================", film.Id, data.Name)
	log.Println("insert actress data.....................")
	for _, item := range data.Actresses{
		actressId := insertActress(model.TActress{
			Name: item.Name,
			Url: item.Url,
		}, addActress)
		if actressId > 0{
			insertActressFilm(model.TActressFilm{
				ActressId: actressId,
				FilmId: film.Id,
			})
		}
	}
	log.Println("insert genre data......................")
	for _, name := range data.Genres{
		genreId := insertGenre(model.TGenre{Name: name})
		insertGenreFilm(model.TGenreFilm{
			GenreId: genreId,
			FilmId: film.Id,
		})
	}
	InsertFilmLinks(film.Id, data.Links)
	InsertFilmImages(film.Id, data.Images)
	saveFilmImage(film)
}
func InsertFilmLinks(filmId int, links []Link){
	log.Println("insert links........................")
	filmLinks := make([]model.TLink, 0)
	for _, link := range links{
		filmLinks = append(filmLinks, model.TLink{Name: link.Name, Magnet: link.Magnet,
			Size: link.Size, ShareDate: link.ShareDate, FilmId: filmId})
	}
	if len(filmLinks) > 0{
		model.DB.Create(&filmLinks)
	}
}
func InsertFilmImages(filmId int, images []Image){
	log.Println("insert images........................")
	filmImages := make([]model.TImage, 0)
	for _, image := range images{
		filmImages = append(filmImages, model.TImage{Name: image.Name,
			Url: image.Url, SimpleUrl: image.SimpleUrl, FilmId: filmId})
	}
	model.DB.Create(&filmImages)
}

func insertGenre(genre model.TGenre)int{
	db := model.DB.Model(&model.TGenre{}).Where("name = ?", genre.Name).First(&genre)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		model.DB.Create(&genre)
	}
	return genre.Id
}
func insertActress(actress model.TActress, addActress bool)(result int){
	db := model.DB.Model(&model.TActress{}).Where("name = ?", actress.Name).First(&actress)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		if !addActress{
			return
		}
		model.DB.Create(&actress)
	}
	return actress.Id
}
func insertFilm(film model.TFilm)model.TFilm{
	data := model.TFilm{}
	db := model.DB.Model(&model.TFilm{}).Where("name = ?", film.Name).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		model.DB.Create(&film)
		return film
	}
	return data
}
func insertActressFilm(data model.TActressFilm){
	db := model.DB.Model(&model.TActressFilm{}).Where("actress_id = ? and film_id = ?", data.ActressId, data.FilmId).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		model.DB.Create(&data)
	}
}
func insertGenreFilm(data model.TGenreFilm){
	db := model.DB.Model(&model.TGenreFilm{}).Where("genre_id = ? and film_id = ?", data.GenreId, data.FilmId).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		model.DB.Create(&data)
	}
}

func saveFilmImage(film model.TFilm){
	SaveImage(film.Name, film.Image)
	imageList := make([]model.TImage, 0)
	model.DB.Where("film_id = ?", film.Id).Find(&imageList)
	for index, image := range imageList{
		log.Println(index, len(imageList))
		go SaveImage(image.Name, image.Url)
		go SaveImage(image.Name + "-simple", image.SimpleUrl)
		time.Sleep(200 * time.Millisecond)
	}
}
func SaveImage(name, url string){
	filename := config.Config.ImagePath + name + ".jpg"
	_, err := os.Stat(filename)
	if err == nil{
		fmt.Printf("%s is already exists", name)
		return
	}
	log.Println("get image ", name, url)
	uri, err := url2.Parse(config.Config.Proxy)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: http.ProxyURL(uri),
	}

	client := &http.Client{Transport: tr}
	var resp *http.Response
	flag := 0
	for {
		resp, err = client.Get(url)
		if err != nil{
			if flag >= 3{
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
			flag += 1
			continue
		}
		break
	}

	out, err := os.Create(filename)
	if err != nil{
		log.Println(err)
		return
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil{
		log.Println(err)
	}
}