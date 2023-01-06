package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Board struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name,omitempty" json:"name"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty" json:"created_ts"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty" json:"modified_ts"`
	OwnerID    primitive.ObjectID `bson:"owner_id,omitempty" json:"owner_id"`
}
