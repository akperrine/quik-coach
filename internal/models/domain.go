package models

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)




type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
 }

 type UserService interface {
	FindAll() ([]byte, error)
	CreateUser(user User) (*mongo.InsertOneResult, error)
	FindOne(email, password string) map[string]interface{}
}

type Token struct {
	UserID string
	Name string
	Email string
	*jwt.RegisteredClaims
}

