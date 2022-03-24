package actress

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"movie/libs"
	"movie/model"
	"movie/utils"
	"net/http"
	"sync"
)

type Handler struct {
	libs.Controller
}

var handler *Handler
var lock *sync.Mutex
var activeMap map[int]bool

func init() {
	lock = &sync.Mutex{}
	activeMap = make(map[int]bool)
	handler = &Handler{}
	handler.Data = &model.TActress{}
}

func (c *Handler) Film(request *gin.Context) {
	id := c.GetParamId(request)
	page, pageSize := c.GetPagination(request)
	actressFilms := make([]model.TActressFilm, 0)

	var total int64

	db := model.DB.Model(&model.TActressFilm{}).Where("actress_id", id).Count(&total)
	if db.Error != nil{
		zap.L().Error("db.First error: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&actressFilms)
	if db.Error != nil{
		zap.L().Error("db.First error: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	filmIds := make([]int, 0)
	for _, item := range actressFilms{
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

func (c *Handler) Delete(request *gin.Context){
	id := c.GetParamId(request)
	db := model.DB.Where("actress_id = ?", id).Delete(&model.TActressFilm{})
	if db.Error != nil{
		zap.L().Error("获取数据异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	db = model.DB.Delete(&model.TActress{}, id)
	if db.Error != nil{
		zap.L().Error("获取数据异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}

func (c *Handler) Films(request *gin.Context){
	id := c.GetParamId(request)
	if activeMap[id]{
		return
	}
	activeMap[id] = true
	defer delete(activeMap, id)
	lock.Lock()
	defer lock.Unlock()
	actress := model.TActress{}
	db := model.DB.First(&actress, id)
	if db.Error != nil{
		zap.L().Error("select actress error => " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	service, err := utils.NewService()
	if err != nil {
		zap.L().Error("new service error => " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	defer service.Stop()

	wd, err := utils.NewWindow()
	if err != nil {
		zap.L().Error("new window error => " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	defer wd.Close()
	actressController := utils.NewActressController(actress.Url, wd)
	data := actressController.Data()
	zap.L().Info("")
	if actressController.Error() != nil{
		zap.L().Error("get actress info error => " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	for _, item := range data.Films{
		zap.L().Info("get film ===> " + item)
		var count int64
		model.DB.Model(&model.TFilm{}).Where("name = ?", item).Count(&count)
		if count > 0{
			continue
		}
		filmController := utils.NewFilm(item, wd)
		filmData := filmController.Data()
		if filmController.Error() != nil{
			zap.L().Error("get filmData error => " + filmController.Error().Error())
			continue
		}
		zap.L().Info("insert film ===> " + item)
		utils.InsertFilmData(filmData)
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}