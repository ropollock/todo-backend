package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/mail"
	"strings"
	"todo/db"
	"todo/services"
	"todo/services/auth"
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
	userResult, err := auth.GetCurrentUser(ctx)
	if err != nil || !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	results, err := services.GetUsers()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to get users.")
	}

	return ctx.JSON(http.StatusOK, results)
}

func FindUserById(ctx echo.Context) error {
	var req, err = bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	reqObjectID, err := db.StringToObjectID(req.ID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := auth.GetCurrentUser(ctx)
	if err != nil || (reqObjectID != userResult.ID && !userResult.IsAdmin) {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	result, err := services.FindUserById(req.ID)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	services.ScrubUserForAPI(&result)

	return ctx.JSON(http.StatusOK, result)
}

func CreateUser(ctx echo.Context) error {
	var req, err = bindUserRequest(ctx)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	var userRecord services.User
	userRecord.Name = strings.TrimSpace(strings.ToLower(req.Name))
	userRecord.Username = strings.TrimSpace(strings.ToLower(req.Username))

	if userRecord.Username == "" {
		return ctx.String(http.StatusBadRequest, "bad request. missing or empty username.")
	}

	if !services.ValidateUsername(userRecord.Username) {
		return ctx.String(http.StatusBadRequest, "bad request. invalid username.")
	}

	_, err = services.FindUserByUsername(userRecord.Username)
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

	if !services.ValidatePassword(req.Password) {
		return ctx.String(http.StatusBadRequest, "bad request. invalid password.")
	}

	hashedPassword, hashErr := services.HashPassword(req.Password)

	if hashErr != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userRecord.Password = hashedPassword
	userRecord.IsAdmin = req.IsAdmin

	_, insertErr := services.CreateUser(&userRecord)

	if insertErr != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create user.")
	}

	resultUser, _ := services.FindUserByUsername(userRecord.Username)
	services.ScrubUserForAPI(&resultUser)
	return ctx.JSON(http.StatusOK, resultUser)
}

func DeleteUser(ctx echo.Context) error {
	userResult, err := auth.GetCurrentUser(ctx)
	if err != nil || !userResult.IsAdmin {
		return ctx.String(http.StatusUnauthorized, "user is not authorized.")
	}

	req, err := bindUserRequest(ctx)

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
