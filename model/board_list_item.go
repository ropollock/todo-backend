package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BoardList struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Order      int32              `bson:"order,omitempty" json:"order,omitempty"`
	CreatedTS  time.Time          `bson:"created_ts,omitempty" json:"created_ts"`
	ModifiedTS time.Time          `bson:"modified_ts,omitempty" json:"modified_ts"`
	BoardID    primitive.ObjectID `bson:"board_id,omitempty" json:"board_id,omitempty"`
}
