package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo/data"
	"todo/model"
	"todo/service"
)

type listsController struct {
	listService  service.ListServiceInterface
	authService  service.AuthServiceInterface
	boardService service.BoardServiceInterface
}

func ListsController(listService service.ListServiceInterface, authService service.AuthServiceInterface, boardService service.BoardServiceInterface) *listsController {
	return &listsController{listService, authService, boardService}
}

func (controller *listsController) RegisterListsRoutes(e *echo.Echo) {
	e.GET("/boards/:board_id/lists", controller.GetLists)
	e.GET("/boards/:board_id/lists/:id", controller.FindListById)
	e.POST("/boards/:board_id/lists", controller.CreateList)
	e.PUT("/boards/:board_id/lists/:id", controller.UpdateList)
	e.DELETE("/boards/:board_id/lists/:id", controller.DeleteList)
	fmt.Println("Registered /lists routes.")
}

func (controller *listsController) GetLists(ctx echo.Context) error {
	var req, err = controller.bindListRequest(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	boardRecord, err := controller.boardService.FindBoardById(req.BoardID)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "board not found.")
	}

	if boardRecord.OwnerID != userResult.ID && !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	results, err := controller.listService.GetLists(req.BoardID)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to get lists.")
	}

	return ctx.JSON(http.StatusOK, results)
}

func (controller *listsController) FindListById(ctx echo.Context) error {
	var req, err = controller.bindListRequest(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	listRecord, err := controller.listService.FindListById(req.ID)
	if err != nil || listRecord.BoardID.Hex() != req.BoardID {
		return ctx.String(http.StatusBadRequest, "list not found.")
	}

	boardRecord, err := controller.boardService.FindBoardById(listRecord.BoardID.Hex())

	if err != nil {
		return ctx.String(http.StatusBadRequest, "board not found.")
	}

	if boardRecord.OwnerID != userResult.ID && !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	result, err := controller.listService.FindListById(req.ID)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (controller *listsController) CreateList(ctx echo.Context) error {
	var req, err = controller.bindListRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	boardResult, err := controller.boardService.FindBoardById(req.BoardID)

	if err != nil || boardResult.ID.Hex() != req.BoardID {
		return ctx.String(http.StatusBadRequest, "board not found.")
	}

	if boardResult.OwnerID != userResult.ID && !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	var listRecord model.BoardList
	if len(listRecord.Name) > 100 {
		listRecord.Name = req.Name[0:100]
	} else {
		listRecord.Name = req.Name
	}
	listRecord.BoardID = boardResult.ID
	listRecord.Order = req.Order

	resultBoard, insertErr := controller.listService.CreateList(&listRecord)

	if insertErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create list.")
	}

	return ctx.JSON(http.StatusOK, resultBoard)
}

func (controller *listsController) UpdateList(ctx echo.Context) error {
	var req, err = controller.bindListRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	listRecord, err := controller.listService.FindListById(req.ID)
	if err != nil || listRecord.BoardID.Hex() != req.BoardID {
		return ctx.String(http.StatusBadRequest, "list not found.")
	}

	boardRecord, err := controller.boardService.FindBoardById(listRecord.BoardID.Hex())

	if err != nil {
		return ctx.String(http.StatusBadRequest, "board not found.")
	}

	if boardRecord.OwnerID != userResult.ID && !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	if len(listRecord.Name) > 100 {
		listRecord.Name = req.Name[0:100]
	} else {
		listRecord.Name = req.Name
	}
	listRecord.Order = req.Order

	resultList, updateErr := controller.listService.UpdateList(&listRecord)

	if updateErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to update list.")
	}

	return ctx.JSON(http.StatusOK, resultList)
}

func (controller *listsController) DeleteList(ctx echo.Context) error {
	req, err := controller.bindListRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	listResult, err := controller.listService.FindListById(req.ID)
	if err != nil || listResult.BoardID.Hex() != req.BoardID {
		return ctx.String(http.StatusBadRequest, "list not found.")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	boardRecord, err := controller.boardService.FindBoardById(listResult.BoardID.Hex())

	if err != nil {
		return ctx.String(http.StatusBadRequest, "board not found.")
	}

	if boardRecord.OwnerID != userResult.ID && !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	var listRecord model.BoardList
	listRecord.ID, err = data.StringToObjectID(req.ID)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	deleteErr := controller.listService.DeleteList(&listRecord)

	if deleteErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete list.")
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (controller *listsController) bindListRequest(ctx echo.Context) (*model.ListRequest, error) {
	var req model.ListRequest

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
