package main

import (
	"crypto/tls"
	"errors"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"movie/model"
	"movie/utils"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

func main() {
	saveImages()
	return
	service, err := utils.NewService()
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop()

	wd, err := utils.NewWindow()
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Close()
	nameList := getMovies()
	log.Println(nameList)
	for _, name := range nameList{
		log.Println("get name: ", name)
		film := utils.NewFilm(name, wd)
		data := film.Data()
		if film.Error() != nil{
			continue
			log.Fatal(film.Error())
		}
		insertData(data)
	}

	time.Sleep(10 * time.Second)

}

func saveImages(){
	imageList := make([]model.TImage, 0)
	model.DB.Find(&imageList)
	for index, image := range imageList{
		if index < 940{
			continue
		}
		log.Println(index, len(imageList))
		saveImage(image.Name, image.Url)
		saveImage(image.Name + "-simple", image.SimpleUrl)
	}
}
func saveImage(name, url string){
	log.Println("get image ", name, url)
	uri, err := url2.Parse("http://127.0.0.1:7890")
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

	out, err := os.Create("./static/image/" + name + ".jpg")
	if err != nil{
		log.Println(err)
		return
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil{
		log.Println(err)
	}
}


func getMovies()(result []string){
	nameMap := make(map[string]bool)
	files, err := ioutil.ReadDir("E:\\FFOutput\\static")
	if err != nil{
		log.Fatalf("read dir error: %s", err.Error())
	}
	for _, file := range files{
		name := file.Name()
		if _, ok := nameMap[name]; ok{
			continue
		}
		name = strings.ToUpper(strings.Split(name, ".")[0])
		nameMap[name] = true
		data := model.TFilm{}
		db := model.DB.Model(&model.TFilm{}).Where("name = ?", name).First(&data)
		if db.Error == nil{
			continue
		}
		result = append(result, name)
	}
	return
}

func insertData(data utils.FilmData){
	log.Println("insert film data......................")
	log.Println(data)
	log.Println(data.Name)
	film := model.TFilm{
		Name: data.Name,
		Title: data.Title,
		Image: data.ImageUrl,
		ReleaseDate: data.ReleaseDate,
		Length: data.Length,
	}
	filmId := insertFilm(film)
	log.Println("filmId==========================", filmId, data.Name)
	log.Println("insert actress data.....................")
	for _, name := range data.Actresses{
		actressId := insertActress(model.TActress{
			Name: name,
		})
		insertActressFilm(model.TActressFilm{
			ActressId: actressId,
			FilmId: filmId,
		})
	}
	log.Println("insert genre data......................")
	for _, name := range data.Genres{
		genreId := insertGenre(model.TGenre{Name: name})
		insertGenreFilm(model.TGenreFilm{
			GenreId: genreId,
			FilmId: filmId,
		})
	}
	log.Println("insert links........................")
	filmLinks := make([]model.TLink, 0)
	for _, link := range data.Links{
		filmLinks = append(filmLinks, model.TLink{Name: link.Name, Magnet: link.Magnet,
			Size: link.Size, ShareDate: link.ShareDate, FilmId: filmId})
	}
	model.DB.Create(&filmLinks)
	log.Println("insert images........................")
	filmImages := make([]model.TImage, 0)
	for _, image := range data.Images{
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
func insertActress(actress model.TActress)(result int){
	db := model.DB.Model(&model.TActress{}).Where("name = ?", actress.Name).First(&actress)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		model.DB.Create(&actress)
	}
	return actress.Id
}
func insertFilm(film model.TFilm)(result int){
	data := model.TFilm{}
	db := model.DB.Model(&model.TFilm{}).Where("name = ?", film.Name).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		model.DB.Create(&film)
		return film.Id
	}
	return data.Id
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