package controllers

import (
	"encoding/json"
	// "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockGoalsService struct {
	mock.Mock
}

var mockGoalJSON  string = `{
	"id": "123456789",
	"user_email": "user@example.com",
	"name": "Run Marathon",
	"target_distance": 42.2,
	"start_date": 1677840400,
	"target_date": 1735689599,
	"modality": "Run",
	"workouts": [
	  {
		"id": "987654321",
		"date": 1678000000,
		"distance": 10.0,
		"duration": 3600
	  }
	]
  }`

func (m *MockGoalsService) FindUserGoals(userEmail string) ([]byte, error) {
	args := m.Called(userEmail)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockGoalsService) CreateGoal(goal domain.Goal) (*mongo.InsertOneResult, error) {
	args := m.Called(goal)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockGoalsService) UpdateGoal(goal domain.Goal) (*mongo.UpdateResult, error) {
	args := m.Called(goal)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockGoalsService) DeleteGoal(goal domain.Goal) (*mongo.DeleteResult, error) {
	args := m.Called(goal)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestGetUserGoals_UserExist(t *testing.T) {
	// Mock the GoalService
	mockGoalsService := &MockGoalsService{}  // Adjust this line based on your package structure
    controller := &GoalsController{
        GoalService: mockGoalsService,
    }

	// Test when user exists
	expectedResponse := mockGoalJSON
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedMap)
	if err != nil {
		t.Errorf("Error unmarshalling expectedResponse: %s", err)
	}
	
	
	mockGoalsService.On("FindUserGoals", "existing@example.com").Return([]byte(expectedResponse), nil)

	req := httptest.NewRequest(http.MethodGet, "/goals/user/existing@example.com", nil)
	res := httptest.NewRecorder()

	controller.GetAllUserGoals(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}

	actualResponse := res.Body.String()
	if assert.Equal(t, expectedResponse, actualResponse) {
		var actualMap map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &actualMap)
		assert.NoError(t, err) 
	}
	mockGoalsService.AssertExpectations(t)
}

func TestGetUserGoals_UserNotExist(t *testing.T) {
	// Mock the GoalService
	mockGoalsService := &MockGoalsService{}  // Adjust this line based on your package structure
    controller := &GoalsController{
        GoalService: mockGoalsService,
    }
	
	mockGoalsService.On("FindUserGoals", "nonexistent@example.com").Return([]byte("null"), nil)

	req := httptest.NewRequest(http.MethodGet, "/goals/user/nonexistent@example.com", nil)
	res := httptest.NewRecorder()

	controller.GetAllUserGoals(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, res.Code)
	}

	assert.Equal(t, res.Body.String(), "User not found\n")
}

// func TestCreateGoal(t *testing.T) {
//     // Create an instance of the mock
//     mock := new(MockGoalsService)

//     // Set up expected input and output for the test case
//     expectedGoal := domain.Goal{
// 		UserEmail:      "a@example.com",
// 		Name:           "Run Marathon",
// 		TargetDistance: 26.6,
// 		StartDate:      1672531199,
// 		TargetDate:     1735689599,
// 		Modality:       "bike",
//     }

//     expectedResult := map[string]interface{}{}

//     // Configure the mock to return the expected result when CreateGoal is called
//     mock.On("CreateGoal", expectedGoal).Return(expectedResult, nil)

//     // Call the method you are testing with the mock
//     result := mock.CreateGoal(expectedGoal)

//     // Assert that the result matches the expected result
//     assert.Equal(t, expectedResult, result)

//     mock.AssertExpectations(t)
// }