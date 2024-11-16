package services

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"io"
	"log"
	"movie/database"
	"movie/models"
	"movie/utils"
	"os"
)

func GetFilmDetail(name string) (*models.TFilm, error) {
	result := models.TFilm{}
	if e := result.GetDetailByName(name); e != nil {
		return nil, e
	}
	return &result, nil
}

func CreateFilm(name string) error {
	// 1. 获取电影详情
	service, err := utils.NewService()
	if err != nil {
		return err
	}
	defer service.Stop()
	wd, e := utils.NewWindow()
	if e != nil {
		return e
	}
	log.Println("获取电影信息----->")
	film := utils.NewFilm(name, wd)
	if e := film.Error(); e != nil {
		return e
	}
	log.Println("保存电影信息--->")
	f, e := saveFilm(film.Data())
	if e != nil {
		return e
	}

	log.Println("保存链接--->")
	if e := saveLinks(f.Id, film.Data().Links); e != nil {
		log.Fatalln(e.Error())
	}
	log.Println("保存logo")
	if e := utils.SaveImage("/logo/"+film.Name, film.Data().ImageUrl); e != nil {
		log.Println("保存logo失败:", e.Error())
	}
	log.Println("剪切logo")
	cutImage("/logo/" + film.Name)

	log.Println("保存图片--->")
	if e := saveImages(f.Id, film.Data().Images); e != nil {
		log.Println(e.Error())
	}
	log.Println("保存演员信息--->")
	log.Println(film.Data().Actresses)
	actress := saveActress(f, film.Data().Actresses, wd)
	if e1 := database.DB.Model(f).Association("Actresses").Replace(actress); e1 != nil {
		log.Println("创建关联失败: ", e1.Error())
	}
	return nil
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
		fmt.Println(actress)
		if e := database.DB.Where("name = ?", d.Name).First(&actress).Error; errors.Is(e, gorm.ErrRecordNotFound) {
			// 如果没有，则把这个电影做为封面
			if e := copyFile("./build/images/logo/"+film.Name+"-right.jpg", "./build/images/actress/"+actress.Name+".jpg"); e != nil {
				log.Println("设置演员头像失败--->", e.Error())
			}
			if e := database.DB.Create(&actress).Error; e != nil {
				log.Println("创建获取演员信息失败: ", c.Error())
				continue
			}
		}
		result = append(result, &actress)
	}
	return result
}

func copyFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// 复制源文件内容到目标文件
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// 强制写入到磁盘
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func cutImage(name string) {
	filename := "./build/images" + name + ".jpg"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	bounds := img.Bounds()
	// 定义要截取的区域 (x, y, width, height)
	right := 442
	if bounds.Dx() == 800 {
		right = 421
	}
	rect := image.Rect(right, 0, bounds.Dx(), bounds.Dy()) // 裁剪 50,50 到 200,200 区域
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)
	// 创建输出文件
	outFilename := "./build/images" + name + "-right.jpg"
	outFile, err := os.Create(outFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()
	// 将裁剪后的图片保存为 JPEG 格式
	err = jpeg.Encode(outFile, croppedImg, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
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
		e1 := utils.SaveImage("/big/"+v.Name, v.Url)
		if e1 != nil {
			log.Println(e1.Error(), e1.Error())
			continue
		}
		if _, ok := currentImageM[v.Name]; ok {
			continue
		}
		zipImage(v.Name)
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
func zipImage(name string) {
	// 打开原始图片文件
	filename := "./build/images/big/" + name + ".jpg"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("无法解码图片:", err)
		return
	}
	outFilename := "./build/images/small/" + name + ".jpg"
	// 创建输出文件
	outFile, err := os.Create(outFilename)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer outFile.Close()

	// 设置压缩质量，值越低压缩越厉害，质量也越差
	options := jpeg.Options{
		Quality: 30, // 0-100，50为中等压缩率
	}

	// 保存压缩后的图片
	err = jpeg.Encode(outFile, img, &options)
	if err != nil {
		fmt.Println("无法保存图片:", err)
		return
	}

	fmt.Println("图片已成功压缩并保存为small", name)
}
