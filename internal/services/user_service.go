package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/akperrine/quik-coach/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


type userService struct{
	userRepoistory domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) domain.UserService {
	return &userService{
		userRepoistory: userRepository,
	}
}


func (s *userService) FindAll() ([]byte, error) {
	users, err := s.userRepoistory.FindAll()
	if  err != nil {
		fmt.Println(err)
		return nil, err
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		log.Println("Failed to marshal users to JSON")
		return nil, err
	}

	return jsonUsers, nil
}

func (s *userService) FindOne(email, password string) map[string]interface{}{
	fmt.Println(password)
	user := &domain.User{}

	user, err := s.userRepoistory.FindOneByEmail(email)
	log.Println("use err ", user, err, email)

	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		// No matching user found
		log.Println("User not found")
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	} else if err != nil {
		log.Fatal(err)
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	} else {
		log.Printf("Found user: %+v\n", &user)
	}
	log.Println(user.Password, password)
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	log.Println(passwordErr)
    if passwordErr != nil {
        // Passwords don't match
        fmt.Println("Invalid password")
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
    }

	// tokenString, tokenError := CreateToken(*user, w) 

	// if tokenError != nil {
	// 	var resp = map[string]interface{}{"status": false, "message": "Error creating token", "error": tokenError}
	// 	return resp
	// }

	user.Password = ""
	
	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	// resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}

func (s *userService) CreateUser(user domain.User) (*mongo.InsertOneResult, error) {
	email := user.Email
	existingUser, _ := s.userRepoistory.FindOneByEmail(email)

	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	
	return s.userRepoistory.Create(user)
}

