package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BoardRequest struct {
	ID      string `param:"id" query:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func RegisterBoardsRoutes(e *echo.Echo) {
	e.GET("/boards", GetBoards)
	e.GET("/boards/:id", FindBoardsById)
	e.POST("/boards", CreateBoard)
	e.PUT("/boards/:id", UpdateBoard)
	e.DELETE("/boards/:id", DeleteBoard)
	fmt.Println("Registered /boards routes.")
}

func GetBoards(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "GET, Boards!")
}

func FindBoardsById(ctx echo.Context) error {
	var req BoardRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	return ctx.String(http.StatusOK, "FIND, Board!")
}

func CreateBoard(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "CREATE, Board!")
}

func UpdateBoard(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "UPDATE, Board!")
}

func DeleteBoard(ctx echo.Context) error {
	return ctx.String(http.StatusNoContent, "DELETE, Board!")
}
