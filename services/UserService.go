package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"todo/db"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Username    string             `bson:"username,omitempty"`
	Password    string             `bson:"password,omitempty"`
	CreatedTS   time.Time          `bson:"created_ts,omitempty"`
	LastLoginTS time.Time          `bson:"last_login_ts,omitempty"`
	IsAdmin     bool               `bson:"is_admin,omitempty"`
}

func CreateUser(user *User) (*mongo.InsertOneResult, error) {
	user.CreatedTS = time.Now()
	result, err := db.UsersCollection.InsertOne(db.MongoContext, user)
	return result, err
}

func DeleteUser(user *User) (*mongo.DeleteResult, error) {
	res, err := db.UsersCollection.DeleteOne(db.MongoContext, bson.M{"_id": user.ID})
	return res, err
}

func FindUserById(id string) (User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	result := db.UsersCollection.FindOne(context.Background(), bson.M{"_id": objectId})
	resultUser := User{}
	err = result.Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return resultUser, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultUser, nil
}

func GetUsers() ([]User, error) {
	var results []User
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := db.UsersCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("Finding all users ERROR:", err)
		defer cursor.Close(ctx)
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		defer cursor.Close(ctx)
		return results, err
	}

	return results, nil
}
