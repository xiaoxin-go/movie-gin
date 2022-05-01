package register

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"movie/libs"
	"movie/model"
	"net/http"
)

func register(request *gin.Context){
	requestData := struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := request.BindJSON(&requestData)
	if err != nil{
		request.JSON(http.StatusOK, libs.ParamsError("解析参数异常"))
		return
	}
	user := model.TUser{Username: requestData.Username, Password: requestData.Password}
	db := model.DB.Where("username = ?", requestData.Username).First(&user)
	fmt.Println("----------------------", db.Error.Error())
	if !errors.Is(db.Error, gorm.ErrRecordNotFound){
		request.JSON(http.StatusOK, libs.ParamsError("服务器异常"))
		return
	}
	fmt.Printf("%+v\n", user)
	db = model.DB.Create(&user)
	fmt.Println("db.error---------", db.Error)
	if db.Error != nil{
		request.JSON(http.StatusOK, libs.ParamsError("服务器异常"))
		return
	}
	request.JSON(http.StatusOK, libs.Success(nil, "注册成功"))
}