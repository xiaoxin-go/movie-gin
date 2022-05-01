package login

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"movie/libs"
	"movie/model"
	"movie/utils"
	"net/http"
)

func login(request *gin.Context){
	requestData := struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := request.BindJSON(&requestData)
	if err != nil{
		request.JSON(http.StatusOK, libs.ParamsError("解析参数异常"))
		return
	}
	fmt.Println("requestData----", requestData)
	user := model.TUser{Username: requestData.Username}
	db := model.DB.Where("username = ?", requestData.Username).First(&user)
	if errors.Is(db.Error, gorm.ErrRecordNotFound){
		request.JSON(http.StatusOK, libs.ParamsError("用户名或密码错误"))
		return
	}
	fmt.Printf("user: %+v\n", user)
	if db.Error != nil{
		request.JSON(http.StatusOK, libs.ParamsError("服务器异常"))
		return
	}
	if requestData.Password != user.Password{
		request.JSON(http.StatusOK, libs.ParamsError("用户名或密码错误"))
		return
	}
	u4 := uuid.NewV4()
	u4Str := u4.String()
	r := utils.NewRedisDefault()
	r.SetJson(u4Str, user)
	if r.Error != nil{
		fmt.Println("存入redis error: ", r.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError("服务器异常"))
	}
	request.SetCookie("movie_cookie", u4Str, 3600, "/", "localhost", false, true)
	request.JSON(http.StatusOK, libs.Success(nil, "登录成功"))
}