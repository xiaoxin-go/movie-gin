package follow

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
	actressIdList := make([]int, 0)
	db := model.DB.Model(&model.TUserFollow{}).Where("user_id = ?", user.Id).Select("actress_id").Find(&actressIdList)
	if db.Error != nil{
		zap.L().Info(fmt.Sprintf("查询user follow异常, %s", err.Error()))
		request.JSON(http.StatusOK, libs.ServerError("获取数据异常"))
		return
	}
	dataList := make([]model.TActress, 0)
	var total int64
	db = model.DB.Model(&model.TActress{}).Where("id in ?", actressIdList).Count(&total).
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