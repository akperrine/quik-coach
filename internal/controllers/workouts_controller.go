package controllers

import (
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
}

func (c *WorkoutsController) Updateworkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (c *WorkoutsController) Deleteworkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}