package genre

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"movie/libs"
	"movie/model"
	"net/http"
)

type Handler struct {
	libs.Controller
}

var handler *Handler

func init() {
	handler = &Handler{}
	handler.Data = &model.TGenre{}
}

func (c *Handler) Film(request *gin.Context) {
	id := c.GetParamId(request)
	page, pageSize := c.GetPagination(request)
	genreFilms := make([]model.TGenreFilm, 0)

	var total int64

	db := model.DB.Model(&model.TGenreFilm{}).Where("genre_id", id).Count(&total)
	if db.Error != nil{
		zap.L().Error("db.First error: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&genreFilms)
	if db.Error != nil{
		zap.L().Error("db.First error: " + db.Error.Error())
		request.JSON(http.StatusOK, libs.ServerError(db.Error.Error()))
		return
	}
	filmIds := make([]int, 0)
	for _, item := range genreFilms{
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