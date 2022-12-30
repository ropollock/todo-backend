package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"todo/model"
	"todo/service"
)

type authController struct {
	userService service.UserServiceInterface
	authService service.AuthServiceInterface
}

func AuthController(userService service.UserServiceInterface, authService service.AuthServiceInterface) *authController {
	return &authController{userService, authService}
}

func (controller *authController) JWTErrorChecker(err error, c echo.Context) error {
	return c.String(http.StatusForbidden, "")
}

func (controller *authController) RegisterLoginRoutes(e *echo.Echo) {
	e.POST("/login", controller.HandleLogin)
	fmt.Println("Registered authentication routes.")
}

func (controller *authController) HandleLogin(ctx echo.Context) error {
	var req model.LoginRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	userResult, err := controller.userService.FindUserByUsername(req.Username)

	if err != nil {
		return ctx.String(http.StatusUnauthorized, "username or password is incorrect.")
	}

	if !checkPasswordHash(req.Password, userResult.Password) {
		return ctx.String(http.StatusUnauthorized, "username or password is incorrect.")
	}

	token, refreshToken, tokenErr := controller.authService.GenerateTokensAndSetCookies(&userResult, ctx)

	if tokenErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "token is incorrect")
	}

	return ctx.JSON(http.StatusOK, model.LoginResponse{Token: token, RefreshToken: refreshToken})
}

func (controller *authController) TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Get("user") == nil {
			return next(c)
		}

		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*model.Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
			rc, err := c.Cookie(controller.authService.GetRefreshTokenCookieName())

			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(controller.authService.GetRefreshJWTSecret()), nil
				})

				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					var _, _, _ = controller.authService.GenerateTokensAndSetCookies(&model.User{
						Username: claims.Username,
					}, c)
				}
			}
		}

		return next(c)
	}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
