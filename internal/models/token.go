package models

import jwt "github.com/golang-jwt/jwt/v5"

// import jwt

type Token struct {
	UserID string
	Name string
	Email string
	*jwt.RegisteredClaims
}




