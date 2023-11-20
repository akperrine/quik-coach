package controllers

import (
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"strings"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/akperrine/quik-coach/internal/db"
	"github.com/akperrine/quik-coach/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoalsController struct {
	GoalService domain.GoalService
}

func NewGoalsController(goalsCollection , workoutCollection *mongo.Collection) *GoalsController {
	goalRepository := db.NewGoalRepository(goalsCollection, workoutCollection)
	goalService := services.NewGoalService(goalRepository)

	return &GoalsController{
		GoalService: goalService,
	}
}



func (c *GoalsController) GetAllUserGoals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Incorrect path", http.StatusBadRequest)
		return
	}
	userEmail := parts[3]
	
	goalJSON, err := c.GoalService.FindUserGoals(userEmail)
	if err != nil {
		http.Error(w, "Error getting user goals", http.StatusBadRequest)
	} else if string(goalJSON) == "null" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// log.Println(s

	writeJSONResponse(w, http.StatusOK, goalJSON)

}

func (c *GoalsController) AddGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
		
	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)

	response, err := c.GoalService.CreateGoal(*goal)

	if err != nil {
		http.Error(w, "Error adding goal", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(response)
}

func (c *GoalsController) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)

	_, err := c.GoalService.UpdateGoal(*goal)

	if err != nil {
		http.Error(w, "Error updating user", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode("User updated succesfuly")
}

func (c *GoalsController) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)

	result, err := c.GoalService.DeleteGoal(*goal)
	
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %s", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

