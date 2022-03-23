package film

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"movie/libs"
	"movie/model"
	"movie/utils"
	"net/http"
	"os"
)

type Handler struct {
	libs.Controller
}

var handler *Handler

func init() {
	handler = &Handler{}
	handler.OrderFields = []string{"-release_date"}
	handler.SearchFields = []string{"name", "title"}
	handler.Data = &model.TFilm{}
}

func (c *Handler) Get(request *gin.Context) {
	id := c.GetParamId(request)
	db := model.DB.Model(&model.TFilm{})
	data := model.TFilm{}
	result := db.Select(c.GetFields).Preload("Genres").Preload("Actresses").First(&data, id)
	fmt.Printf("%+v", data)

	if result.Error != nil{
		zap.L().Error("db.First error: " + result.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(result.Error.Error()))
		return
	}
	request.JSON(http.StatusOK, libs.Success(data, "ok"))
}
func (c *Handler) Delete(request *gin.Context){
	id := c.GetParamId(request)
	db := model.DB.Delete(&model.TFilm{}, id)
	if db.Error != nil{
		zap.L().Error("删除数据异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	model.DB.Where("film_id = ?", id).Delete(&model.TActressFilm{})
	model.DB.Where("film_id = ?", id).Delete(&model.TGenreFilm{})
	model.DB.Where("film_id = ?", id).Delete(&model.TLink{})
	images := make([]model.TImage, 0)
	model.DB.Where("film_id = ?", id).Find(&images).Delete(&model.TImage{})
	for _, image := range images{
		os.Remove("E:/FFOutput/static/images/" + image.Name + ".jpg")
		os.Remove("E:/FFOutput/static/images/" + image.Name + "-simple.jpg")
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}
func (c *Handler) Like(request *gin.Context) {
	id := c.GetParamId(request)

	genreFilms := make([]model.TGenreFilm, 0)

	var total int64

	db := model.DB.Where("film_id = ?", id).Order("-id").Find(&genreFilms)
	if db.Error != nil{
		zap.L().Error("获取数据异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	genreIds := make([]int, 0)
	for _, item := range genreFilms{
		genreIds = append(genreIds, item.GenreId)
	}
	fmt.Println(genreIds)
	db = model.DB.Where("genre_id in ?", genreIds[:2]).Order("-id").Group("film_id").Limit(6).Offset(0).Find(&genreFilms)
	filmIds := make([]int, 0)
	for _, item := range genreFilms{
		filmIds = append(filmIds, item.FilmId)
	}

	results := make([]model.TFilm, 0)
	db = model.DB.Where("id in ?", filmIds).Find(&results)
	if db.Error != nil{
		zap.L().Error("db.First error: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	data := map[string]interface{}{
		"data_list": results,
		"total":     total,
	}
	request.JSON(http.StatusOK, libs.Success(data, "ok"))
}
func (c *Handler) Post(request *gin.Context){
	data := struct{
		Name string `json:"name"`
	}{}
	c.ShouldBindJSON(request, &data)
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
	film := utils.NewFilm(data.Name, wd)
	if film.Error() != nil{
		request.JSON(http.StatusOK, libs.ServerError(film.Error().Error()))
		return
	}
	utils.InsertFilmData(film.Data())
	request.JSON(http.StatusOK, libs.Success(data, "ok"))
}
func (c *Handler) Cover(request *gin.Context){
	id := c.GetParamId(request)
	film := model.TFilm{}
	db := model.DB.First(&film, id)
	if db.Error != nil{
		zap.L().Error("获取film异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	actressFilms := make([]model.TActressFilm, 0)
	db = model.DB.Where("film_id = ?", id).Find(&actressFilms)
	if db.Error != nil{
		zap.L().Error("获取actress film异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	for _, item := range actressFilms{
		db = model.DB.Model(&model.TActress{}).Where("id = ?", item.ActressId).Update("image", film.Name)
		if db.Error != nil{
			zap.L().Error("获取actress film异常: " + db.Error.Error())
			request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
			return
		}
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}