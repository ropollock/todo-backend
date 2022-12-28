package routes

import (
	"fmt"
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func GetJWTSecret() string {
	return config.AppConfig.JWTSecretKey
}

func GetRefreshJWTSecret() string {
	return config.AppConfig.JWTRefreshSecretKey
}

func generateAccessToken(user *services.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
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

func JWTErrorChecker(err error, c echo.Context) error {
	return c.String(http.StatusForbidden, "")
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

func RegisterLoginRoutes(e *echo.Echo) {
	e.POST("/login", HandleLogin)
	e.POST("/logout", HandleLogout)
	fmt.Println("Registered authentication routes.")
}

func HandleLogin(ctx echo.Context) error {
	var req LoginRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := services.FindUserByUsername(req.Username)

	if err != nil {
		return ctx.String(http.StatusUnauthorized, "username or password is incorrect.")
	}

	if !services.CheckPasswordHash(req.Password, userResult.Password) {
		return ctx.String(http.StatusUnauthorized, "username or password is incorrect.")
	}

	token, refreshToken, tokenErr := GenerateTokensAndSetCookies(&userResult, ctx)

	if tokenErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "token is incorrect")
	}

	return ctx.JSON(http.StatusOK, LoginResponse{Token: token, RefreshToken: refreshToken})
}

func HandleLogout(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "logout")
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Get("user") == nil {
			return next(c)
		}

		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
			rc, err := c.Cookie(refreshTokenCookieName)

			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})

				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					var _, _, _ = GenerateTokensAndSetCookies(&services.User{
						Name: claims.Name,
					}, c)
				}
			}
		}

		return next(c)
	}
}
