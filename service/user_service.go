package service

import (
	"regexp"
	"time"
	"todo/dao"
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
	userDao dao.UserDaoInterface
}

func UserService(userDao dao.UserDaoInterface) *userService {
	return &userService{userDao}
}

func (userService *userService) CreateUser(user *model.User) (*model.User, error) {
	user.CreatedTS = time.Now()
	return userService.userDao.CreateUser(user)
}

func (userService *userService) DeleteUser(user *model.User) error {
	return userService.userDao.DeleteUser(user)
}

func (userService *userService) FindUserById(id string) (model.User, error) {
	return userService.userDao.FindUserById(id)
}

func (userService *userService) FindUserByUsername(username string) (model.User, error) {
	return userService.userDao.FindUserByUsername(username)
}

func (userService *userService) GetUsers() ([]model.User, error) {
	return userService.userDao.GetUsers()
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
