package model

import "fmt"

type Database struct {
}

func (db *Database) Init() {
	fmt.Println("model is up")
}
