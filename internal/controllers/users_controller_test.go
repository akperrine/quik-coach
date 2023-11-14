package controllers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/akperrine/quik-coach/internal/models"
// )

// // MockUserService is a mock implementation of the services.UserService interface for testing purposes.
// type MockUserService struct{}

// func (m *MockUserService) FindAll() ([]models.User, error) {
// 	// Implement mock behavior for FindAll
// 	return []models.User{}, nil
// }

// func (m *MockUserService) CreateUser(user models.User) (string, error) {
// 	// Implement mock behavior for CreateUser
// 	return "user_id", nil
// }

// func (m *MockUserService) FindOne(email, password string) map[string]interface{} {
// 	// Implement mock behavior for FindOne
// 	return map[string]interface{}{"status": true, "message": "User found"}
// }

// func TestGetAllUsers(t *testing.T) {
// 	// Create a new instance of the controller with the mock service
// 	controller := &UserController{UserService: &MockUserService{}}

// 	// Create a request to simulate a GET request
// 	req, err := http.NewRequest("GET", "/users", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a ResponseRecorder to record the response
// 	rr := httptest.NewRecorder()
// 	fmt.Println(rr)

// 	// Call the handler function
// 	http.HandlerFunc(controller.GetAllUsers).ServeHTTP(rr, req)

// 	// Check the response status code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// TODO: Add more assertions for the response body or other details
// }

// func TestRegisterUser(t *testing.T) {
// 	// Create a new instance of the controller with the mock service
// 	controller := &UserController{UserService: &MockUserService{}}

// 	// Create a user for the test
// 	user := models.User{
// 		// TODO: Set user details for testing
// 	}

// 	// Convert user to JSON
// 	userJSON, err := json.Marshal(user)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a request to simulate a POST request with user data
// 	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a ResponseRecorder to record the response
// 	rr := httptest.NewRecorder()

// 	// Call the handler function
// 	http.HandlerFunc(controller.registerUser).ServeHTTP(rr, req)

// 	// Check the response status code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// TODO: Add more assertions for the response body or other details
// }

// func TestLoginUser(t *testing.T) {
// 	// Create a new instance of the controller with the mock service
// 	controller := &UserController{UserService: &MockUserService{}}

// 	// Create a user for the test
// 	user := models.User{
// 		// TODO: Set user details for testing
// 	}

// 	// Convert user to JSON
// 	userJSON, err := json.Marshal(user)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a request to simulate a POST request with user data
// 	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a ResponseRecorder to record the response
// 	rr := httptest.NewRecorder()

// 	// Call the handler function
// 	http.HandlerFunc(controller.loginUser).ServeHTTP(rr, req)

// 	// Check the response status code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// TODO: Add more assertions for the response body or other details
// }
