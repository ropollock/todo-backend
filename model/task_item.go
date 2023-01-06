package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Order      int32              `bson:"order,omitempty" json:"order,omitempty"`
	Content    string             `bson:"content,omitempty" json:"content,omitempty"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty" json:"created_ts"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty" json:"modified_ts"`
	ListID     primitive.ObjectID `bson:"list_id,omitempty" json:"list_id,omitempty"`
}
