package link

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
	handler.FilterFields = []string{"film_id"}
	handler.Data = &model.TLink{}
}
