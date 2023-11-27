package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/akperrine/quik-coach/internal"
	"github.com/akperrine/quik-coach/internal/db"
	"github.com/akperrine/quik-coach/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	UserService domain.UserService 
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := c.UserService.FindAll()
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, users)
	}

	writeJSONResponse(w, http.StatusOK, users)
}

func (c *UserController) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &domain.User{}
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &domain.User{}
	err := json.NewDecoder(r.Body).Decode(user) 

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	response := c.UserService.FindOne(user.Email, user.Password)
	log.Println("maptime ",response["user"])
    userValue := reflect.TypeOf(response["user"])
	log.Println("log ",userValue)

	userMap, ok := response["user"].(*domain.User)
	if ok {
		log.Println("good ", response["user"])
	}
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.ID = userMap.ID
	user.FirstName = userMap.FirstName
	user.Email = userMap.Email
	
	tokenErr := services.CreateToken(*user, w) 

	if tokenErr != nil {
		// Handle the error
		http.Error(w, "Internal Server Error creating token", http.StatusInternalServerError)
		return
	}
	// if tokenError != nil {
	// 	var response = map[string]interface{}{"status": false, "message": "Error creating token", "error": tokenError}
	// 	return http.Error(w, "response", http.StatusUnauthorized)
	// }
	// token := response["token"] 
	// http.SetCookie(w, &token)
	json.NewEncoder(w).Encode(response)
}


func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Cookies())
		cookie, err := r.Cookie("auth")
		log.Println("hi ", cookie, err)
		if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

		tokenString := cookie.Value
		if err := services.VerifyToken(tokenString); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
		}

		next.ServeHTTP(w, r)
	}) 
}