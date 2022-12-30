package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Board struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty"`
	OwnerID    primitive.ObjectID `bson:"owner_id,omitempty"`
}
