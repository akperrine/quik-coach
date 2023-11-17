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

type Goal struct {
	ID 			 string `json:"id" bson:"_id,omitempty"`
	UserEmail	 string `json:"user_email" bson:"user_email"`
	Name		 string `json:"name"`
	TargetDistance	float64 `json:"target_distance"`
	StartDate		int `json:"start_date"`
	TargetDate		int	`json:"target_date"`
	Device	   	 string `json:"device,omitempty"`
	Workouts  []Workout `json:"workouts,omitempty"`
}

type GoalDto struct {
	ID 			 string `json:"id" bson:"_id,omitempty"`
	UserEmail	 string `json:"user_email" bson:"user_email"`
	Name		 string `json:"name"`
	TargetDistance	float64 `json:"target_distance"`
	CurrentDistance float64 `json:"current_distance"`
	StartDate		int `json:"start_date"`
	TargetDate		int	`json:"target_date"`
	Device	   	 string `json:"device,omitempty"`
	Workouts  []Workout `json:"workouts,omitempty"`
}

type GoalRepository interface {
	FindAll() ([]Goal, error)
	FindOne(id string) (*Goal, error)
	Create(goal Goal) (*mongo.InsertOneResult, error)
	Update(id string) (*mongo.UpdateResult, error)
	Delete(id string) (*mongo.DeleteResult, error)
 }

type Workout struct {
	ID       	string 		`json:"id" bson:"_id,omitempty"`
	GoalID   	string 		`json:"goal_id" bson:"goal_id,omitempty"`
	UserEmail 	string 		`json:"user_email" bson:"user_email"`
	Distance 	float64    	`json:"distance"`
	Date     	int    		`json:"date"`
	Modality	string 		`json:"modality"`
}

type Modality string

const (
	Run     Modality = "Run"
	Walk    Modality = "Walk"
	Bike    Modality = "Bike"
	Row     Modality = "Row"
	Climb   Modality = "Climb"
	Swim    Modality = "Swim"
	Other   Modality = "Other"
)