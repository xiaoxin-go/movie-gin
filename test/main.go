package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"log"
	"movie/model"
	"movie/utils"
	"strings"
)

func main() {
	//getActressUrl()
	//return
	//return
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
	//nameList := getMovies()
	//nameList := []string{"FSDSS-373"}
	//saveFilms(nameList, wd)
	actressList := getActresses()
	for _, actress := range actressList{
		log.Printf("actress => %+v", actress)
		actress.Url = "https://www.javbus.com/star/okq"
		controller := utils.NewActressController(actress.Url, wd)
		actressData := controller.Data()
		fmt.Println("-------------------------")
		fmt.Println(actressData)
		fmt.Println(len(actressData.Films))
		break
		if controller.Error() != nil{
			log.Fatal(controller.Error().Error())
		}
		fmt.Printf("%+v\n", actressData)
		actress.Age = actressData.Age
		actress.Birthday = actressData.Birthday
		actress.Cup = actressData.Cup
		actress.Height = actressData.Height
		model.DB.Save(&actress)
	}
}

func saveActressUrl(){
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
	actresses := make([]model.TActress, 0)
	db := model.DB.Where("url is null").Find(&actresses)
	if db.Error != nil{
		log.Fatalf("get actress error, %s", db.Error.Error())
	}
	for _, actress := range actresses{
		log.Println("actress: ", actress.Name)
		actressFilm := model.TActressFilm{}
		db = model.DB.Where("actress_id = ?", actress.Id).First(&actressFilm)
		if db.Error != nil{
			fmt.Println("get actress error---------")
			continue
			log.Fatalf("get actress error, %s", db.Error.Error())
		}
		film := model.TFilm{}
		db = model.DB.Where("id = ?", actressFilm.FilmId).First(&film)
		if db.Error != nil{
			model.DB.Delete(&model.TActressFilm{}, actressFilm.Id)
			continue
			log.Fatalf("get film error, %s", db.Error.Error())
		}
		filmData := utils.NewFilm(film.Name, wd)
		if filmData.Error() != nil{
			continue
			//log.Fatal(film.Error())
		}
		nameUrlMap := map[string]string{}
		for _, item := range filmData.Data().Actresses{
			nameUrlMap[item.Name] = item.Url
		}
		actress.Url = nameUrlMap[actress.Name]
		fmt.Println("save-------------", actress.Url)
		model.DB.Save(&actress)
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

func getActresses()(result []model.TActress){
	result = make([]model.TActress, 0)
	db := model.DB.Where("url is not null").Find(&result)
	if db.Error != nil{
		log.Fatalf("find actress error => %s", db.Error.Error())
	}
	return
}

func saveFilms(names []string, wd selenium.WebDriver){
	for _, name := range names{
		log.Println("get name ===> ", name)
		film := utils.NewFilm(name, wd)
		data := film.Data()
		if film.Error() != nil{
			log.Fatal(film.Error())
		}
		utils.InsertFilmData(data)
	}
}