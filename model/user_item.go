package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Username    string             `bson:"username,omitempty" json:"username,omitempty"`
	Password    string             `bson:"password,omitempty" json:"-"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	CreatedTS   time.Time          `bson:"created_ts,omitempty" json:"created_ts"`
	LastLoginTS time.Time          `bson:"last_login_ts,omitempty" json:"last_login_ts"`
	IsAdmin     bool               `bson:"is_admin" json:"is_admin,omitempty"`
}
