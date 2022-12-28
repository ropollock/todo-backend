package routes

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"todo/services"
	"todo/services/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func JWTErrorChecker(err error, c echo.Context) error {
	return c.String(http.StatusForbidden, "")
}

func RegisterLoginRoutes(e *echo.Echo) {
	e.POST("/login", HandleLogin)
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

	token, refreshToken, tokenErr := auth.GenerateTokensAndSetCookies(&userResult, ctx)

	if tokenErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "token is incorrect")
	}

	return ctx.JSON(http.StatusOK, LoginResponse{Token: token, RefreshToken: refreshToken})
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Get("user") == nil {
			return next(c)
		}

		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*auth.Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
			rc, err := c.Cookie(auth.GetRefreshTokenCookieName())

			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(auth.GetRefreshJWTSecret()), nil
				})

				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					var _, _, _ = auth.GenerateTokensAndSetCookies(&services.User{
						Name: claims.Name,
					}, c)
				}
			}
		}

		return next(c)
	}
}
