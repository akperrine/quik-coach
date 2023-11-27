package controllers

import (
	"fmt"
	"log"
	// "log"
	"net/http"

	"github.com/akperrine/quik-coach/internal/db"
	// "github.com/gorilla/handlers"
	// "github.com/rs/cors"
)

func HandleRequests() {
	db := db.Connect()

	var userCollection = db.Collection("users")
	var goalsCollection = db.Collection("goals")
	var workoutsCollection = db.Collection("workouts")


	userController := NewUserController(userCollection)
	goalsController := NewGoalsController(goalsCollection, workoutsCollection)
	workoutsController := NewWorkoutsController(workoutsCollection)

	protectedRoute := http.NewServeMux()


	http.HandleFunc("/health_check", healthCheck)

	http.HandleFunc("/users", userController.GetAllUsers)
	http.HandleFunc("/users/register", userController.registerUser)
	http.HandleFunc("/users/login", userController.loginUser)
	// router.HandleFunc("/users/login", userController.loginUser)
	// router.HandleFunc("/goals/user/", goalsController.GetAllUserGoals)

	protectedRoute.HandleFunc("/goals/user/", goalsController.GetAllUserGoals)
	http.Handle("/goals/user/", AuthenticationMiddleware(CORSMiddleware(protectedRoute)))
	// http.HandleFunc("/goals/user/", goalsController.GetAllUserGoals)
	http.HandleFunc("/goals/create", goalsController.AddGoal)
	http.HandleFunc("/goals/update", goalsController.UpdateGoal)
	http.HandleFunc("/goals/delete", goalsController.DeleteGoal)

	http.HandleFunc("/workouts/user/", workoutsController.GetUserWorkouts)
	http.HandleFunc("/workouts/goal/", workoutsController.GetGoalWorkouts)
	http.HandleFunc("/workouts/create", workoutsController.Addworkout)
	http.HandleFunc("/workouts/update", workoutsController.Updateworkout)
	http.HandleFunc("/workouts/delete", workoutsController.Deleteworkout)

	// handler := cors.Default().Handler(http.DefaultServeMux)
	handler := CORSMiddleware(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8000", handler))
	// http.Handle("/", router)
	http.ListenAndServe(":8000",nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, *")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			requestedHeaders := r.Header.Get("Access-Control-Request-Headers")
			if requestedHeaders != "" {
				log.Println("req", requestedHeaders)
				w.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}