package handler

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	echo *echo.Echo
}

func (handler *Handler) Init() {
	handler.echo = echo.New()
	err := handler.echo.Start("127.0.0.1:7000")
	if err != nil {
		return
	}
}
