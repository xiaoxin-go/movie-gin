package film

import (
	"movie/libs"
	"movie/model"
)

type Handler struct {
	libs.Controller
}

var handler *Handler

func init() {
	handler = &Handler{}
	handler.OrderFields = []string{"-release_date"}
	handler.SearchFields = []string{"name"}
	handler.Data = &model.TFilm{}
}