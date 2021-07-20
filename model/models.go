package model

import "image"

type User struct {
	Address   string `json:"address"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Balance   int    `json:"balance"`
}

type Admin struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Product struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	SoldNumber int    `json:"sold_number"`
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
