package routes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Username    string             `bson:"username,omitempty"`
	Password    string             `bson:"password,omitempty"`
	CreatedTS   time.Time          `bson:"created_ts,omitempty"`
	LastLoginTS time.Time          `bson:"last_login_ts,omitempty"`
}
