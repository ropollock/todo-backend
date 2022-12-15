package routes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BoardList struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Order      int32              `bson:"order,omitempty"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty"`
	BoardID    primitive.ObjectID `bson:"board_id,omitempty"`
}
