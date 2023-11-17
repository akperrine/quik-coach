package services

import (
	"fmt"
	"time"

	"github.com/akperrine/quik-coach/internal"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user domain.User) (string, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Minute * 100000))



	claims := &domain.Token{
		UserID: user.ID,
		Name:	user.FirstName,
		Email:	user.Email,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	 })
	 
	 if err != nil {
		return err
	 }
	
	 if !token.Valid {
		return fmt.Errorf("invalid token")
	 }
	
	 return nil
}