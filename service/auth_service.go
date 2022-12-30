package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"todo/config"
	"todo/model"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

type authService struct {
	userService UserServiceInterface
}

type AuthServiceInterface interface {
	GetJWTSecret() string
	GetRefreshJWTSecret() string
	GetAccessTokenCookieName() string
	GetRefreshTokenCookieName() string
	GenerateTokensAndSetCookies(user *model.User, c echo.Context) (string, string, error)
	GetCurrentUser(ctx echo.Context) (model.User, error)
}

func AuthService(userService UserServiceInterface) *authService {
	return &authService{userService}
}

func (srv *authService) GetJWTSecret() string {
	return config.AppConfig.JWTSecretKey
}

func (srv *authService) GetRefreshJWTSecret() string {
	return config.AppConfig.JWTRefreshSecretKey
}

func (srv *authService) GetAccessTokenCookieName() string {
	return accessTokenCookieName
}

func (srv *authService) GetRefreshTokenCookieName() string {
	return refreshTokenCookieName
}

func (srv *authService) generateAccessToken(user *model.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return srv.generateToken(user, expirationTime, []byte(srv.GetJWTSecret()))
}

func (srv *authService) generateToken(user *model.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &model.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func (srv *authService) GenerateTokensAndSetCookies(user *model.User, c echo.Context) (string, string, error) {
	accessToken, exp, err := srv.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	srv.setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	srv.setUserCookie(user, exp, c)

	refreshToken, exp, err := srv.generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}
	srv.setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return accessToken, refreshToken, nil
}

func (srv *authService) generateRefreshToken(user *model.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	return srv.generateToken(user, expirationTime, []byte(srv.GetRefreshJWTSecret()))
}

func (srv *authService) setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func (srv *authService) setUserCookie(user *model.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func (srv *authService) GetCurrentUser(ctx echo.Context) (model.User, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.Claims)

	userResult, err := srv.userService.FindUserByUsername(claims.Username)
	if err != nil {
		return model.User{}, err
	}

	return userResult, nil
}
