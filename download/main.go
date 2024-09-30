package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"movie/config"
	"movie/database"
	"movie/models"
	"movie/utils"
)

func main() {
	// 下载电影
	log.Println("加载配置--->")
	config.LoadConfig("config.json")

	// 连接数据库
	log.Println("连接数据库--->")
	database.ConnectDatabase(config.AppConfig.Database)
	// 1. 获取电影详情
	service, err := utils.NewService()
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop()
	wd, e := utils.NewWindow()
	if e != nil {
		log.Fatalln(e.Error())
	}
	log.Println("获取电影信息----->")
	film := utils.NewFilm("ABF-073", wd)
	if e := film.Error(); e != nil {
		log.Fatalln(e.Error())
	}
	log.Println("保存电影信息--->")
	f, e := saveFilm(film.Data())
	if e != nil {
		log.Fatalln(e.Error())
	}

	fmt.Println(film.Error())
	log.Println("保存链接--->")
	if e := saveLinks(f.Id, film.Data().Links); e != nil {
		log.Fatalln(e.Error())
	}
	log.Println("保存logo")
	if e := utils.SaveImage("\\logo\\"+film.Name, film.Data().ImageUrl); e != nil {
		log.Println("保存logo失败:", e.Error())
	}

	log.Println("保存图片--->")
	if e := saveImages(f.Id, film.Data().Images); e != nil {
		log.Fatalln(e.Error())
	}
	log.Println("保存演员信息--->")
	actress := saveActress(f, film.Data().Actresses, wd)
	if e1 := database.DB.Model(f).Association("Actresses").Replace(actress); e1 != nil {
		log.Fatalln("创建关联失败: ", e1.Error())
	}
}

func saveActress(film *models.TFilm, actresses []utils.Actress, wd selenium.WebDriver) []*models.TActress {
	result := make([]*models.TActress, 0)
	for _, v := range actresses {
		c := utils.NewActressController(v.Url, wd)
		d := c.Data()
		if c.Error() != nil {
			log.Println("获取演员信息失败: ", c.Error())
			continue
		}
		actress := models.TActress{
			Name:     d.Name,
			Cup:      d.Cup,
			Height:   d.Height,
			Birthday: d.Birthday,
			Films:    []*models.TFilm{film},
		}
		if e := database.DB.FirstOrCreate(&actress).Error; e != nil {
			log.Println("创建获取演员信息失败: ", c.Error())
			continue
		}
		result = append(result, &actress)
	}
	return result
}

func saveFilm(f utils.FilmData) (*models.TFilm, error) {
	film := models.TFilm{}
	if e := film.FirstByName(f.Name); e == nil {
		return &film, nil
	}
	film.Name = f.Name
	film.Title = f.Title
	film.Length = f.Length
	film.ReleaseDate = f.ReleaseDate
	if e := film.Create(nil); e != nil {
		return nil, e
	}
	return &film, nil
}

func saveLinks(filmId int, data []utils.Link) error {
	if e := new(models.TLink).DeleteByFilmId(nil, filmId); e != nil {
		return e
	}
	links := make([]*models.TLink, 0)
	for _, v := range data {
		links = append(links, &models.TLink{
			FilmId:    filmId,
			Name:      v.Name,
			Magnet:    v.Magnet,
			Size:      v.Size,
			ShareDate: v.ShareDate,
		})
	}
	if e := new(models.TLink).BulkCreate(nil, links); e != nil {
		return e
	}
	return nil
}

func saveImages(filmId int, data []utils.Image) error {
	currentImages, e := new(models.TImage).FindByFilmId(filmId)
	if e != nil {
		return e
	}
	currentImageM := make(map[string]*models.TImage)
	for _, v := range currentImages {
		currentImageM[v.Name] = v
	}
	images := make([]*models.TImage, 0)
	for _, v := range data {
		e1 := utils.SaveImage("\\big\\"+v.Name, v.Url)
		e2 := utils.SaveImage("\\small\\"+v.Name, v.SimpleUrl)
		if e1 != nil && e2 != nil {
			log.Println(e1.Error(), e1.Error())
			continue
		}
		if _, ok := currentImageM[v.Name]; ok {
			continue
		}
		images = append(images, &models.TImage{
			FilmId: filmId,
			Name:   v.Name,
		})
	}
	if e := new(models.TImage).BulkCreate(nil, images); e != nil {
		return e
	}
	return nil
}
