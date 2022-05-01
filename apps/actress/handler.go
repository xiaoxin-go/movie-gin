package actress

import (
	"fmt"
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
		request.JSON(http.StatusOK, libs.ParamsError("已在请求中...."))
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
		utils.InsertFilmData(filmData, false)
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}
func (c *Handler) Follow(request *gin.Context){
	id := c.GetParamId(request)
	user, err := utils.GetCookieUser(request)
	if err != nil{
		zap.L().Info(fmt.Sprintf("获取用户信息异常, %s", err.Error()))
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}
	count := isFollow(id, user.Id)
	if count > 0{
		request.JSON(http.StatusOK, libs.Success(nil, "收藏成功"))
		return
	}
	userFollow := model.TUserFollow{ActressId: id, UserId: user.Id}
	db := model.DB.Create(&userFollow)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("创建user follow异常, %s, data: %+v", db.Error.Error(), userFollow))
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
		return
	}
	request.JSON(http.StatusOK, libs.Success(nil, "关注成功"))
	return
}
func (c *Handler) UnFollow(request *gin.Context){
	id := c.GetParamId(request)
	user, err := utils.GetCookieUser(request)
	if err != nil{
		zap.L().Error(fmt.Sprintf("获取用户信息异常, %s", err.Error()))
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}
	userFollow := model.TUserFollow{}
	db := model.DB.Where("actress_id = ? and user_id = ?", id, user.Id).First(&userFollow)
	if db.Error != nil{
		zap.L().Error(fmt.Sprintf("查询user follow error: %s, actress_id: %d, user_id: %d", db.Error.Error(), id, user.Id))
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
		return
	}
	db = model.DB.Delete(&userFollow)
	if db.Error != nil{
		zap.L().Error(fmt.Sprintf("删除user follow error: %s, userFollow: %+v", db.Error.Error(), userFollow))
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
		return
	}
	request.JSON(http.StatusOK, libs.Success(nil, "ok"))
}
func (c *Handler) IsFollow(request *gin.Context){
	id := c.GetParamId(request)
	user, err := utils.GetCookieUser(request)
	if err != nil{
		zap.L().Info(fmt.Sprintf("获取用户信息异常, %s", err.Error()))
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}
	count := isFollow(id, user.Id)
	request.JSON(http.StatusOK, libs.Success(count, "ok"))
	return
}

func isFollow(actressId, userId int)(result int64){
	var count int64
	db := model.DB.Model(&model.TUserFollow{}).Where("actress_id = ? and user_id = ?", actressId, userId).Count(&count)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("查询user follow异常, %s, actress_id: %d, user_id: %d", db.Error.Error(), actressId, userId))
		return
	}
	return count
}