package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo/data"
	"todo/model"
	"todo/service"
)

type boardsController struct {
	boardService service.BoardServiceInterface
	authService  service.AuthServiceInterface
}

func BoardsController(boardService service.BoardServiceInterface, authService service.AuthServiceInterface) *boardsController {
	return &boardsController{boardService, authService}
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
	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	results, err := controller.boardService.GetBoards(userResult.ID.Hex())
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to get boards.")
	}

	return ctx.JSON(http.StatusOK, results)
}

func (controller *boardsController) FindBoardsById(ctx echo.Context) error {
	var req, err = controller.bindBoardRequest(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	// validate board is owned by current user or is admin
	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil || (req.ID != userResult.ID.Hex() && !userResult.IsAdmin) {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	result, err := controller.boardService.FindBoardById(req.ID)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (controller *boardsController) CreateBoard(ctx echo.Context) error {
	var req, err = controller.bindBoardRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	// TODO validate board name does not already exist

	var boardRecord model.Board
	boardRecord.Name = req.Name
	boardRecord.OwnerID = userResult.ID

	resultBoard, insertErr := controller.boardService.CreateBoard(&boardRecord)

	if insertErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create board.")
	}

	return ctx.JSON(http.StatusOK, resultBoard)
}

func (controller *boardsController) UpdateBoard(ctx echo.Context) error {
	var req, err = controller.bindBoardRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	boardRecord, err := controller.boardService.FindBoardById(req.ID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "board not found.")
	}

	if boardRecord.OwnerID != userResult.ID && !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	boardRecord.Name = req.Name

	resultBoard, updateErr := controller.boardService.UpdateBoard(&boardRecord)

	if updateErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to update board.")
	}

	return ctx.JSON(http.StatusOK, resultBoard)
}

func (controller *boardsController) DeleteBoard(ctx echo.Context) error {
	req, err := controller.bindBoardRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil || (req.ID != userResult.ID.Hex() && !userResult.IsAdmin) {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	var boardRecord model.Board
	boardRecord.ID, err = data.StringToObjectID(req.ID)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	deleteErr := controller.boardService.DeleteBoard(&boardRecord)

	if deleteErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete board.")
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (controller *boardsController) bindBoardRequest(ctx echo.Context) (*model.BoardRequest, error) {
	var req model.BoardRequest

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
