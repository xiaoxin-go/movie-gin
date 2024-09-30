package libs

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"movie/database"
	"net/http"
	"strconv"
	"strings"
)

type LogContent struct {
	Action string
	Before any
	After  any
}

var AddLog func(ctx *gin.Context, content LogContent)

type Restfuller interface {
	Get(c *gin.Context)
	List(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Instance interface {
	GetId() int
}

type Controller struct {
	ListFunc    func() any
	ModelFunc   func() Instance
	OrderFilter func(db *gorm.DB) *gorm.DB
	QueryFilter func(db *gorm.DB) *gorm.DB
	db          *gorm.DB
}

func AdvancedQuery(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range ctx.Request.URL.Query() {
			// 排除page size
			if k == "page" || k == "size" {
				continue
			}
			if !strings.Contains(k, "__") {
				db = db.Where(fmt.Sprintf("%s = ?", k), v[0])
				continue
			}
			// name__contains  取出name和contains
			ks := strings.Split(k, "__")
			k1 := ks[0]
			kType := ks[1]

			// ?name=lisi&age__gt=18&id__in=1,2,3
			switch kType {
			case "eq":
				db = db.Where(fmt.Sprintf("%s = ?", k1), v[0])
			case "neq":
				db = db.Where(fmt.Sprintf("%s != ?", k1), v[0])
			case "contains":
				db = db.Where(fmt.Sprintf("%s like ?", k1), "%"+v[0]+"%")
			case "in":
				db = db.Where(fmt.Sprintf("%s in ?", k1), strings.Split(v[0], ","))
			case "not_in":
				db = db.Where(fmt.Sprintf("%s not in ?", k1), strings.Split(v[0], ","))
			case "gte":
				db = db.Where(fmt.Sprintf("%s >= ?", k1), v[0])
			case "lte":
				db = db.Where(fmt.Sprintf("%s <= ?", k1), v[0])
			case "gt":
				db = db.Where(fmt.Sprintf("%s > ?", k1), v[0])
			case "lt":
				db = db.Where(fmt.Sprintf("%s < ?", k1), v[0])
			}
		}
		return db
	}
}

func Pagination(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 20
		}
		db.Offset((page - 1) * size).Limit(size)
		return db
	}
}

func (c *Controller) orderBy(db *gorm.DB) *gorm.DB {
	if c.OrderFilter != nil {
		return c.OrderFilter(db)
	}
	return db.Order("-id")
}

func (c *Controller) QuerySet(ctx *gin.Context) (db *gorm.DB) {
	m := c.ModelFunc()
	db = database.DB.Model(m).Scopes(AdvancedQuery(ctx))
	if c.QueryFilter != nil {
		db.Scopes(c.QueryFilter)
	}
	return
}
func (c *Controller) Count(db *gorm.DB, count *int64) *gorm.DB {
	return db.Count(count)
}
func (c *Controller) OrderBy(db *gorm.DB) *gorm.DB {
	return c.orderBy(db)
}
func (c *Controller) Pagination(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, size := c.GetPagination(ctx)
	return Pagination(page, size)
}
func (c *Controller) Response(ctx *gin.Context, results interface{}, total int64, err error) {
	if err != nil {
		HttpServerError(ctx, err.Error())
		return
	}
	data := map[string]interface{}{
		"data":  results,
		"total": total,
	}
	ctx.JSON(http.StatusOK, Success(data, "ok"))
}
func (c *Controller) QueryListData(ctx *gin.Context) (results any, total int64, err error) {
	db := c.QuerySet(ctx).Count(&total).Scopes(c.orderBy, c.Pagination(ctx))
	results = c.ListFunc()
	if e := db.Find(results).Error; e != nil {
		err = fmt.Errorf("查询数据异常, err: %w", e)
		return
	}
	return
}
func (c *Controller) checkModelFunc() error {
	if c.ModelFunc == nil {
		return errors.New("该接口未定义instance，请联系管理员")
	}
	return nil
}
func (c *Controller) List(ctx *gin.Context) {
	results, total, err := c.QueryListData(ctx)
	if err != nil {
		HttpServerError(ctx, err.Error())
		return
	}
	HttpListSuccess(ctx, results, total)
}
func (c *Controller) Get(ctx *gin.Context) {
	if e := c.checkModelFunc(); e != nil {
		HttpServerError(ctx, e.Error())
		return
	}
	id, e := c.GetId(ctx)
	if e != nil {
		HttpParamsError(ctx, e.Error())
		return
	}
	data := c.ModelFunc()
	if err := database.DB.Where("id = ?", id).First(data).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		HttpParamsError(ctx, fmt.Sprintf("根据ID<%d>未获取到数据", id))
		return
	} else if err != nil {
		HttpServerError(ctx, fmt.Sprintf("获取数据异常, id: %d, err: %s", id, err.Error()))
		return
	}
	HttpSuccess(ctx, data, "ok")
}
func (c *Controller) Create(ctx *gin.Context) {
	if c.ModelFunc == nil {
		HttpServerError(ctx, "该接口未定义instance，请联系管理员")
	}
	data := c.ModelFunc()
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		HttpParamsError(ctx, fmt.Sprintf("参数解析异常: <%s>", err.Error()))
		return
	}
	if err := database.DB.Save(data).Error; err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		HttpParamsError(ctx, "添加失败：数据已存在")
		return
	} else if err != nil {
		HttpServerError(ctx, fmt.Sprintf("添加数据异常: <%s>", err.Error()))
		return
	}
	if AddLog != nil {
		content := LogContent{
			Action: "insert",
			After:  data,
		}
		AddLog(ctx, content)
	}

	HttpSuccess(ctx, data, "ok")
}
func (c *Controller) Update(ctx *gin.Context) {
	data := c.ModelFunc()
	if e := ctx.ShouldBindJSON(data); e != nil {
		HttpParamsError(ctx, "参数解析失败, err: %s", e.Error())
		return
	}
	id := data.GetId()
	oldM := c.ModelFunc()
	if e := database.DB.First(oldM, id).Error; e != nil {
		HttpServerError(ctx, "获取数据失败, err: %s", e.Error())
		return
	}
	if e := database.DB.Model(data).Where("id = ?", id).Updates(data).Error; e != nil {
		HttpServerError(ctx, "更新数据异常, err: %s", e.Error())
		return
	}

	if AddLog != nil {
		content := LogContent{
			Action: "update",
			Before: oldM,
			After:  data,
		}
		AddLog(ctx, content)
	}
	HttpSuccess(ctx, data, "ok")
}
func (c *Controller) Delete(ctx *gin.Context) {
	id, e := c.GetId(ctx)
	if e != nil {
		HttpParamsError(ctx, e.Error())
		return
	}
	m := c.ModelFunc()
	if err := database.DB.First(m, id).Error; err != nil {
		HttpServerError(ctx, "获取数据异常, err: %s", err.Error())
		return
	}
	if err := database.DB.Delete(m).Error; err != nil {
		HttpServerError(ctx, "删除数据异常, err: %s", err.Error())
		return
	}
	if AddLog != nil {
		content := LogContent{
			Action: "delete",
			Before: m,
		}
		AddLog(ctx, content)
	}
	HttpSuccess(ctx, m, "ok")
}
func (c *Controller) GetId(ctx *gin.Context) (int, error) {
	idStr := ctx.Query("id")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		return 0, errors.New("id不能是字符串")
	}
	if id <= 0 {
		return 0, errors.New("ID必须大于0")
	}
	return id, nil
}
func (c *Controller) BatchDelete(ctx *gin.Context) {
	params := struct {
		DataIds []int `json:"data_ids"`
	}{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		HttpParamsError(ctx, fmt.Sprintf("读取请求参数异常: <%s>", err.Error()))
		return
	}
	m := c.ModelFunc()

	if err := database.DB.Where("id in ?", params.DataIds).Delete(m).Error; err != nil {
		HttpServerError(ctx, fmt.Sprintf("批量删除数据异常: <%s>", err.Error()))
		return
	}
	if AddLog != nil {
		content := LogContent{
			Action: "batch_delete",
			Before: params.DataIds,
		}
		AddLog(ctx, content)
	}
	HttpSuccess(ctx, nil, "删除成功")
}

// GetPagination 获取分页内容
func (c *Controller) GetPagination(ctx *gin.Context) (page, pageSize int) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("size", "20")
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
func (c *Controller) StopRun() {
	panic("stop run")
}
