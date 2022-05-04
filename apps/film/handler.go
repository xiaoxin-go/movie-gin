package film

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	config "movie/conf"
	"movie/libs"
	"movie/model"
	"movie/utils"
	"net/http"
	"os"
	"strings"
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
	utils.InsertFilmData(film.Data(), true)
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
func (c *Handler) AddLink(request *gin.Context){
	id := c.GetParamId(request)
	data := model.TFilm{}
	db := model.DB.First(&data, id)
	if db.Error != nil{
		zap.L().Error("获取film异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
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
	utils.InsertFilmLinks(data.Id, film.Data().Links)
	request.JSON(http.StatusOK, libs.Success(data, "ok"))
}
func (c *Handler) AddImage(request *gin.Context){
	id := c.GetParamId(request)
	data := model.TFilm{}
	db := model.DB.First(&data, id)
	if db.Error != nil{
		zap.L().Error("获取film异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
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
	utils.InsertFilmImages(data.Id, film.Data().Images)
	request.JSON(http.StatusOK, libs.Success(data, "ok"))
}
func (c *Handler) SaveImage(request *gin.Context){
	id := c.GetParamId(request)
	images := make([]model.TImage, 0)
	db := model.DB.Where("film_id = ?", id).Find(&images)
	if db.Error != nil{
		zap.L().Error("获取film异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	for _, image := range images{
		utils.SaveImage(image.Name + "-simple", image.SimpleUrl)
		utils.SaveImage(image.Name, image.Url)
	}
}
func (c *Handler) Collect(request *gin.Context){
	id := c.GetParamId(request)
	user, err := utils.GetCookieUser(request)
	if err != nil{
		zap.L().Info(fmt.Sprintf("获取用户信息异常, %s", err.Error()))
		return
	}
	count := isCollect(id, user.Id)
	if count > 0{
		request.JSON(http.StatusOK, libs.Success(nil, "收藏成功"))
		return
	}
	userCollect := model.TUserCollect{FilmId: id, UserId: user.Id}
	db := model.DB.Create(&userCollect)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("创建user collect异常, %s, data: %+v", db.Error.Error(), userCollect))
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
		return
	}
	request.JSON(http.StatusOK, libs.Success(nil, "收藏成功"))
	return
}
func (c *Handler) UnCollect(request *gin.Context){
	id := c.GetParamId(request)
	user, err := utils.GetCookieUser(request)
	if err != nil{
		zap.L().Error(fmt.Sprintf("获取用户信息异常, %s", err.Error()))
		return
	}
	userCollect := model.TUserCollect{}
	db := model.DB.Where("film_id = ? and user_id = ?", id, user.Id).First(&userCollect)
	if db.Error != nil{
		zap.L().Error(fmt.Sprintf("查询user collect error: %s, film_id: %d, user_id: %d", db.Error.Error(), id, user.Id))
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
		return
	}
	db = model.DB.Delete(&userCollect)
	if db.Error != nil{
		zap.L().Error(fmt.Sprintf("删除user collect error: %s, userCollect: %+v", db.Error.Error(), userCollect))
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
		return
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}
func (c *Handler) IsCollect(request *gin.Context){
	var count int64
	id := c.GetParamId(request)
	user, err := utils.GetCookieUser(request)
	if err == nil{
		count = isCollect(id, user.Id)
	}
	request.JSON(http.StatusOK, libs.Success(count, "ok"))
	return
}
func (c *Handler) IsPlayer(request *gin.Context){
	id := c.GetParamId(request)
	data := model.TFilm{}
	db := model.DB.First(&data, id)
	if db.Error != nil{
		zap.L().Error("获取film异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	filename := fmt.Sprintf("%s%s.mp4", config.Config.MoviePath, data.Name)
	_, err := os.Stat(filename)
	if err == nil{
		request.JSON(http.StatusOK, libs.Success(1, "ok"))
		return
	}
	request.JSON(http.StatusOK, libs.Success(0, "ok"))
}

func isCollect(filmId, userId int)(result int64){
	var count int64
	db := model.DB.Model(&model.TUserCollect{}).Where("film_id = ? and user_id = ?", filmId, userId).Count(&count)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("查询user collect异常, %s, film_id: %d, user_id: %d", db.Error.Error(), filmId, userId))
		return
	}
	return count
}

func (c *Handler)Player(request *gin.Context){
	dir, err := os.ReadDir(config.Config.MoviePath)
	page, pageSize := c.GetPagination(request)
	if err != nil{
		request.JSON(http.StatusOK, libs.ServerError("读取文件失败"))
		return
	}
	nameList := make([]string, 0)
	for _, item := range dir{
		nameList = append(nameList, strings.Split(item.Name(), ".")[0])
	}
	var total int64
	dataList := make([]model.TFilm, 0)
	db := model.DB.Model(&model.TFilm{}).Where("name in ?", nameList).Count(&total).
			Limit(pageSize).Offset((page - 1) * pageSize).Find(&dataList)
	if db.Error != nil{
		zap.L().Error(fmt.Sprintf("获取film信息异常: %s", db.Error.Error()))
		request.JSON(http.StatusOK, libs.ServerError("获取电影信息失败"))
		return
	}
	result := map[string]interface{}{
		"total": total,
		"data_list": dataList,
	}
	request.JSON(http.StatusOK, libs.Success(result, "ok"))
}