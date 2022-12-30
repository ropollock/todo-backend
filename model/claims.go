package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"name"`
	jwt.StandardClaims
}
