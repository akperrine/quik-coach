package controllers

import (
	"encoding/json"
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
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Incorrect path", http.StatusBadRequest)
		return
	}
	userEmail := parts[3]
	
	goalJSON, err := c.GoalService.FindUserGoals(userEmail)
	if err != nil {
		http.Error(w, "Error getting user goals", http.StatusBadRequest)
	}

	writeJSONResponse(w, http.StatusOK, goalJSON)

}

func (c *GoalsController) AddGoal(w http.ResponseWriter, r *http.Request) {
	goal := &domain.Goal{}
	json.NewDecoder(r.Body).Decode(goal)

	response, err := c.GoalService.CreateGoal(*goal)

	if err != nil {
		http.Error(w, "Error adding goal", http.StatusBadRequest)
	}

	// if _, ok := domain.ModalitySet[goal.Modality]; !ok {
	// 	http.Error(w, "Invalid modality chosen", http.StatusBadRequest)
	// 	return
	// }

	// goal.ID = uuid.NewString()
	
	// createdGoal, err := c.goalsCollection.InsertOne(context.TODO(), goal)

	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
	// 	return
	// }

	json.NewEncoder(w).Encode(response)
}

// func (c *GoalsController) UpdateGoal(w http.ResponseWriter, r *http.Request) {
// 	goal := &domain.Goal{}
// 	json.NewDecoder(r.Body).Decode(goal)
// 	log.Println(goal)
// 	if _, ok := domain.ModalitySet[goal.Modality]; !ok {
// 		http.Error(w, "Invalid modality chosen", http.StatusBadRequest)
// 		return
// 	}
// 	fmt.Println(reflect.TypeOf(goal.TargetDistance))
	
// 	updateData := bson.M{
// 		"$set": bson.M{
// 			"name":            goal.Name,
// 			"target_distance": float64(goal.TargetDistance),
// 			"start_date":      int(goal.StartDate),
// 			"target_date":     int(goal.TargetDate),
// 			"modality":        goal.Modality,
// 		},
// 	}
// 	log.Println(goal.ID)
// 	_, err := c.goalsCollection.UpdateByID(context.TODO(), goal.ID, updateData)

// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode("User updated succesfuly")
// }

// func (c *GoalsController) DeleteGoal(w http.ResponseWriter, r *http.Request) {
// 	goal := &domain.Goal{}
// 	json.NewDecoder(r.Body).Decode(goal)

// 	result, err := c.goalsCollection.DeleteOne(context.TODO(), bson.M{"_id": goal.ID})
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, fmt.Sprintf("Error deleting user: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(result)
// }