package main

import (
	"ie-project-back/handler"
	"ie-project-back/model"
)

func main() {
	db := model.Database{}
	db.Init()
	//db.AddReceipt("محصول یک", 2, "kjanparvari@gmail.com", "کامیار", "جان پروری", "تهران", 3, "2021-01-12", "testCode1", "درحال انجام")
	//db.AddReceipt("محصول یک", 3, "kjanparvari@gmail.com", "کامیار", "جان پروری", "تهران", 3, "2021-01-12", "testCode2", "درحال انجام")
	//db.AddReceipt("محصول دو", 1, "kjanparvari@gmail.com", "کامیار", "جان پروری", "تهران", 3, "2021-01-12", "testCode3", "درحال انجام")
	//db.AddReceipt("محصول سه", 1, "kjanparvari@gmail.com", "کامیار", "جان پروری", "تهران", 3, "2021-01-12", "testCode4", "درحال انجام")
	//db.AddReceipt("محصول چهار", 6, "kjanparvari@gmail.com", "کامیار", "جان پروری", "تهران", 3, "2021-01-12", "testCode5", "درحال انجام")
	//fmt.Println(handler.HashFunc("kj12345678"))
	//db.AddUser("c@c.co", "7afcd74c08ff96179aa4be2c4d2cc41f","کامیار","جان پروری", 500, "تهران")
	server := handler.Handler{}
	server.Init(&db)
}
