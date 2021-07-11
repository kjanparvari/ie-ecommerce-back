package main

import (
	"ie-project-back/handler"
	"ie-project-back/model"
)

func main() {
	db := model.Database{}
	db.Init()
	server := handler.Handler{}
	server.Init()
}
