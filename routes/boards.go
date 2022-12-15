package routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type Board struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty"`
	OwnerID    primitive.ObjectID `bson:"owner_id,omitempty"`
}

func Boards(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "GET, Boards!")
}
