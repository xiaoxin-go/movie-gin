package root 

import (
	"fip/apps/system/menu"
	"fip/apps/system/role"
	user2 "fip/apps/system/user"
	"fip/model"
	"fmt"
	"fip/libs"
	"fip/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type Handler struct{
	libs.Controller
}

var handler *Handler

func init(){
	handler = &Handler{}
}

func (h *Handler)Index(request *gin.Context){
	request.HTML(http.StatusOK, "index.html", gin.H{"title": "fip平台", "ce": "123456"})
}

// 获取菜单
func (h *Handler) GetMenu(request *gin.Context){
	fmt.Println("---------------------------!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// 1. 获取用户信息
	userInfo, _ := h.GetUserInfo(request)
	fmt.Println("---------------------------------")
	user := user2.User{Username: userInfo.Username}
	db := model.DB.First(&user)
	if db.Error != nil{
		zap.L().Error("获取用户信息，查询user表异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError("获取用户信息异常: " + db.Error.Error()))
		return
	}

	// 2. 根据用户角色ID，获取关联的菜单ID列表
	roleMenuList := make([]role.TRoleMenu, 0)
	db = model.DB.Where("role_id = ?", user.RoleId).Find(&roleMenuList)
	if db.Error != nil{
		zap.L().Error("获取菜单，查询role_menu表异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError("获取菜单信息异常: " + db.Error.Error()))
		return
	}
	// 3. 根据菜单ID,获取菜单信息
	menuIdList := make([]int, 0)
	for _, item := range roleMenuList{
		menuIdList = append(menuIdList, item.MenuId)
	}
	menuList := make([]menu.Menu, 0)
	db = model.DB.Order("sort").Find(&menuList, menuIdList)
	if db.Error != nil{
		zap.L().Error("获取菜单，查询menu表异常: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError("获取菜单信息异常: " + db.Error.Error()))
		return
	}
	// 4. 循环菜单列表，区分parent和child，并将child挂到parent下
	parentMenu := make([]menu.Menu, 0)
	childMap := make(map[int][]menu.Menu)
	for _, item := range menuList{
		if item.Type == "parent"{
			parentMenu = append(parentMenu, item)
			continue
		}
		// 如果是子菜单，将子菜单信息挂到父菜单ID下，后面循环父菜单信息使用
		if _, ok := childMap[item.ParentId]; ok{
			childMap[item.ParentId] = append(childMap[item.ParentId], item)
		}else{
			childMap[item.ParentId] = []menu.Menu{item}
		}
	}

	// 5. 循环父菜单信息，将子菜单信息从childMap取出放到child里面
	results := make([]menu.Menu, 0)
	for _, item := range parentMenu{
		item.Child = childMap[item.Id]
		results = append(results, item)
	}

	// 6. 返回菜单信息
	data := map[string]interface{}{
		"data_list": results,
	}
	request.JSON(http.StatusOK, libs.Success(data, "ok"))
}

// 获取用户信息
func (h *Handler) GetUser(request *gin.Context){
	userInfo, err := h.GetUserInfo(request)
	if err != nil{
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}
	request.JSON(http.StatusOK, libs.Success(userInfo, "ok"))
}

// 获取加密公钥，保存私钥
func(h *Handler) GetPublicKey(request *gin.Context){
	// 获取私钥需要传入一个申请人
	author := request.Query("author")
	if author == ""{
		request.JSON(http.StatusOK, libs.ParamsError("the author cannot be empty"))
		return
	}
	// 生成公私钥
	privateKey, publicKey, err := utils.GenerateKey()
	if err != nil{
		zap.L().Error("GenerateKey error: " + err.Error())
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}

	// 保存公私钥
	secretKey := fmt.Sprintf("%s_secret_key", author)
	r := utils.NewRedisDefault()
	r.SetString(secretKey, privateKey)
	r.SetExpire(secretKey, 60)
	if r.Error != nil{
		zap.L().Error("redis.SetString error: " + r.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(r.Error.Error()))
		return
	}

	result := map[string]string{
		"public_key": publicKey,
	}

	request.JSON(http.StatusOK, libs.Success(result, "ok"))
}

// 创建新项目
func(h *Handler) CreateApp(request *gin.Context){
	// 创建app
	// 1. 获取app name，创建文件夹
	appName := request.Query("app_name")
	dirPath := fmt.Sprintf("apps/%s", appName)
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil{
		zap.L().Error("os.Mkdir error" + err.Error())
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}
	// 创建handler文件
	f, err := os.Create(fmt.Sprintf("%s/handle.go", dirPath))
	if err != nil{
		zap.L().Error("os.Create error" + err.Error())
		request.JSON(http.StatusOK, libs.ServerError(err.Error()))
		return
	}
	_, _ = f.WriteString("package " + appName)
}