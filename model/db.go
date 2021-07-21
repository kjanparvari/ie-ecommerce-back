package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
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
	// the bellow commented code creates tables in database
	db.postgres.AutoMigrate(User{}, Admin{}, Product{}, Category{}, Receipt{})
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
func (db *Database) ModifyProduct(name string, category string, price int, stock int) {
	db.postgres.Model(Product{}).Where("name = ?", name).Updates(Product{Category: category, Price: price, Stock: stock})
}

// InsertCategory COMPLETE
func (db *Database) AddProduct(name string, category string, price int, stock int, soldNumber int) {
	products := make([]Product, 10)
	db.postgres.Find(&products, "Name =?", name)
	if (len(products)) > 0 {
		print("there is another category with same name")
		return
	}
	products = []Product{
		{Name: name, Category: category, Price: price, Stock: stock, SoldNumber: soldNumber},
	}
	for _, prods := range products {
		db.postgres.Create(&prods)
	}
}

// InsertCategory COMPLETE
func (db *Database) AddCategory(categoryName string) int {
	categories := make([]Category, 10)
	db.postgres.Find(&categories, "Name =?", categoryName)
	if (len(categories)) > 0 {
		print("there is another category with same name")
		return 0
	}
	categories = []Category{
		{Name: categoryName},
	}
	for _, categs := range categories {
		db.postgres.Create(&categs)
	}
	return 1
}

func (db *Database) GetReceipt(email string) []Receipt {
	receipts := make([]Receipt, 10)
	if len(email) == 0 {
		db.postgres.Find(&receipts)
	} else {
		db.postgres.Where("customerEmail = ?", email).Find(&receipts)
	}
	return receipts
}

// RiseBalance COMPLETE
func (db *Database) RiseBalance(email string, amount int) {
	users := make([]User, 20)
	db.postgres.Where("email = ?", email).Find(&users)
	db.postgres.Model(User{}).Where("email = ?", email).Updates(User{Balance: users[0].Balance + amount})
}

// ModifyCategory COMPLETE
func (db *Database) ModifyCategory(newName string, oldName string) int {
	categories := make([]Category, 20)
	db.postgres.Where("name = ?", newName).Find(&categories)
	if len(categories) > 0 {
		return 0
	}
	db.postgres.Model(Category{}).Where("name = ?", newName).Updates(Category{Name: newName})
	db.postgres.Model(Product{}).Where("category = ?", oldName).Updates(Product{Category: newName})
	return 1
}

// ExistCategory COMPLETE
func (db *Database) ExistCategory(name string) int {
	categories := make([]Category, 20)
	db.postgres.Where("name = ?", name).Find(&categories)
	if len(categories) == 0 {
		return 0
	}
	return 1
}

// GetCategories COMPLETE
func (db *Database) GetCategories() []string {
	allCategories := make([]Category, 20)
	db.postgres.Find(&allCategories)
	result := make([]string, 0)
	for _, cat := range allCategories {
		if cat.Name != "" {
			result = append(result, cat.Name)
		}
	}
	return result
}

// DeleteCategory COMPLETE
func (db *Database) DeleteCategory(categoryName string) {
	db.postgres.Model(Product{}).Where("category = ?", categoryName).Updates(Product{Category: "NO CATEGORY"})
	db.postgres.Where("name = ?", categoryName).Delete(&Category{})
}

// SeeAllReceipt COMPLETE
func (db *Database) SeeAllReceipt() []Receipt {
	receipts := make([]Receipt, 10)
	result := db.postgres.Find(&receipts)
	if result.Error != nil {
		panic(result.Error)
	}
	return receipts
}

// SeeReceiptByCode COMPLETE
func (db *Database) SeeReceiptByCode(code string) []Receipt {
	receipts := make([]Receipt, 10)
	result := db.postgres.Where("tracingCode=?", code).Find(&receipts)
	if result.Error != nil {
		panic(result.Error)
	}
	return receipts
}

// GetProductSort COMPLETE
func (db *Database) GetProductSort(name string, sortType string, categories []string, maxPrice int, minPrice int) []Product {
	products := make([]Product, 10)
	arrayProducts := make([]Product, 0)
	var result *gorm.DB
	//fmt.Println(sortType)
	//fmt.Println(categories)
	//fmt.Println(maxPrice)
	//fmt.Println(minPrice)
	if len(name) == 0 {
		result = db.postgres.Order(sortType).Where("price>? AND price<?", minPrice, maxPrice).Find(&products)
	} else {
		result = db.postgres.Order(sortType).Where("price>? AND price<? AND name like ?", minPrice, maxPrice, name).Find(&products)
	}
	if len(categories) == 0 {
		return products
	}
	for _, prods := range products {
		for _, categs := range categories {
			if prods.Category == categs {
				arrayProducts = append(arrayProducts, prods)
			}
		}
	}
	if result.Error != nil {
		panic(result.Error)
	}
	return arrayProducts
}

// ModifyUser COMPLETE
func (db *Database) ModifyUser(email string, address string, password string, firstName string, lastName string, balance int) {
	db.postgres.Model(User{}).Where("email = ?", email).Updates(User{Address: address, Password: password, Firstname: firstName, Lastname: lastName, Balance: balance})
}

// DeleteProduct COMPLETE
func (db *Database) DeleteProduct(name string) {
	db.postgres.Where("name = ?", name).Delete(&Product{})
}

// AddReceipt COMPLETE
func (db *Database) AddReceipt(productName string, soldNumber int, customerEmail string, customerFirstname string, customerLastname string, customerAddress string, amount int, date string, tracingCode string, status string) {
	receipts := []Receipt{
		{ProductName: productName, SoldNumber: soldNumber, CustomerAddress: customerAddress, CustomerEmail: customerEmail, CustomerFirstname: customerFirstname, CustomerLastname: customerLastname, Amount: amount, Date: date, TracingCode: tracingCode, Status: status},
	}
	for _, rescps := range receipts {
		db.postgres.Create(&rescps)
	}
}

// ChangeReceiptStatus COMPLETE
func (db *Database) ChangeReceiptStatus(code string, status string) {
	db.postgres.Model(Receipt{}).Where("tracingCode = ?", code).Updates(Receipt{Status: status})
}

func (db *Database) seeClientReceipt(email string) []Receipt {
	receipts := make([]Receipt, 10)
	result := db.postgres.Where("costumerEmail=?", email).Find(&receipts)
	if result.Error != nil {
		panic(result.Error)
	}
	return receipts
}

func (db *Database) AddUser(email string, password string, firstname string, lastname string, balance int, Address string) (int, string) {
	log.Println("[Database]: request to add user: ", email)
	users := make([]User, 10)
	db.postgres.Find(&users, "Email =?", email)
	if (len(users)) > 0 {
		log.Println("there is another users with same email")
		return -1, "there is another users with same email"
	}
	users = []User{
		{Email: email, Password: password, Firstname: firstname, Lastname: lastname, Balance: balance, Address: Address},
	}
	for _, us := range users {
		db.postgres.Create(&us)
	}
	return 1, "done"
}

func (db *Database) BuyProduct(email string, name string, number int) string {
	products := make([]Product, 10)
	db.postgres.Find(&products, "name =?", name)
	users := make([]User, 10)
	db.postgres.Find(&users, "email = ?", email)
	if len(products) == 0 {
		return "There Is No Such Product"
	}
	if products[0].Stock < number {
		return "There Is Not Enough Stocks"
	}
	if users[0].Balance < number*products[0].Price {
		return "Not Enough Money"
	}
	db.postgres.Model(Product{}).Where("name = ?", name).Updates(Product{Stock: products[0].Stock - number, SoldNumber: products[0].SoldNumber + number})
	db.postgres.Model(User{}).Where("email = ?", email).Updates(User{Balance: users[0].Balance - number*products[0].Price})
	// Using time.Now() function.
	dt := time.Now()
	formatedTime := dt.Format(time.RFC1123)
	db.AddReceipt(name, number, email, users[0].Firstname, users[0].Lastname, users[0].Address, number*products[0].Price, formatedTime, formatedTime+email, "در حال انجام")
	return "Done"
}

func (db *Database) GetUser(email string) *User {
	var user User
	db.postgres.Where("email = ?", email).First(&user)
	if user.Email == "" {
		return nil
	}
	return &user
}

func (db *Database) GetAdmin(email string) *Admin {
	var admin Admin
	db.postgres.Where("email = ?", email).First(&admin)
	if admin.Email == "" {
		return nil
	}
	return &admin
}
