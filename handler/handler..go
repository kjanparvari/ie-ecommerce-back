package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ie-project-back/model"
	"log"
	"net/http"
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
	handler.echo.POST("/api/logout", handler.handleLogout)
	err := handler.echo.Start("127.0.0.1:7000")
	if err != nil {
		return
	}
}


func (handler *Handler) handlerGetProducts(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for categories"))
	sort := context.QueryParam("sort")
	minPrice , _ := strconv.Atoi(context.QueryParam("minPrice"))
	maxPrice , _:= strconv.Atoi(context.QueryParam("maxPrice"))

	raw := handler.db.GetProductSort(sort,,maxPrice,minPrice)
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
	var json map[string]string = map[string]string{}
	err := context.Bind(&json)
	if err != nil {
		log.Println(err)
		return context.String(http.StatusBadRequest, "")
	}
	log.Println("[Server]: user info: ", json)
	hashedStr := HashFunc(json["password"])
	ok, msg := handler.db.InsertUser(json["email"], hashedStr, json["firstname"], json["lastname"], 0, json["address"])
	if ok == -1 {
		return context.String(http.StatusBadRequest, msg)
	}
	return context.String(http.StatusOK, "you have been registered")
}

func (handler *Handler) handleLogin(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for login"))
	var json map[string]string = map[string]string{}

	if err := context.Bind(&json); err != nil {
		log.Println(err)
		return context.String(http.StatusBadRequest, "")
	}
	log.Println("[Server]: user info: ", json)
	user := handler.db.GetUser(json["email"])
	if user == nil {
		log.Println("[Server]: user not found")
		return context.String(http.StatusNotFound, "user not found")
	}
	if user.Password != HashFunc(json["password"]) {
		return context.String(http.StatusBadRequest, "incorrect password")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	if token, err := claims.SignedString([]byte(handler.secretKey)); err != nil {
		log.Println(err)
		return context.String(http.StatusInternalServerError, "could not login")
	} else {
		log.Println("[Server]: user ", json["email"], " logged in")
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
	user := handler.db.GetUser(claims.Issuer)

	return context.JSON(http.StatusOK, *user)
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
