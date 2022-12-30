package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Order      int32              `bson:"order,omitempty"`
	Content    string             `bson:"content,omitempty"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty"`
	ListID     primitive.ObjectID `bson:"list_id,omitempty"`
}
