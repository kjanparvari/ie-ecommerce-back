package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Database struct {
	postgres *gorm.DB
}

func (db *Database) Init() {
	var err error
	// database should be created in pgAdmin
	db.postgres, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=ie-project-db sslmode=disable password=62442")
	if err != nil {
		log.Println("[Database]: failed to connect database")
		log.Println(err)
		os.Exit(-1)
	}
	log.Println("[Database]: db is up")
	// the bellow commented code creates database
	//db.createTables()
}

func (db *Database) createTables() {
	db.postgres.CreateTable(User{})
	db.postgres.CreateTable(Admin{})
	db.postgres.CreateTable(Product{})
	db.postgres.CreateTable(Category{})
	db.postgres.CreateTable(Receipt{})
}

func (db *Database) tmp() {
	db.postgres.Exec("")
}
