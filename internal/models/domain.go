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

 type UserRepository interface {
	FindAll() ([]User, error)
	FindOneByEmail(email string) (*User, error)
	Create(user User) (*mongo.InsertOneResult, error)
 }

 type UserService interface {
	FindAll() ([]byte, error)
	FindOne(email, password string) map[string]interface{}
	CreateUser(user User) (*mongo.InsertOneResult, error)
}

type Token struct {
	UserID string
	Name string
	Email string
	*jwt.RegisteredClaims
}

