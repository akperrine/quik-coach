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
	goalsCollection *mongo.Collection
	workoutCollection *mongo.Collection
}

var goalsCollection *mongo.Collection
var workoutCollection *mongo.Collection


func NewGoalsController(db *mongo.Database) *GoalsController {
	goalsCollection = db.Collection("goals")
	workoutCollection = db.Collection("workouts")
	return &GoalsController{
		goalsCollection: goalsCollection,
		workoutCollection: workoutCollection,
	}
}

func (c *GoalsController) GetAllUserGoals(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connected to collection: %s", c.goalsCollection.Name())
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


	cursor, err := c.goalsCollection.Find(ctx, filter)
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
	log.Printf("Connected to collection: %s\n", c.goalsCollection.Name())
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	// Empty filter to retrieve all goals
	filter := bson.M{}

	findOptions := options.Find()
	log.Println("Fetching all goals...", filter, "\n", findOptions)

	cursor, err := c.goalsCollection.Find(ctx, filter)
	if err != nil {
		// Handle the error, e.g., log it and return an error response
		http.Error(w, "Error fetching all goals", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var goals []models.Goal
	for cursor.Next(ctx) {
		log.Println(cursor)
		var goal models.Goal
		if err := cursor.Decode(&goal); err != nil {
			// Handle the decoding error, e.g., log it and skip the current document
			log.Println("cant parse")
			continue
		}
		log.Println(goal.UserEmail)
		goals = append(goals, goal)
	}

	log.Println("ggs",goals)

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

func (c *GoalsController) GetAllWods(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connected to collection: %s\n", c.workoutCollection.Name())
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	// Empty filter to retrieve all goals
	filter := bson.M{}

	findOptions := options.Find()
	log.Println("Fetching all goals...", filter, "\n", findOptions)

	cursor, err := c.workoutCollection.Find(ctx, filter)
	if err != nil {
		// Handle the error, e.g., log it and return an error response
		http.Error(w, "Error fetching all goals", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var workouts []models.Workout
	for cursor.Next(ctx) {
		log.Println(cursor)
		var workout models.Workout
		if err := cursor.Decode(&workout); err != nil {
			// Handle the decoding error, e.g., log it and skip the current document
			log.Println("cant parse")
			continue
		}
		log.Println(workout.UserEmail)
		workouts = append(workouts, workout)
	}

	log.Println("ggs",workouts)

	if err := cursor.Err(); err != nil {
		// Handle the cursor error, e.g., log it and return an error response
		http.Error(w, "Error fetching all goals", http.StatusInternalServerError)
		return
	}

	// Marshal the goals to JSON and send the response
	responseJSON, err := json.Marshal(workouts)
	if err != nil {
		// Handle the JSON marshaling error, e.g., log it and return an error response
		http.Error(w, "Error encoding all goals", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, responseJSON)
}