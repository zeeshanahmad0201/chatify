package models

import "github.com/dgrijalva/jwt-go"

type SignedDetails struct {
	Name   string
	Email  string
	UserId string
	jwt.StandardClaims
}
