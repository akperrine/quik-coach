package controllers

import (
	"fmt"
	"log"
	"net/http"
)

func HandleRequests() {
	userController := NewUserController()

	http.HandleFunc("/users", userController.GetAllUsers)
	http.HandleFunc("/users/register", userController.registerUser)
	http.HandleFunc("/users/login", userController.loginUser)
	http.HandleFunc("/health_check", healthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}