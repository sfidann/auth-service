package utils

import "github.com/dgrijalva/jwt-go"

type AccessToken struct {
	Userid string
	jwt.StandardClaims
}
