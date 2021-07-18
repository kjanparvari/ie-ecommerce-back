package main

import (
	"ie-project-back/handler"
	"ie-project-back/model"
)

func main() {
	db := model.Database{}
	db.Init()
	//db.InsertCategory("دسته بندی یک")
	//db.InsertCategory("دسته بندی دو")
	//db.InsertCategory("دسته بندی سه")
	//db.InsertCategory("دسته بندی چهار")
	//db.InsertCategory("categ01")
	//db.InsertProduct("soup","categ01",120, 100,10)
	//db.InsertProduct("آش","دسته بندی پنج",150, 20,30)
	//db.InsertProduct("مرغ","دسته بندی پنج",300, 50,40)
	//var trash = make([]string, 10)
	//trash = append(trash, "categ01")
	//db.GetProductSort("Price asc", trash, 290, 100)
	//db.GetCategories()
	//db.InsertUser("saeed.maroof@ymail.com", "12345678", "saeed", "maroof", 0, "tehran")
	server := handler.Handler{}
	server.Init(&db)
}
