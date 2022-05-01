package libs

import (
	"fmt"
	"movie/model"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler interface{
	List(request *gin.Context)
	Get(request *gin.Context)
	Post(request *gin.Context)
	Put(request *gin.Context)
	Delete(request *gin.Context)
}

type Controller struct {
	Data         interface{}
	DataList interface{}		// 查询列表
	PreloadList []string 		// 外键查询使用
	SearchFields []string      // 关键字搜索, 参数为search, 对应的模糊查询列表
	FilterFields []string      // 添加的过滤条件
	GetFields    []string      // 查询时，返回的字段
	OrderFields  []string      // 排序字段
	UpdateSelect []interface{} // 更新时，只更新的字段
	UpdateOmit   []string      // 更新时，不更新的字段
	isParseData  bool          // 是否已经序列化数据
}

func (c *Controller) QuerySet(request *gin.Context)(db *gorm.DB, total int64){
	page, pageSize := c.GetPagination(request)
	db = model.DB.Model(c.Data)

	// 获取指定的字段
	if c.GetFields != nil {
		db = db.Select(c.GetFields)
	}

	// 若is_delete存在，则只获取is_delete为0的值
	if c.isDeleteExists() {
		db = db.Where("is_delete = ?", 0)
	}

	// 排序
	if c.OrderFields != nil {
		for _, field := range c.OrderFields {
			db = db.Order(field)
		}
	} else {
		db = db.Order("-id")
	}

	// 获取匹配查询字符串
	fields := map[string]string{}
	for _, field := range c.FilterFields {
		value := request.Query(field)
		if value == "" {
			continue
		}
		fields[field] = value
	}

	// 匹配查询
	//db = db.Where(fields)
	for field, value := range fields {
		db = db.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// 模糊查询
	// 最后拼接成的sql是 (username LIKE "%search%" or key LIKE "%search%")
	search := request.Query("search")
	if search != "" {
		orSql := ""
		searchArgs := make([]interface{}, 0)
		for _, field := range c.SearchFields {
			if orSql == ""{
				orSql += fmt.Sprintf("%s LIKE ?", field)
			}else{
				orSql += fmt.Sprintf(" or %s LIKE ?", field)
			}
			searchArgs = append(searchArgs, "%"+search+"%")
		}
		orSql += ""
		db = db.Where(orSql, searchArgs...)
	}

	db = db.Count(&total)
	if request.DefaultQuery("all", "false") == "false"{
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	// 关联查询
	for _, item := range c.PreloadList{
		db = db.Preload(item)
	}
	return
}
func (c *Controller) Response(request *gin.Context, results interface{}, total int64, err error){
	if err != nil {
		zap.L().Error("db.Limit error: " + err.Error())
		request.JSON(http.StatusOK, ServerError(err.Error()))
		return
	}

	data := map[string]interface{}{
		"data_list": results,
		"total":     total,
	}
	request.JSON(http.StatusOK, Success(data, "ok"))
}
func (c *Controller) List(request *gin.Context) {
	db, total := c.QuerySet(request)
	if c.DataList != nil {
		fmt.Println("=============================")
		db = db.Find(c.DataList)
		c.Response(request, c.DataList, total, db.Error)
	}else{
		results := make([]map[string]interface{}, 0)
		db = db.Find(&results)
		c.Response(request, results, total, db.Error)
	}
}
func (c *Controller) Get(request *gin.Context) {
	id := c.GetParamId(request)
	db := model.DB.Model(c.Data)
	data := map[string]interface{}{}
	result := db.Select(c.GetFields).First(&data, id)
	if result.Error != gorm.ErrRecordNotFound && result.Error != nil {
		zap.L().Error("db.First error: " + result.Error.Error())
		request.JSON(http.StatusOK, ServerError(result.Error.Error()))
		return
	}
	request.JSON(http.StatusOK, Success(data, "ok"))
}
func (c *Controller) Post(request *gin.Context) {
	zap.L().Info("<!--添加数据----------------------")
	defer zap.L().Info("添加数据end-----------------!>")
	fmt.Printf("%+v, %+v\n", c.Data, c.isParseData)
	defer func() { c.isParseData = false }()
	if !c.isParseData {
		c.ShouldBindJSON(request, c.Data)
	}


	db := model.DB.Create(c.Data)
	zap.L().Info(fmt.Sprintf("request_data: %+v", c.Data))
	if db.Error != nil {
		if strings.Contains(db.Error.Error(), "Duplicate entry"){
			request.JSON(http.StatusOK, ServerError("添加失败: 数据已存在"))
		}else{
			zap.L().Error("model.db.Create error: " + db.Error.Error())
			request.JSON(http.StatusOK, ServerError("插入数据异常: "+db.Error.Error()))
		}
		return
	}
	request.JSON(http.StatusOK, Success(c.Data, "ok"))
}
func (c *Controller) Put(request *gin.Context) {
	zap.L().Info("<!--更新数据----------------------")
	defer zap.L().Info("更新数据end-----------------!>")
	defer func() { c.isParseData = false }()
	if !c.isParseData {
		//c.Data = make(map[string]interface{})
		c.ShouldBindJSON(request, c.Data)
	}
	fmt.Printf("%+v\n", c.Data)
	// 设置操作人
	zap.L().Info(fmt.Sprintf("request_data: %+v", c.Data))
	id := c.GetParamId(request)
	zap.L().Info(fmt.Sprintf("update_id: %d", id))
	db := model.DB.Model(c.Data)
	db = db.Where("id = ?", id)

	if c.UpdateSelect != nil {
		db = db.Select(c.UpdateSelect[0], c.UpdateSelect[1:]...)
	}

	if c.UpdateOmit != nil {
		db = db.Omit(c.UpdateOmit...)
	}

	db = db.Updates(c.Data)
	if db.Error != nil {
		zap.L().Error("db.Updates error: " + db.Error.Error())
		request.JSON(http.StatusOK, ServerError(db.Error.Error()))
		return
	}
	request.JSON(http.StatusOK, Success(db.RowsAffected, "ok"))
}
func (c *Controller) Delete(request *gin.Context) {
	zap.L().Info("<!--删除数据----------------------")
	defer zap.L().Info("删除数据end-----------------!>")
	id := c.GetParamId(request)
	zap.L().Info(fmt.Sprintf("delete_id: %d", id))
	var db *gorm.DB
	// 若model存在is_delete，则设置is_delete为1
	if c.isDeleteExists() {
		db = model.DB.Model(c.Data).Where("id = ?", id).Update("is_delete", 1)
	} else {
		db = model.DB.Delete(c.Data, id)
	}

	if db.Error != nil {
		zap.L().Error("model.DB.Delete error: " + db.Error.Error())
		request.JSON(http.StatusOK, ServerError(db.Error.Error()))
		return
	}
	request.JSON(http.StatusOK, Success(db.RowsAffected, "ok"))
}

// 判断是否存在is_delete
func (c *Controller) isDeleteExists() bool {
	if c.Data == nil {
		return false
	}
	te := reflect.ValueOf(c.Data)
	te = te.Elem()
	fe := te.FieldByName("is_delete")
	return fe.IsValid()
}

// 初始化ID
func (c *Controller) initId() {
	if c.Data == nil {
		return
	}
	te := reflect.ValueOf(c.Data)
	te = te.Elem()
	fe := te.FieldByName("Id")
	if fe.IsValid() {
		fe.SetInt(0)
	}
}


// ShouldBindJSON 获取前端传来的数据
func (c *Controller) ShouldBindJSON(request *gin.Context, data interface{}) {
	c.isParseData = true
	c.initId()
	err := request.ShouldBindJSON(data)
	if err != nil {
		zap.L().Error("request.ShouldBindJson error: " + err.Error())
		request.JSON(http.StatusOK, ParamsError("解析参数异常: "+err.Error()))
		c.stopRun()
	}
}

// GetParamId 获取参数中传的ID
func (c *Controller) GetParamId(request *gin.Context) int {
	idStr := request.Param("id")
	idStr = strings.Trim(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		request.JSON(http.StatusOK, ParamsError(err.Error()))
		c.stopRun()
	}
	return id
}

// GetUserInfo 获取用户信息
func (c *Controller) GetUserInfo(request *gin.Context) ( err error) {
	// 获取cookie
	sessionId, _ := request.Cookie("sso_session_id")
	if sessionId == "" {
	} else {
	}
	return
}

// GetPagination 获取分页内容
func (c *Controller) GetPagination(request *gin.Context) (page, pageSize int) {
	pageStr := request.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	pageSizeStr := request.Query("page_size")
	if pageSizeStr == "" {
		pageSizeStr = "20"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}
	return
}

func (c *Controller) stopRun() {
	panic("stop run")
}
