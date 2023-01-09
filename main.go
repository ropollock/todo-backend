package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
	"net/http"
	"os"
	"todo/config"
	"todo/controller"
	"todo/dao"
	"todo/data"
	"todo/model"
	"todo/service"
)

// go run main.go
func main() {
	fmt.Println("Loading config.")

	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables config.", err)
	}

	config.AppConfig = &conf

	databaseProvider := data.MongoDBProvider()
	databaseProvider.Connect(conf.DBUri)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
	}))

	// Middleware
	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(log.DEBUG),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Logger = logger

	e.Use(middleware.RequestID())
	e.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Routes
	e.GET("/api/healthcheck", healthcheck)

	userDao := dao.UserDao(databaseProvider)
	userSerivce := service.UserService(userDao)
	authService := service.AuthService(userSerivce)

	boardDao := dao.BoardDao(databaseProvider)
	boardsService := service.BoardService(boardDao)
	boardsController := controller.BoardsController(boardsService, authService)
	boardsController.RegisterBoardsRoutes(e)

	usersController := controller.UsersController(userSerivce, authService)
	usersController.RegisterUserRoutes(e)

	authController := controller.AuthController(userSerivce, authService)
	authController.RegisterLoginRoutes(e)

	listDao := dao.ListDao(databaseProvider)
	listsService := service.ListService(listDao)
	listsController := controller.ListsController(listsService, authService, boardsService)
	listsController.RegisterListsRoutes(e)

	taskDao := dao.TaskDao(databaseProvider)
	tasksService := service.TaskService(taskDao)
	tasksController := controller.TasksController(tasksService, authService, boardsService, listsService)
	tasksController.RegisterTasksRoutes(e)

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &model.Claims{},
		SigningKey:              []byte(authService.GetJWTSecret()),
		TokenLookup:             "cookie:access-token,header:Authorization",
		ErrorHandlerWithContext: authController.JWTErrorChecker,
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" {
				return true
			}
			return false
		},
	}))

	e.Use(authController.TokenRefresherMiddleware)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
