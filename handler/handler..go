package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"ie-project-back/model"
	"log"
	"net/http"
)

type Handler struct {
	echo *echo.Echo
	db   *model.Database
}

func (handler *Handler) Init(db *model.Database) {
	handler.db = db
	handler.echo = echo.New()
	handler.echo.GET("/api/categories/all", handler.handleGetCategories)
	handler.echo.POST("/api/signup", handler.handleSignup)
	handler.echo.POST("/api/login", handler.handleLogin)
	err := handler.echo.Start("127.0.0.1:7000")
	if err != nil {
		return
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
func NewSHA256(str string) string {
	data := []byte(str)
	hash := sha256.Sum256(data)
	return string(hash[:])
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
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(json["password"]), 14)
	hashedStr := hex.EncodeToString(hashedBytes)
	ok, msg := handler.db.InsertUser(json["email"], hashedStr, json["firstname"], json["lastname"], 0, json["address"])
	if ok == -1 {
		return context.String(http.StatusAccepted, msg)
	}
	return context.String(http.StatusOK, "you have been registered")
}

func (handler *Handler) handleLogin(context echo.Context) error {
	log.Println(fmt.Sprintf("[Server]: requested for login"))
	var json map[string]string = map[string]string{}
	err := context.Bind(&json)
	if err != nil {
		log.Println(err)
		return context.String(http.StatusBadRequest, "")
	}
	log.Println("[Server]: user info: ", json)
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(json["password"]), 14)
	hashedStr := hex.EncodeToString(hashedBytes)
	ok, msg := handler.db.InsertUser(json["email"], hashedStr, json["firstname"], json["lastname"], 0, json["address"])
	if ok == -1 {
		return context.String(http.StatusAccepted, msg)
	}
	return context.String(http.StatusOK, "you have been registered")
}
