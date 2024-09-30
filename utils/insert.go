package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"movie/database"
	model "movie/models"
	"net/http"
	"net/url"
	"os"
	"time"
)

func InsertFilmData(data FilmData, addActress bool) {
	log.Println("insert film data......................")
	log.Println(data.Name)
	film := model.TFilm{
		Name:        data.Name,
		Title:       data.Title,
		ReleaseDate: data.ReleaseDate,
		Length:      data.Length,
	}
	if len(data.Actresses) > 3 {
		return
	}
	film = insertFilm(film)
	log.Println("filmId==========================", film.Id, data.Name)
	log.Println("insert actress data.....................")
	for _, item := range data.Actresses {
		actressId := insertActress(model.TActress{
			Name: item.Name,
		}, addActress)
		if actressId > 0 {
			insertActressFilm(model.TActressFilm{
				ActressId: actressId,
				FilmId:    film.Id,
			})
		}
	}
	log.Println("insert genre data......................")
	for _, name := range data.Genres {
		genreId := insertGenre(model.TGenre{Name: name})
		insertGenreFilm(model.TGenreFilm{
			GenreId: genreId,
			FilmId:  film.Id,
		})
	}
	InsertFilmLinks(film.Id, data.Links)
	InsertFilmImages(film.Id, data.Images)
	saveFilmImage(film, "")
}
func InsertFilmLinks(filmId int, links []Link) {
	log.Println("insert links........................")
	filmLinks := make([]model.TLink, 0)
	for _, link := range links {
		filmLinks = append(filmLinks, model.TLink{Name: link.Name, Magnet: link.Magnet,
			Size: link.Size, ShareDate: link.ShareDate, FilmId: filmId})
	}
	if len(filmLinks) > 0 {
		database.DB.Create(&filmLinks)
	}
}
func InsertFilmImages(filmId int, images []Image) {
	log.Println("insert images........................")
	filmImages := make([]model.TImage, 0)
	for _, image := range images {
		filmImages = append(filmImages, model.TImage{Name: image.Name, FilmId: filmId})
	}
	database.DB.Create(&filmImages)
}

func insertGenre(genre model.TGenre) int {
	db := database.DB.Model(&model.TGenre{}).Where("name = ?", genre.Name).First(&genre)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		database.DB.Create(&genre)
	}
	return genre.Id
}
func insertActress(actress model.TActress, addActress bool) (result int) {
	db := database.DB.Model(&model.TActress{}).Where("name = ?", actress.Name).First(&actress)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		if !addActress {
			return
		}
		database.DB.Create(&actress)
	}
	return actress.Id
}
func insertFilm(film model.TFilm) model.TFilm {
	data := model.TFilm{}
	db := database.DB.Model(&model.TFilm{}).Where("name = ?", film.Name).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		database.DB.Create(&film)
		return film
	}
	return data
}
func insertActressFilm(data model.TActressFilm) {
	db := database.DB.Model(&model.TActressFilm{}).Where("actress_id = ? and film_id = ?", data.ActressId, data.FilmId).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		database.DB.Create(&data)
	}
}
func insertGenreFilm(data model.TGenreFilm) {
	db := database.DB.Model(&model.TGenreFilm{}).Where("genre_id = ? and film_id = ?", data.GenreId, data.FilmId).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		database.DB.Create(&data)
	}
}

func saveFilmImage(film model.TFilm, url string) {
	SaveImage(film.Name, url)
	imageList := make([]model.TImage, 0)
	database.DB.Where("film_id = ?", film.Id).Find(&imageList)
	for index, image := range imageList {
		log.Println(index, len(imageList))
		go SaveImage(image.Name, "image.Url")
		go SaveImage(image.Name+"-simple", "image.SimpleUrl")
		time.Sleep(200 * time.Millisecond)
	}
}
func SaveImage(name, url1 string) error {
	filename := "f:\\static\\images" + name + ".jpg"
	_, err := os.Stat(filename)
	if err == nil {
		fmt.Printf("%s is already exists", name)
		return nil
	}
	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.javbus.com")
	req.Header.Set("Host", "https://www.javbus.com")
	log.Println("get image ", name, url1)
	uri, err := url.Parse("127.0.0.1:10809")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(uri),
	}

	client := &http.Client{Transport: tr, Timeout: time.Second * 3}
	var resp *http.Response
	flag := 0
	for {
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("get image %s error: %v", name, err)
			if flag >= 3 {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
			flag += 1
			continue
		}
		break
	}

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file error: %w", err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}
