package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"todo/config"
	"todo/routes"
)

var (
	ctx         context.Context
	mongoclient *mongo.Client
)

func init() {
	fmt.Println("Loading config.")

	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()
	mongoconn := options.Client().ApplyURI(conf.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected.")
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

	// Routes
	e.GET("/api/healthcheck", healthcheck)
	e.GET("/boards", routes.Boards)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
