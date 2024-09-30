package models

type TGenre struct {
	BaseModel
	Name string `gorm:"size:50;unique" json:"name"`
}
