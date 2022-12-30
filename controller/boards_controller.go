package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type boardsController struct {
}

func BoardsController() *boardsController {
	return &boardsController{}
}

type BoardRequest struct {
	ID      string `param:"id" query:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func (controller *boardsController) RegisterBoardsRoutes(e *echo.Echo) {
	e.GET("/boards", controller.GetBoards)
	e.GET("/boards/:id", controller.FindBoardsById)
	e.POST("/boards", controller.CreateBoard)
	e.PUT("/boards/:id", controller.UpdateBoard)
	e.DELETE("/boards/:id", controller.DeleteBoard)
	fmt.Println("Registered /boards routes.")
}

func (controller *boardsController) GetBoards(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "GET, Boards!")
}

func (controller *boardsController) FindBoardsById(ctx echo.Context) error {
	var req BoardRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	return ctx.String(http.StatusOK, "FIND, Board!")
}

func (controller *boardsController) CreateBoard(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "CREATE, Board!")
}

func (controller *boardsController) UpdateBoard(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "UPDATE, Board!")
}

func (controller *boardsController) DeleteBoard(ctx echo.Context) error {
	return ctx.String(http.StatusNoContent, "DELETE, Board!")
}
