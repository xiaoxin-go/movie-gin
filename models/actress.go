package models

type TActress struct {
	BaseModel
	Name     string   `gorm:"size:50;unique" json:"name"`
	Height   string   `gorm:"size:50" json:"height"`
	Cup      string   `gorm:"size:50" json:"cup"`
	Birthday string   `gorm:"size:50" json:"birthday"`
	Films    []*TFilm `gorm:"many2many:t_actress_film;" json:"films"`
}
