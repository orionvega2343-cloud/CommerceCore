package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId string
	Email  string
	jwt.RegisteredClaims
}
