package main

import (
	"ie-project-back/model"
)
func main() {
	db := model.Database{}
	db.Init()
	db.SeeAllReceipt()
	//server := handler.Handler{}
	//server.Init()
}
