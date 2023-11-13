package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/akperrine/quik-coach/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
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


func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}

func HandleRequests() {
	http.HandleFunc("/", GetAllUsers)
	http.HandleFunc("/health_check", healthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseJSON, _ := json.Marshal(map[string]string{"error": message})
	w.Write(responseJSON)
}