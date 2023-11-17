package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akperrine/quik-coach/internal/db"
	"github.com/akperrine/quik-coach/internal/models"
	"github.com/akperrine/quik-coach/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	UserService models.UserService 
}

func NewUserController(collection *mongo.Collection) *UserController {

	// Initialize repository with the MongoDB collection
	userRepository := db.NewUserRepository(collection)

	// Initialize service with the repository
	userService := services.NewUserService(userRepository)

	// Create UserController with the initialized service
	return &UserController{
		UserService: userService,
	}
}


func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.UserService.FindAll()
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, users)
	}

	writeJSONResponse(w, http.StatusOK, users)
}

func (c *UserController) registerUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	_, err := c.UserService.CreateUser(*user)

	if err != nil {
		fmt.Printf("Error creating new user: %s", err)
		http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("user successfully created")
}

func (c *UserController) loginUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) 

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	response := c.UserService.FindOne(user.Email, user.Password)

	json.NewEncoder(w).Encode(response)
}


