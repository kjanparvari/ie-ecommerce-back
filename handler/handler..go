package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ie-project-back/model"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	echo      *echo.Echo
	db        *model.Database
	secretKey string
}

func HashFunc(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (handler *Handler) Init(db *model.Database) {
	handler.db = db
	handler.echo = echo.New()
	handler.secretKey = "secret-key"
	handler.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		//AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	handler.echo.GET("/api/categories/all", handler.handleGetCategories)
	handler.echo.POST("/api/signup", handler.handleSignup)
	handler.echo.POST("/api/login", handler.handleLogin)
	handler.echo.GET("/api/user", handler.handleGetUser)
	handler.echo.GET("/api/admin", handler.handleGetAdmin)
	handler.echo.POST("/api/logout", handler.handleLogout)
	handler.echo.GET("/api/products", handler.handlerGetProducts)
	handler.echo.POST("/api/products/modify", handler.handlerModifyProduct)
	handler.echo.POST("/api/products/delete", handler.handlerDeleteProduct)
	handler.echo.POST("/api/user/modify", handler.handlerModifyUser)
	handler.echo.POST("/api/categories/modify", handler.handlerModifyCategory)
	handler.echo.POST("/api/categories/add", handler.handlerAddCategory)
	handler.echo.POST("/api/categories/delete", handler.handlerDeleteCategory)
	handler.echo.POST("/api/user/risePrice", handler.handlerRisePrice)
	handler.echo.POST("/api/admin/products/add", handler.handlerAddProduct)

	err := handler.echo.Start("127.0.0.1:7000")
	if err != nil {
		return
	}
}
func (handler *Handler) handlerRisePrice(context echo.Context) error {
	email := context.QueryParam("email")
	amount, _ := strconv.Atoi(context.QueryParam("amount"))
	handler.db.RiseBalance(email, amount)
	return context.String(http.StatusOK, "OK")
}
func (handler *Handler) handlerModifyProduct(context echo.Context) error {
	name := context.QueryParam("name")
	category := context.QueryParam("category")
	price, _ := strconv.Atoi(context.QueryParam("price"))
	stock, _ := strconv.Atoi(context.QueryParam("stock"))
	handler.db.ModifyProduct(name, category, price, stock)
	return context.String(http.StatusOK, "OK")
}
func (handler *Handler) handlerDeleteProduct(context echo.Context) error {
	name := context.QueryParam("name")
	handler.db.DeleteProduct(name)
	return context.String(http.StatusOK, "OK")
}
func (handler *Handler) handlerDeleteCategory(context echo.Context) error {
	name := context.QueryParam("name")
	handler.db.DeleteCategory(name)
	return context.String(http.StatusOK, "OK")
}
func (handler *Handler) handlerAddCategory(context echo.Context) error {
	name := context.QueryParam("name")
	correctness := handler.db.AddCategory(name)
	if correctness == 0 {
		return context.String(http.StatusOK, "Can Not Add Category")
	}
	return context.String(http.StatusOK, "OK")
}

func (handler *Handler) handlerAddProduct(context echo.Context) error {
	name := context.QueryParam("name")
	category := context.QueryParam("category")
	price, _ := strconv.Atoi(context.QueryParam("price"))
	stock, _ := strconv.Atoi(context.QueryParam("stock"))
	exist := handler.db.ExistCategory(category)
	returnString := "OK"
	if exist == 0 {
		category = "دسته بندی نشده"
		returnString = "Category Changes"
	}
	handler.db.AddProduct(name, category, price, stock, 0)
	return context.String(http.StatusOK, returnString)

}

func (handler *Handler) handlerModifyCategory(context echo.Context) error {
	newName := context.QueryParam("newName")
	oldName := context.QueryParam("oldName")
	handler.db.ModifyCategory(newName, oldName)
	return context.String(http.StatusOK, "OK")
}

func (handler *Handler) handlerModifyUser(context echo.Context) error {
	address := context.QueryParam("address")
	email := context.QueryParam("email")
	password := context.QueryParam("password")
	firstName := context.QueryParam("firstName")
	lastName := context.QueryParam("lastName")
	balance, _ := strconv.Atoi(context.QueryParam("balance"))
	handler.db.ModifyUser(email, address, password, firstName, lastName, balance)
	return context.String(http.StatusOK, "OK")
}
func (handler *Handler) handlerGetProducts(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for Products"))
	sort := context.QueryParam("sort")
	name := context.QueryParam("name")
	categories := context.QueryParams()["category"]
	minPrice, _ := strconv.Atoi(context.QueryParam("minPrice"))
	maxPrice, _ := strconv.Atoi(context.QueryParam("maxPrice"))
	raw := handler.db.GetProductSort(name, sort, categories, maxPrice, minPrice)
	_json, err := json.Marshal(raw)
	if err != nil {
		log.Println(err)
		return context.String(http.StatusServiceUnavailable, "")
	} else {
		log.Println(fmt.Sprintf("[Server]: categories: %s", string(_json)))
		return context.String(http.StatusOK, string(_json))
	}
}

func (handler *Handler) handleGetCategories(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for categories"))
	raw := handler.db.GetCategories()
	_json, err := json.Marshal(raw)
	if err != nil {
		log.Println(err)
		return context.String(http.StatusServiceUnavailable, "")
	} else {
		log.Println(fmt.Sprintf("[Server]: categories: %s", string(_json)))
		return context.String(http.StatusOK, string(_json))
	}
}

func (handler *Handler) handleSignup(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for signup"))
	var _json map[string]string = map[string]string{}
	err := context.Bind(&_json)
	if err != nil {
		log.Println(err)
		return context.String(http.StatusBadRequest, "")
	}
	log.Println("[Server]: user info: ", _json)
	hashedStr := HashFunc(_json["password"])
	ok, msg := handler.db.AddUser(_json["email"], hashedStr, _json["firstname"], _json["lastname"], 0, _json["address"])
	if ok == -1 {
		return context.String(http.StatusBadRequest, msg)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    _json["email"],
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	if token, err := claims.SignedString([]byte(handler.secretKey)); err != nil {
		log.Println(err)
		return context.String(http.StatusInternalServerError, "could not login")
	} else {
		log.Println("[Server]: user ", _json["email"], " logged in")
		cookie := new(http.Cookie)
		cookie.Name = "jwt"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * time.Hour) // 1 day
		cookie.HttpOnly = true
		context.SetCookie(cookie)
		return context.String(http.StatusOK, "you have been registered\"")
	}
}

func (handler *Handler) handleLogin(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for login"))
	var _json map[string]string = map[string]string{}

	if err := context.Bind(&_json); err != nil {
		log.Println(err)
		return context.String(http.StatusBadRequest, "")
	}
	log.Println("[Server]: user info: ", _json)

	var _pass string
	var _email string
	if admin := handler.db.GetAdmin(_json["email"]); admin != nil {
		_email = admin.Email
		_pass = admin.Password
	} else {
		user := handler.db.GetUser(_json["email"])
		if user == nil {
			log.Println("[Server]: user not found")
			return context.String(http.StatusNotFound, "user not found")
		}
		_email = user.Email
		_pass = user.Password
	}

	if _pass != HashFunc(_json["password"]) {
		return context.String(http.StatusBadRequest, "incorrect password")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    _email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	if token, err := claims.SignedString([]byte(handler.secretKey)); err != nil {
		log.Println(err)
		return context.String(http.StatusInternalServerError, "could not login")
	} else {
		log.Println("[Server]: user ", _json["email"], " logged in")
		cookie := new(http.Cookie)
		cookie.Name = "jwt"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * time.Hour) // 1 day
		cookie.HttpOnly = true
		context.SetCookie(cookie)
		return context.String(http.StatusOK, "logged in!")
	}
}

func (handler *Handler) authenticate(context echo.Context) (bool, *jwt.Token) {
	cookie, err1 := context.Cookie("jwt")
	if err1 != nil {
		log.Println(err1)
		return false, nil
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.secretKey), nil
	})
	if err != nil {
		return false, nil
	}
	return true, token
}
func (handler *Handler) handleGetUser(context echo.Context) error {
	isAuth, token := handler.authenticate(context)
	if !isAuth {
		return context.String(http.StatusUnauthorized, "unauthenticated")
	}

	claims := token.Claims.(*jwt.StandardClaims)

	if user := handler.db.GetUser(claims.Issuer); user != nil {
		return context.JSON(http.StatusOK, *user)
	} else if admin := handler.db.GetAdmin(claims.Issuer); admin != nil {
		return context.JSON(http.StatusOK, *admin)
	} else {
		return context.String(http.StatusInternalServerError, "")
	}

}
func (handler *Handler) handleGetAdmin(context echo.Context) error {
	isAuth, token := handler.authenticate(context)
	if !isAuth {
		return context.String(http.StatusUnauthorized, "unauthenticated")
	}

	claims := token.Claims.(*jwt.StandardClaims)

	if admin := handler.db.GetAdmin(claims.Issuer); admin != nil {
		return context.JSON(http.StatusOK, *admin)
	} else {
		return context.String(http.StatusInternalServerError, "")
	}

}
func (handler *Handler) handleLogout(context echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.HttpOnly = true
	context.SetCookie(cookie)
	return context.String(http.StatusOK, "logged out")
}
