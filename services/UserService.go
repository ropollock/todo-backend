package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"regexp"
	"time"
	"todo/db"
	"unicode"
)

const (
	USERNAME_REGEX_STRING = "^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$"
)

var (
	USERNAME_REGEX *regexp.Regexp
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

func init() {
	USERNAME_REGEX, _ = regexp.Compile(USERNAME_REGEX_STRING)
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

func FindUserByUsername(username string) (User, error) {
	result := db.UsersCollection.FindOne(context.Background(), bson.M{"username": username})
	resultUser := User{}
	err := result.Decode(&resultUser)
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

func ValidatePassword(s string) bool {
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

func ValidateUsername(s string) bool {
	var firstLetter = []rune(s)
	if (len(s) > 40 || len(s) < 4) || !USERNAME_REGEX.MatchString(s) || !unicode.IsLetter(firstLetter[0]) {
		return false
	}
	return true
}

func ScrubUserForAPI(u *User) {
	u.Password = ""
}
