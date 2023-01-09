package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo/data"
	"todo/model"
	"todo/service"
)

type tasksController struct {
	taskService  service.TaskServiceInterface
	authService  service.AuthServiceInterface
	boardService service.BoardServiceInterface
	listService  service.ListServiceInterface
}

func TasksController(taskService service.TaskServiceInterface, authService service.AuthServiceInterface,
	boardService service.BoardServiceInterface, listService service.ListServiceInterface) *tasksController {
	return &tasksController{taskService, authService, boardService, listService}
}

func (controller *tasksController) RegisterTasksRoutes(e *echo.Echo) {
	e.GET("/boards/:board_id/lists/:list_id/tasks", controller.GetTasks)
	e.GET("/boards/:board_id/lists/:list_id/tasks/:id", controller.FindTaskById)
	e.POST("/boards/:board_id/lists/:list_id/tasks", controller.CreateTask)
	e.PUT("/boards/:board_id/lists/:list_id/tasks/:id", controller.UpdateTask)
	e.DELETE("/boards/:board_id/lists/:list_id/tasks/:id", controller.DeleteTask)
	fmt.Println("Registered /tasks routes.")
}

func (controller *tasksController) GetTasks(ctx echo.Context) error {
	var req, err = controller.bindTaskRequest(ctx)
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

	results, err := controller.taskService.GetTasks(req.ListID)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to get tasks.")
	}

	return ctx.JSON(http.StatusOK, results)
}

func (controller *tasksController) FindTaskById(ctx echo.Context) error {
	var req, err = controller.bindTaskRequest(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	listRecord, err := controller.listService.FindListById(req.ListID)
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

	result, err := controller.taskService.FindTaskById(req.ID)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (controller *tasksController) CreateTask(ctx echo.Context) error {
	var req, err = controller.bindTaskRequest(ctx)

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

	listResult, err := controller.listService.FindListById(req.ListID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "list not found.")
	}

	var taskRecord model.Task
	if len(taskRecord.Name) > 100 {
		taskRecord.Name = req.Name[0:100]
	} else {
		taskRecord.Name = req.Name
	}
	taskRecord.Content = req.Content
	taskRecord.ListID = listResult.ID
	taskRecord.Order = req.Order

	resultBoard, insertErr := controller.taskService.CreateTask(&taskRecord)

	if insertErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create list.")
	}

	return ctx.JSON(http.StatusOK, resultBoard)
}

func (controller *tasksController) UpdateTask(ctx echo.Context) error {
	var req, err = controller.bindTaskRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	listRecord, err := controller.listService.FindListById(req.ListID)
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

	taskRecord, err := controller.taskService.FindTaskById(req.ID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "task not found.")
	}

	if len(taskRecord.Name) > 100 {
		taskRecord.Name = req.Name[0:100]
	} else {
		taskRecord.Name = req.Name
	}
	taskRecord.Content = req.Content
	taskRecord.Order = req.Order

	resultTask, updateErr := controller.taskService.UpdateTask(&taskRecord)

	if updateErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to update task.")
	}

	return ctx.JSON(http.StatusOK, resultTask)
}

func (controller *tasksController) DeleteTask(ctx echo.Context) error {
	req, err := controller.bindTaskRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	listResult, err := controller.listService.FindListById(req.ListID)
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

	var taskRecord model.Task
	taskRecord.ID, err = data.StringToObjectID(req.ID)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	deleteErr := controller.taskService.DeleteTask(&taskRecord)

	if deleteErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete task.")
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (controller *tasksController) bindTaskRequest(ctx echo.Context) (*model.TaskRequest, error) {
	var req model.TaskRequest

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
