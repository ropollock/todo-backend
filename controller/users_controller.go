package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"strings"
	"todo/data"
	"todo/model"
	"todo/service"
)

type usersController struct {
	userService service.UserServiceInterface
	authService service.AuthServiceInterface
}

func UsersController(userService service.UserServiceInterface, authService service.AuthServiceInterface) *usersController {
	return &usersController{userService, authService}
}

func (controller *usersController) RegisterUserRoutes(e *echo.Echo) {
	e.GET("/users", controller.GetUsers)
	e.GET("/users/:id", controller.FindUserById)
	e.POST("/users", controller.CreateUser)
	e.DELETE("/users/:id", controller.DeleteUser)
	fmt.Println("Registered /users routes.")
}

func (controller *usersController) GetUsers(ctx echo.Context) error {
	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil || !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	results, err := controller.userService.GetUsers()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to get users.")
	}

	return ctx.JSON(http.StatusOK, results)
}

func (controller *usersController) FindUserById(ctx echo.Context) error {
	var req, err = controller.bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	reqObjectID, err := data.StringToObjectID(req.ID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil || (reqObjectID != userResult.ID && !userResult.IsAdmin) {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	result, err := controller.userService.FindUserById(req.ID)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	controller.userService.ScrubUserForAPI(&result)

	return ctx.JSON(http.StatusOK, result)
}

func (controller *usersController) CreateUser(ctx echo.Context) error {
	var req, err = controller.bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	var userRecord model.User
	userRecord.Name = strings.TrimSpace(strings.ToLower(req.Name))
	userRecord.Username = strings.TrimSpace(strings.ToLower(req.Username))

	if userRecord.Username == "" {
		return ctx.String(http.StatusBadRequest, "bad request. missing or empty username.")
	}

	if !controller.userService.ValidateUsername(userRecord.Username) {
		return ctx.String(http.StatusBadRequest, "bad request. invalid username.")
	}

	_, err = controller.userService.FindUserByUsername(userRecord.Username)
	if err == nil {
		return ctx.String(http.StatusBadRequest, "bad request. user by that username already exists.")
	}

	if userRecord.Name == "" {
		userRecord.Name = userRecord.Username
	}

	if len(userRecord.Name) > 40 {
		userRecord.Name = userRecord.Name[0:40]
	}

	userRecord.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if userRecord.Email == "" {
		return ctx.String(http.StatusBadRequest, "bad request. missing or empty email.")
	}

	_, err = mail.ParseAddress(userRecord.Email)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request. invalid email.")
	}

	if !controller.userService.ValidatePassword(req.Password) {
		return ctx.String(http.StatusBadRequest, "bad request. invalid password.")
	}

	hashedPassword, hashErr := hashPassword(req.Password)

	if hashErr != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userRecord.Password = hashedPassword
	userRecord.IsAdmin = req.IsAdmin

	resultUser, insertErr := controller.userService.CreateUser(&userRecord)

	if insertErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create user.")
	}

	controller.userService.ScrubUserForAPI(resultUser)
	return ctx.JSON(http.StatusOK, resultUser)
}

func (controller *usersController) DeleteUser(ctx echo.Context) error {
	userResult, err := controller.authService.GetCurrentUser(ctx)
	if err != nil || !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	req, err := controller.bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	var userRecord model.User
	userRecord.ID, err = data.StringToObjectID(req.ID)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	deleteErr := controller.userService.DeleteUser(&userRecord)

	if deleteErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete user.")
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (controller *usersController) bindUserRequest(ctx echo.Context) (*model.UserRequest, error) {
	var req model.UserRequest

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
