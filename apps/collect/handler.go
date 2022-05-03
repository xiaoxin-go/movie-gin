package collect

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"movie/libs"
	"movie/model"
	"movie/utils"
	"net/http"
)

type Handler struct {
	libs.Controller
}

var handler *Handler

func init() {
	handler = &Handler{}
}

func (h Handler) Get(request *gin.Context){
	page, pageSize := h.GetPagination(request)
	user, err := utils.GetCookieUser(request)
	if err != nil{
		zap.L().Info(fmt.Sprintf("获取用户信息异常, %s", err.Error()))
		return
	}
	filmIdList := make([]int, 0)
	db := model.DB.Model(&model.TUserCollect{}).Where("user_id = ?", user.Id).Select("film_id").Find(&filmIdList)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("查询user collect异常, %s", err.Error()))
		request.JSON(http.StatusOK, libs.ServerError("获取数据异常"))
		return
	}
	dataList := make([]model.TFilm, 0)
	var total int64
	db = model.DB.Model(&model.TFilm{}).Where("id in ?", filmIdList).Count(&total).
		Limit(pageSize).Offset((page - 1) * pageSize).Find(&dataList)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("查询film异常, %s", err.Error()))
		request.JSON(http.StatusOK, libs.ServerError("获取数据异常"))
		return
	}
	result := map[string]interface{}{
		"total": total,
		"data_list": dataList,
	}
	request.JSON(http.StatusOK, libs.Success(result, "ok"))
}