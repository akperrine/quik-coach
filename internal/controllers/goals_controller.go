package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/akperrine/quik-coach/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoalsController struct {
	collection *mongo.Collection
}

var goalsCollection *mongo.Collection

func GoalsCollection(db *mongo.Database) *mongo.Collection {
	goalsCollection = db.Collection("goals")
	log.Printf("Connected to database: %s, collection: %s\n", db.Name(), collection.Name())

	return goalsCollection
}

func NewGoalsController(db *mongo.Database) *GoalsController {
	collection := GoalsCollection(db)

	return &GoalsController{
		collection: collection,
	}
}

func (c *GoalsController) GetAllUserGoals(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "Incorrect path", http.StatusBadRequest)
		return
	}
	userEmail := parts[2]
	
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	
	filter := bson.M{"user_email": userEmail} 
	
	findOptions := options.Find()
	log.Println(parts, userEmail, filter, "\n", findOptions)


	cursor, err := goalsCollection.Find(ctx, filter)
	if err != nil {
		// Handle the error, e.g., log it and return an error response
		http.Error(w, "Error fetching user goals", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	log.Println(ctx)
	var goals []models.Goal
	for cursor.Next(ctx) {
		var goal models.Goal
		if err := cursor.Decode(&goal); err != nil {
			// Handle the decoding error, e.g., log it and skip the current document
			continue
		}
		goals = append(goals, goal)
	}

	if err := cursor.Err(); err != nil {
		// Handle the cursor error, e.g., log it and return an error response
		http.Error(w, "Error fetching user goals", http.StatusInternalServerError)
		return
	}

	// Marshal the goals to JSON and send the response
	responseJSON, err := json.Marshal(goals)
	if err != nil {
		// Handle the JSON marshaling error, e.g., log it and return an error response
		http.Error(w, "Error encoding user goals", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, responseJSON)

}

func (c *GoalsController) GetAllGoals(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	// Empty filter to retrieve all goals
	filter := bson.M{}

	findOptions := options.Find()
	log.Println("Fetching all goals...", filter, "\n", findOptions)

	cursor, err := goalsCollection.Find(ctx, filter)
	if err != nil {
		// Handle the error, e.g., log it and return an error response
		http.Error(w, "Error fetching all goals", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var goals []models.Goal
	for cursor.Next(ctx) {
		var goal models.Goal
		if err := cursor.Decode(&goal); err != nil {
			// Handle the decoding error, e.g., log it and skip the current document
			continue
		}
		goals = append(goals, goal)
	}

	if err := cursor.Err(); err != nil {
		// Handle the cursor error, e.g., log it and return an error response
		http.Error(w, "Error fetching all goals", http.StatusInternalServerError)
		return
	}

	// Marshal the goals to JSON and send the response
	responseJSON, err := json.Marshal(goals)
	if err != nil {
		// Handle the JSON marshaling error, e.g., log it and return an error response
		http.Error(w, "Error encoding all goals", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, responseJSON)
}
