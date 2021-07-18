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

// InsertCategory COMPLETE
func (db *Database) InsertProduct(name string, category string, price int, stock int, soldNumber int) {
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
func (db *Database) InsertCategory(categoryName string) {
	categories := make([]Category, 10)
	db.postgres.Find(&categories, "Name =?", categoryName)
	if (len(categories)) > 0 {
		print("there is another category with same name")
		return
	}
	categories = []Category{
		{Name: categoryName},
	}
	for _, categs := range categories {
		db.postgres.Create(&categs)
	}
}

// ModifyCategory COMPLETE
func (db *Database) RiseBalance(email string, amount int) {
	users := make([]User, 20)
	db.postgres.Where("email = ?", email).Find(&users)
	db.postgres.Model(User{}).Where("email = ?", email).Updates(User{Balance: users[0].Balance + amount})
}

// ModifyCategory COMPLETE
func (db *Database) ModifyCategory(newName string, oldName string) {
	db.postgres.Model(Category{}).Where("name = ?", oldName).Updates(Category{Name: newName})
	db.postgres.Model(Product{}).Where("category = ?", oldName).Updates(Product{Category: newName})
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
func (db *Database) GetProductSort(sortType string, categories []string, maxPrice int, minPrice int) []Product {
	products := make([]Product, 10)
	//db.postgres.Raw("SELECT * FROM products WHERE category IN ?", []string{"دسته بندی پنج", "jinzhu 2"}).Scan(&products)
	//db.postgres.Not(map[string]interface{}{"category": []string{"دسته بندی سه", "دسته بندی دو"}}).Find(&products)

	//fmt.Println(products)
	result := db.postgres.Order(sortType).Where("price>? AND price<?", minPrice, maxPrice).Find(&products)
	if result.Error != nil {
		panic(result.Error)
	}
	return products
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

func (db *Database) InsertUser(email string, password string, firstname string, lastname string, balance int, Address string) (int, string) {
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
	db.postgres.Model(Product{}).Where("name = ?", name).Updates(Product{Stock: products[0].Stock - number})
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
