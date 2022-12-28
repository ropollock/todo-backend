package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo/db"
	"todo/services"
)

type UserRequest struct {
	ID       string `param:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

func RegisterUserRoutes(e *echo.Echo) {
	e.GET("/users", GetUsers)
	e.GET("/users/:id", FindUserById)
	e.POST("/users", CreateUser)
	e.DELETE("/users/:id", DeleteUser)
	fmt.Println("Registered /users routes.")
}

func GetUsers(ctx echo.Context) error {
	results, err := services.GetUsers()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to get users.")
	}
	return ctx.JSON(http.StatusOK, results)
}

func FindUserById(ctx echo.Context) error {
	var req, err = bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	result, err := services.FindUserById(req.ID)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	result.Password = ""

	return ctx.JSON(http.StatusOK, result)
}

func CreateUser(ctx echo.Context) error {
	var req, err = bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	var userRecord services.User
	userRecord.Name = req.Name
	userRecord.Username = req.Username
	hashedPassword, hashErr := services.HashPassword(req.Password)

	if hashErr != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userRecord.Password = hashedPassword
	userRecord.IsAdmin = req.IsAdmin

	result, insertErr := services.CreateUser(&userRecord)

	if insertErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create user.")
	}

	return ctx.JSON(http.StatusOK, result)
}

func DeleteUser(ctx echo.Context) error {
	var req, err = bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	var userRecord services.User
	userRecord.ID, err = db.StringToObjectID(req.ID)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	result, deleteErr := services.DeleteUser(&userRecord)

	if deleteErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete user.")
	}

	return ctx.JSON(http.StatusOK, result)
}

func bindUserRequest(ctx echo.Context) (*UserRequest, error) {
	var req UserRequest

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
