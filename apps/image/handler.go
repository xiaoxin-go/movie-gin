package image

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
	handler.OrderFields = []string{"id"}
	handler.FilterFields = []string{"film_id"}
	handler.Data = &model.TImage{}
}
