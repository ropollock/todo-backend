package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Board struct {
	Uuid string
	Name string
}

func Boards(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "GET, Boards!")
}
