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
	//db.InsertCategory("دسته بندی پنج")
	db.GetCategories()
	server := handler.Handler{}
	server.Init(&db)
}
