package domain

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

var ModalitySet = map[string]struct{}{
	"run":   {},
	"walk":  {},
	"bike":  {},
	"row":   {},
	"climb": {},
	"swim":  {},
	"other": {},
}

type Goal struct {
	ID 			 	string 		`json:"id" bson:"_id,omitempty"`
	UserEmail	 	string  	`json:"user_email" bson:"user_email"`
	Name		 	string  	`json:"name"`
	TargetDistance	float64 	`json:"target_distance" bson:"target_distance"`
	StartDate		int 		`json:"start_date" bson:"start_date"`
	TargetDate		int			`json:"target_date" bson:"target_date"`
	Modality	   	string  	`json:"modality,omitempty"`
	Workouts  	 	[]Workout   `json:"workouts,omitempty"`
}

type GoalDto struct {
	ID 			 string `json:"id" bson:"_id,omitempty"`
	UserEmail	 string `json:"user_email" bson:"user_email"`
	Name		 string `json:"name"`
	TargetDistance	float64 `json:"target_distance" bson:"target_distance"`
	CurrentDistance float64 `json:"current_distance"`
	StartDate		int `json:"start_date" bson:"start_date"`
	TargetDate		int	`json:"target_date" bson:"target_date"`
	Modality	   	 string `json:"modality,omitempty"`
	Workouts  []Workout `json:"workouts,omitempty"`
}

type GoalService interface {
	FindUserGoals(email string) ([]byte, error)
	CreateGoal(goal Goal) (*mongo.InsertOneResult, error)
	UpdateGoal(goal Goal) (*mongo.UpdateResult, error)
	DeleteGoal(goal Goal) (*mongo.DeleteResult, error)
}

type GoalRepository interface {
	FindGoalsByEmail(email string) ([]GoalDto, error)
	Create(goal Goal) (*mongo.InsertOneResult, error)
	Update(goal Goal) (*mongo.UpdateResult, error)
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

type WorkoutService interface {
	FindUserWorkouts(email string) ([]byte, error)
	FindGoalWorkouts(id string) ([]byte, error)
	CreateWorkout(Wwrkout Workout) (*mongo.InsertOneResult, error)
	UpdateWorkout(workout Workout) (*mongo.UpdateResult, error)
	DeleteWorkout(workout Workout) (*mongo.DeleteResult, error)
}

type WorkoutRepository interface {
	FindByEmail(email string) ([]Workout, error)
	FindByGoal(id string) ([]Workout, error)
	Create(worokout Workout) (*mongo.InsertOneResult, error)
	Update(workout Workout) (*mongo.UpdateResult, error)
	Delete(id string) (*mongo.DeleteResult, error)
}
