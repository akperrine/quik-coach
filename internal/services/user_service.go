package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/akperrine/quik-coach/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}

type userService struct{
	userRepoistory models.UserRepository
}

func NewUserService(userRepository models.UserRepository) models.UserService {
	return &userService{
		userRepoistory: userRepository,
	}
}

func (*userService) FindAll() ([]byte, error) {
	users := []models.User{}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	
	
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Failed to decode user with error: %s", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error during cursor iteration with error: %s", err)
		return nil, err
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		log.Println("Failed to marshal users to JSON")
		return nil, err
	}

	return jsonUsers, nil

}

func (*userService) FindOne(email, password string) map[string]interface{}{
	fmt.Println(password)
	user := &models.User{}


	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(user)

	if err == mongo.ErrNoDocuments {
		// No matching user found
		fmt.Println("User not found")
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	} else if err != nil {
		log.Fatal(err)
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	} else {
		fmt.Printf("Found user: %+v\n", user)
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if passwordErr != nil {
        // Passwords don't match
        fmt.Println("Invalid password")
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
    }

	tokenString, tokenError := CreateToken(*user) 

	if tokenError != nil {
		var resp = map[string]interface{}{"status": false, "message": "Error creating token", "error": tokenError}
		return resp
	}

	user.Password = ""
	
	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}

func (*userService) CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	encrptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.ID = uuid.NewString()
	user.Password = string(encrptedPass)

    fmt.Printf("{ID: %s, FirstName: %s, LastName: %s, Email: %s}", user.ID, user.FirstName, user.LastName, user.Email)


	createdUser, err := collection.InsertOne(context.TODO(), user)
	if createdUser != nil {
		fmt.Println(err)
		return nil, err
	}

	return createdUser, nil
}