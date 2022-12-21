package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MongoContext    context.Context
	MongoClient     *mongo.Client
	TodoDB          *mongo.Database
	UsersCollection *mongo.Collection
)

func Connect(dbURI string) {
	MongoContext = context.TODO()
	mongoconn := options.Client().ApplyURI(dbURI)
	MongoClient, err := mongo.Connect(MongoContext, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := MongoClient.Ping(MongoContext, readpref.Primary()); err != nil {
		panic(err)
	}

	TodoDB = MongoClient.Database("todo")
	UsersCollection = TodoDB.Collection("users")

	fmt.Println("MongoDB successfully connected.")
}

func StringToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
