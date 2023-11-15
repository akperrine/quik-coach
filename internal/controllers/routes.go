package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akperrine/quik-coach/internal/db"
)

func HandleRequests() {
	db := db.Connect()
	userController := NewUserController(db)

	http.HandleFunc("/users", userController.GetAllUsers)
	http.HandleFunc("/users/register", userController.registerUser)
	http.HandleFunc("/users/login", userController.loginUser)
	http.HandleFunc("/health_check", healthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}