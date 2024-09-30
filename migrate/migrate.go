package migrate

import (
	"fmt"
	"movie/database"
	"movie/models"
)

func AutoMigrate() {
	if e := database.DB.AutoMigrate(
		models.TFilm{},
		models.TLink{},
		models.TImage{},
		models.TActress{},
	); e != nil {
		fmt.Println(e.Error())
	}
}
