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
	"todo/services/auth"
)

func init() {
	fmt.Println("Loading config.")

	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables config.", err)
	}
	config.AppConfig = &conf

	db.Connect(conf.DBUri)
}

// go run main.go
func main() {
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

	// Routes
	e.GET("/api/healthcheck", healthcheck)
	routes.RegisterBoardsRoutes(e)
	routes.RegisterUserRoutes(e)
	routes.RegisterLoginRoutes(e)

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &auth.Claims{},
		SigningKey:              []byte(auth.GetJWTSecret()),
		TokenLookup:             "cookie:access-token,header:Authorization",
		ErrorHandlerWithContext: routes.JWTErrorChecker,
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" {
				return true
			}
			return false
		},
	}))

	e.Use(routes.TokenRefresherMiddleware)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
