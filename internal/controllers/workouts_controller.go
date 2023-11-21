package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/akperrine/quik-coach/internal"
	"github.com/akperrine/quik-coach/internal/db"
	"github.com/akperrine/quik-coach/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutsController struct {
	WorkoutService domain.WorkoutService
}


func NewWorkoutsController(collection *mongo.Collection) *WorkoutsController {
	workoutRepository := db.NewWorkoutRepository(collection)
	workoutService := services.NewWorkoutService(workoutRepository)

	return &WorkoutsController{
		WorkoutService: workoutService,
	}
}

func (c *WorkoutsController) GetGoalWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Incorrect path", http.StatusBadRequest)
		return
	}

	goalId := parts[3]

	workoutsJSON, err := c.WorkoutService.FindGoalWorkouts(goalId)
	if err != nil {
		http.Error(w, "Error getting goal workouts", http.StatusBadRequest)
	} else if string(workoutsJSON) == "null" {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, workoutsJSON)
}

func (c *WorkoutsController) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Incorrect path", http.StatusBadRequest)
		return
	}

	goalEmail := parts[3]

	workoutsJSON, err := c.WorkoutService.FindUserWorkouts(goalEmail)
	if err != nil {
		http.Error(w, "Error getting goal workouts", http.StatusBadRequest)
	} else if string(workoutsJSON) == "null" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, workoutsJSON)
}

func (c *WorkoutsController) Addworkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workout := &domain.Workout{}
	json.NewDecoder(r.Body).Decode(workout)

	result, err := c.WorkoutService.CreateWorkout(*workout)
	log.Println("hIIII", result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding workout", http.StatusBadRequest)
	}

	insertedID, ok := result.InsertedID.(string)
	if !ok {
		log.Println(err)
		http.Error(w, "Error asserting InsertedID", http.StatusInternalServerError)
	}
	
	workout.ID = insertedID
	json.NewEncoder(w).Encode(workout)
}

func (c *WorkoutsController) Updateworkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workout := &domain.Workout{}
	json.NewDecoder(r.Body).Decode(workout)

	_, err := c.WorkoutService.UpdateWorkout(*workout)

	if err != nil {
		http.Error(w, "Error updating user", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode("Workout updated succesfuly")
}

func (c *WorkoutsController) Deleteworkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workout := &domain.Workout{}
	json.NewDecoder(r.Body).Decode(workout)

	result, err := c.WorkoutService.DeleteWorkout(*workout)
	
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %s", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}