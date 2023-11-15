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

var collection *mongo.Collection

func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}

var userRepoistory = db.NewUserRepository(collection)
var userService = services.NewUserService(userRepoistory)

func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userService.FindAll()
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, users)
	}

	writeJSONResponse(w, http.StatusOK, users)
}

func (c *UserController) registerUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	_, err := userService.CreateUser(*user)


	if err != nil {
		fmt.Printf("Error creating new user: %s", err)
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

	response := userService.FindOne(user.Email, user.Password)

	json.NewEncoder(w).Encode(response)
}




func NewUserController() *UserController {
    return &UserController{
    }
}