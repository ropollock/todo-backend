package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"todo/config"
	"todo/services"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func GetJWTSecret() string {
	return config.AppConfig.JWTSecretKey
}

func GetRefreshJWTSecret() string {
	return config.AppConfig.JWTRefreshSecretKey
}

func GetAccessTokenCookieName() string {
	return accessTokenCookieName
}

func GetRefreshTokenCookieName() string {
	return refreshTokenCookieName
}

func generateAccessToken(user *services.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateToken(user *services.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &Claims{
		Name: user.Name,
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

func GenerateTokensAndSetCookies(user *services.User, c echo.Context) (string, string, error) {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)

	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return accessToken, refreshToken, nil
}

func generateRefreshToken(user *services.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(user *services.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}
