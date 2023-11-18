package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoalsController struct {
	goalsCollection *mongo.Collection
	workoutCollection *mongo.Collection
}

func NewGoalsController(goalsCollection , workoutCollection *mongo.Collection) *GoalsController {
	return &GoalsController{
		goalsCollection: goalsCollection,
		workoutCollection: workoutCollection,
	}
}

func (c *GoalsController) GetAllUserGoals(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connected to collection: %s", c.goalsCollection.Name())
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Incorrect path", http.StatusBadRequest)
		return
	}
	userEmail := parts[3]
	
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "user_email",Value: userEmail}}}}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "workouts"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "goal_id"},
			{Key: "as", Value: "workouts"},
		}},
	}
	
	pipeline := mongo.Pipeline{matchStage, lookupStage}


	cursor, err := c.goalsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, "Error fetching goals with workouts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var goals []domain.GoalDto
	for cursor.Next(ctx) {
		var goal domain.GoalDto
		if err := cursor.Decode(&goal); err != nil {
			log.Println("Error decoding goal:", err)
			continue
		}
		log.Println(goal)
		var totDistance float64
		for _, workout := range goal.Workouts {
			log.Println("wod ",workout.Distance)
			totDistance += float64(workout.Distance)
		}
		goal.CurrentDistance = totDistance
		goals = append(goals, goal)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error fetching goals with workouts", http.StatusInternalServerError)
		return
	}

	// Marshal the goals to JSON and send the response
	responseJSON, err := json.Marshal(goals)
	if err != nil {
		http.Error(w, "Error encoding goals with workouts", http.StatusInternalServerError)
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

	var goals []domain.Goal
	for cursor.Next(ctx) {
		log.Println(cursor)
		var goal domain.Goal
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

func (c *GoalsController) AddGoal(w http.ResponseWriter, r *http.Request) {
	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)

	if _, ok := domain.ModalitySet[goal.Modality]; !ok {
		http.Error(w, "Invalid modality chosen", http.StatusBadRequest)
		return
	}

	goal.ID = uuid.NewString()
	
	createdGoal, err := c.goalsCollection.InsertOne(context.TODO(), goal)

	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
		return
	}

	// return createdGoal, nil

	json.NewEncoder(w).Encode(createdGoal)
}

func (c *GoalsController) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)
	log.Println(goal)
	if _, ok := domain.ModalitySet[goal.Modality]; !ok {
		http.Error(w, "Invalid modality chosen", http.StatusBadRequest)
		return
	}
	fmt.Println(reflect.TypeOf(goal.TargetDistance))
	
	updateData := bson.M{
		"$set": bson.M{
			"name":            goal.Name,
			"target_distance": float64(goal.TargetDistance),
			"start_date":      int(goal.StartDate),
			"target_date":     int(goal.TargetDate),
			"modality":        goal.Modality,
		},
	}
	log.Println(goal.ID)
	_, err := c.goalsCollection.UpdateByID(context.TODO(), goal.ID, updateData)

	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User updated succesfuly")
}

func (c *GoalsController) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)

	result, err := c.goalsCollection.DeleteOne(context.TODO(), bson.M{"_id": goal.ID})
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Error deleting user: %s", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}