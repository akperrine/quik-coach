package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akperrine/quik-coach/internal/db"
)

func HandleRequests() {
	db := db.Connect()

	var userCollection = db.Collection("users")
	var goalsCollection = db.Collection("goals")
	var workoutsCollection = db.Collection("workouts")


	userController := NewUserController(userCollection)
	goalsController := NewGoalsController(goalsCollection, workoutsCollection)
	workoutsController := NewWorkoutsController(workoutsCollection)

	http.HandleFunc("/health_check", healthCheck)

	http.HandleFunc("/users", userController.GetAllUsers)
	http.HandleFunc("/users/register", userController.registerUser)
	http.HandleFunc("/users/login", userController.loginUser)

	http.HandleFunc("/goals/user/", goalsController.GetAllUserGoals)
	http.HandleFunc("/goals/create", goalsController.AddGoal)
	http.HandleFunc("/goals/update", goalsController.UpdateGoal)
	http.HandleFunc("/goals/delete", goalsController.DeleteGoal)

	http.HandleFunc("/workouts/user/", workoutsController.GetUserWorkouts)
	http.HandleFunc("/workouts/goal/", workoutsController.GetGoalWorkouts)
	http.HandleFunc("/workouts/create", workoutsController.Addworkout)
	http.HandleFunc("/workouts/update", workoutsController.Updateworkout)
	http.HandleFunc("/workouts/delete", workoutsController.Deleteworkout)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}