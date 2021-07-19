package model

import "image"

type User struct {
	Address   string `json:"address"`
	Email     string `json:"email"`
	Password  string `json:"_"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Balance   int    `json:"balance"`
}

type Admin struct {
	Email    string
	Password string
}

type Product struct {
	Name       string
	Category   string
	Price      int
	Stock      int
	SoldNumber int
	image      image.Image
}

type Category struct {
	Name string
}

type Receipt struct {
	ProductName       string
	SoldNumber        int
	CustomerEmail     string
	CustomerFirstname string
	CustomerLastname  string
	CustomerAddress   string
	Amount            int
	Date              string
	TracingCode       string
	Status            string
}
