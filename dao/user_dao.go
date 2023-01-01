package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
	"todo/data"
	"todo/model"
)

type userDao struct {
	databaseProvider data.MongoDBProviderInterface
}

type UserDaoInterface interface {
	CreateUser(user *model.User) (*model.User, error)
	DeleteUser(user *model.User) error
	FindUserById(id string) (model.User, error)
	FindUserByUsername(username string) (model.User, error)
	GetUsers() ([]model.User, error)
}

func UserDao(databaseProvider data.MongoDBProviderInterface) *userDao {
	return &userDao{databaseProvider}
}

func (dao *userDao) CreateUser(user *model.User) (*model.User, error) {
	insertResult, err := dao.databaseProvider.GetUsersCollection().InsertOne(dao.databaseProvider.GetContext(), user)
	result, _ := dao.FindUserById(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return &result, err
}

func (dao *userDao) DeleteUser(user *model.User) error {
	_, err := dao.databaseProvider.GetUsersCollection().DeleteOne(dao.databaseProvider.GetContext(), bson.M{"_id": user.ID})
	return err
}

func (dao *userDao) FindUserById(id string) (model.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	result := dao.databaseProvider.GetUsersCollection().FindOne(context.Background(), bson.M{"_id": objectId})
	resultUser := model.User{}
	err = result.Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return resultUser, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultUser, nil
}

func (dao *userDao) FindUserByUsername(username string) (model.User, error) {
	result := dao.databaseProvider.GetUsersCollection().FindOne(context.Background(), bson.M{"username": username})
	resultUser := model.User{}
	err := result.Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return resultUser, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultUser, nil
}

func (dao *userDao) GetUsers() ([]model.User, error) {
	var results []model.User
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := dao.databaseProvider.GetUsersCollection().Find(context.Background(), bson.D{})
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
