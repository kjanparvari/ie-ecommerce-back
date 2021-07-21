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
	ProductName       string `json:"product_name"`
	SoldNumber        int    `json:"sold_number"`
	CustomerEmail     string `json:"customer_email"`
	CustomerFirstname string `json:"customer_firstname"`
	CustomerLastname  string `json:"customer_lastname"`
	CustomerAddress   string `json:"customer_address"`
	Amount            int    `json:"amount"`
	Date              string `json:"date"`
	TracingCode       string `json:"tracing_code"`
	Status            string `json:"status"`
}
