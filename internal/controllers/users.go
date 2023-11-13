package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func HandleRequests() {
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/users/register", registerUser)
	http.HandleFunc("/health_check", healthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusInternalServerError, "Failed to retrieve users from MongoDB")
		return
	}
	defer cursor.Close(context.Background())
	
	
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Failed to decode user with error: %s", err)
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error during cursor iteration with error: %s", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println("Failed to marshal users to JSON")
		return
	}

	writeJSONResponse(w, http.StatusOK, usersJSON)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	encrptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	user.ID = uuid.NewString()
	user.Password = string(encrptedPass)

    fmt.Printf("{ID: %s, FirstName: %s, LastName: %s, Email: %s}", user.ID, user.FirstName, user.LastName, user.Email)


	createdUser, insErr := collection.InsertOne(context.TODO(), user)

	if insErr != nil {
		fmt.Printf("Error creating new user: %s", insErr)
	}

	json.NewEncoder(w).Encode(createdUser)
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}


