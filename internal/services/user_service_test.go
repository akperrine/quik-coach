package services

import (
	"testing"

	"github.com/akperrine/quik-coach/internal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserRepository struct {
	// You can add fields or methods specific to your test needs.
}

func (m *MockUserRepository) FindAll() ([]domain.User, error) {
	// Mock implementation for FindAll function.
	// Return sample data for testing.
	users := []domain.User{
		{ID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		{ID: "2", FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
	}
	return users, nil
}

func (m *MockUserRepository) FindOneByEmail(email string) (*domain.User, error) {
	// Return a sample user for testing.
	// Password is hashed password123 with 10 cost factor
	if email == "john@example.com" {
		user := &domain.User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "$2y$10$go8ZNp2w5TmjyiuoGRwJOOINZ9Qi1KWGZbAygqZc3/w84IxJSFcn6"}
		return user, nil
	} else {
		return nil, mongo.ErrNilDocument
	}
}

func (m *MockUserRepository) Create(user domain.User) (*mongo.InsertOneResult, error) {
	// Return a sample InsertOneResult for testing.
	result := &mongo.InsertOneResult{InsertedID: user.ID}
	return result, nil
}

func TestUserService_FindAll(t *testing.T) {
	// Create a UserService with the MockUserRepository.
	userService := NewUserService(&MockUserRepository{})

	// Call the FindAll function.
	result, err := userService.FindAll()

	// Assert that there is no error.
	assert.Nil(t, err)

	// Assert that the result is not nil.
	assert.NotNil(t, result)

}

func TestUserService_FindOne_EmailExists(t *testing.T) {

	userService := NewUserService(&MockUserRepository{})

	// Call the FindOne function with an existing email and a valid password.
	result := userService.FindOne("john@example.com", "password123")

	// Assert that the result is not nil.
	assert.NotNil(t, result)

	// Check if the 'user' key is present.
	_, userExists := result["user"]

	assert.True(t, userExists, "Expected 'user' key to be present")
	assert.Equal(t, "John", result["user"].(*domain.User).FirstName)
}

func TestUserService_FindOne_EmailNotExist(t *testing.T) {

	userService := NewUserService(&MockUserRepository{})

	// Call the FindOne function with an existing email and a valid password.
	result := userService.FindOne("j@example.com", "password123")

	// Assert that the result is not nil.
	assert.NotNil(t, result)

	assert.True(t, result["message"] == "Email address not found")
}

func TestUserService_FindOne_InvalidPassword(t *testing.T) {

	userService := NewUserService(&MockUserRepository{})

	// Call the FindOne function with an existing email and a valid password.
	result := userService.FindOne("john@example.com", "password")

	// Assert that the result is not nil.
	assert.NotNil(t, result)

	assert.True(t, result["message"] == "Invalid login credentials. Please try again")
}




func TestUserService_CreateUser(t *testing.T) {
	// Create a UserService with the MockUserRepository.
	userService := NewUserService(&MockUserRepository{})

	// Create a sample user for testing.
	user := domain.User{
		ID:        "3",
		FirstName: "Alice",
		LastName:  "Doe",
		Email:     "alice@example.com",
		Password:  "password456",
	}

	// Call the CreateUser function.
	result, err := userService.CreateUser(user)

	// Assert that there is no error.
	assert.Nil(t, err)

	// Assert that the result is not nil.
	assert.NotNil(t, result)

}

