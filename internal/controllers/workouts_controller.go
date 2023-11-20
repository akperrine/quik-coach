package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/akperrine/quik-coach/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutsController struct {
	workoutCollection *mongo.Collection
}


func NewWorkoutsController(collection *mongo.Collection) *WorkoutsController {
	return &WorkoutsController{
		workoutCollection: collection,
	}
}

func (c *WorkoutsController) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connected to collection: %s\n", c.workoutCollection.Name())
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	// Empty filter to retrieve all goals
	filter := bson.M{}

	// findOptions := options.Find()
	// log.Println("Fetching all workouts...", filter, "\n")

	cursor, err := c.workoutCollection.Find(ctx, filter)
	if err != nil {
		// Handle the error, e.g., log it and return an error response
		http.Error(w, "Error fetching all goals", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	log.Println("Number of documents:", cursor.RemainingBatchLength())


	var workouts []domain.Workout
	for cursor.Next(ctx) {
		log.Println(cursor)
		var workout domain.Workout
		if err := cursor.Decode(&workout); err != nil {
			// Handle the decoding error, e.g., log it and skip the current document
			log.Println("cant parse")
			continue
		}
		log.Println(workout)
		workouts = append(workouts, workout)
	}


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

func (c *WorkoutsController) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	log.Println("adding wod User...")
}

func (c *WorkoutsController) Addworkout(w http.ResponseWriter, r *http.Request) {
	log.Println("adding wod...")
}

func (c *WorkoutsController) Updateworkout(w http.ResponseWriter, r *http.Request) {
	log.Println("updating wod...")

}

func (c *WorkoutsController) Deleteworkout(w http.ResponseWriter, r *http.Request) {
	log.Println("deleting wod...")
}