package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
	"time"
	"todo/db"
	"todo/model"
	"unicode"
)

const (
	USERNAME_REGEX_STRING = "^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$"
)

var (
	USERNAME_REGEX *regexp.Regexp
)

func init() {
	USERNAME_REGEX, _ = regexp.Compile(USERNAME_REGEX_STRING)
}

type UserServiceInterface interface {
	CreateUser(user *model.User) (*model.User, error)
	DeleteUser(user *model.User) error
	FindUserById(id string) (model.User, error)
	FindUserByUsername(username string) (model.User, error)
	GetUsers() ([]model.User, error)
	ValidatePassword(s string) bool
	ValidateUsername(s string) bool
	ScrubUserForAPI(u *model.User)
}

type userService struct {
}

func UserService() *userService {
	return &userService{}
}

func (userService *userService) CreateUser(user *model.User) (*model.User, error) {
	user.CreatedTS = time.Now()
	insertResult, err := db.UsersCollection.InsertOne(db.MongoContext, user)
	result, _ := userService.FindUserById(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return &result, err
}

func (userService *userService) DeleteUser(user *model.User) error {
	_, err := db.UsersCollection.DeleteOne(db.MongoContext, bson.M{"_id": user.ID})
	return err
}

func (userService *userService) FindUserById(id string) (model.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	result := db.UsersCollection.FindOne(context.Background(), bson.M{"_id": objectId})
	resultUser := model.User{}
	err = result.Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return resultUser, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultUser, nil
}

func (userService *userService) FindUserByUsername(username string) (model.User, error) {
	result := db.UsersCollection.FindOne(context.Background(), bson.M{"username": username})
	resultUser := model.User{}
	err := result.Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return resultUser, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultUser, nil
}

func (userService *userService) GetUsers() ([]model.User, error) {
	var results []model.User
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

func (userService *userService) ValidatePassword(s string) bool {
	if len(s) < 8 {
		return false
	}

	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}

func (userService *userService) ValidateUsername(s string) bool {
	var firstLetter = []rune(s)
	if (len(s) > 40 || len(s) < 4) || !USERNAME_REGEX.MatchString(s) || !unicode.IsLetter(firstLetter[0]) {
		return false
	}
	return true
}

func (userService *userService) ScrubUserForAPI(u *model.User) {
	u.Password = ""
}
