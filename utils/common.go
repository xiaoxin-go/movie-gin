package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"movie/model"
	"time"
)

func StrToTime(date string) time.Time {
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date, local)
	return t
}

func StrToDate(date string)time.Time{
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02", date, local)
	return t
}

func GetCookieUser(request *gin.Context)(result model.TUser, err error){
	cookie, err := request.Cookie("movie_cookie")
	if err != nil{
		err = errors.New("用户未登录")
		return
	}
	result = model.TUser{}
	r := NewRedisDefault()
	r.GetJson(cookie, &result)
	if r.Error != nil{
		zap.L().Info(fmt.Sprintf("获取redis cookie异常, %s, cookie:%s", r.Error.Error(), cookie))
		err = errors.New("读取redis信息失败")
		return
	}
	return
}