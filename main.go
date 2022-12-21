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
	"todo/db"
	"todo/routes"
)

var (
	SERVER *echo.Echo
)

func init() {
	fmt.Println("Loading config.")

	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables config.", err)
	}

	db.Connect(conf.DBUri)
}

// go run main.go
func main() {
	e := echo.New()
	SERVER = e
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

	// Routes
	e.GET("/api/healthcheck", healthcheck)
	routes.RegisterBoardsRoutes(e)
	routes.RegisterUserRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
