package main

import (
	"gin_movie/model"
	"log"
)

func main(){
	//insertActress()
	insertFilm()
	//insertLink()
	//insertImage()
}

func insertImage(){
	datas := []model.TImage{
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
		{Name: "张三image", FilmId: 1},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}
func insertLink(){
	datas := []model.TLink{
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
		{Name: "张三link", FilmId: 1, Size: "2GB", ShareDate: utils.StrToDate("2022-02-10")},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}

func insertFilm(){
	datas := []model.TFilm{
		{Name: "李四打虎", ActressId: 2,  Title: "李四打虎之英雄无敌", Image: "1-1.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
		{Name: "李四打狗", ActressId: 2,  Title: "李四打狗之英雄无敌", Image: "1-2.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
		{Name: "李四打狼", ActressId: 2,  Title: "李四打狼之英雄无敌", Image: "1-3.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
		{Name: "李四打熊", ActressId: 2,  Title: "李四打熊之英雄无敌", Image: "1-4.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
		{Name: "李四无敌", ActressId: 2,  Title: "李四打人之英雄无敌", Image: "1-5.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
		{Name: "李四无敌一", ActressId: 2,  Title: "李四打人之英雄无敌", Image: "1-6.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
		{Name: "李四无敌二", ActressId: 2,  Title: "李四打人之英雄无敌", Image: "1-7.jpg", Length: "90", ReleaseDate: utils.StrToDate("2022-01-01")},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}

func insertActress(){
	datas := []model.TActress{
		{Name: "张三", Age: 31, Height: "160cm", Cup: "C", Birthday: utils.StrToDate("1990-01-01")},
		{Name: "李四", Age: 32, Height: "165cm", Cup: "D", Birthday: utils.StrToDate("1989-01-01")},
		{Name: "王二", Age: 24, Height: "170cm", Cup: "B", Birthday: utils.StrToDate("1997-01-01")},
		{Name: "麻子", Age: 28, Height: "173cm", Cup: "B", Birthday: utils.StrToDate("1993-01-01")},
		{Name: "王五", Age: 22, Height: "162cm", Cup: "E", Birthday: utils.StrToDate("1999-01-01")},
		{Name: "赵六", Age: 18, Height: "155cm", Cup: "D", Birthday: utils.StrToDate("2003-01-01")},
		{Name: "张七", Age: 23, Height: "157cm", Cup: "C", Birthday: utils.StrToDate("1998-01-01")},
		{Name: "杨八", Age: 26, Height: "158cm", Cup: "B", Birthday: utils.StrToDate("1995-01-01")},
		{Name: "高九", Age: 19, Height: "171cm", Cup: "C", Birthday: utils.StrToDate("2002-01-01")},
		{Name: "黄十", Age: 20, Height: "169cm", Cup: "B", Birthday: utils.StrToDate("2001-01-01")},
	}
	result := model.DB.Create(datas)
	if result.Error != nil{
		log.Fatal(result.Error)
	}
}
