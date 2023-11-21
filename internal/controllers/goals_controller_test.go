package controllers

import (
	"bytes"
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
	if result, ok := args.Get(0).(*mongo.InsertOneResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockGoalsService) UpdateGoal(goal domain.Goal) (*mongo.UpdateResult, error) {
	args := m.Called(goal)
	if result, ok := args.Get(0).(*mongo.UpdateResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockGoalsService) DeleteGoal(goal domain.Goal) (*mongo.DeleteResult, error) {
	args := m.Called(goal)
	if result, ok := args.Get(0).(*mongo.DeleteResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
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

func TestCreateUserGoal(t *testing.T) {
	// Mock the GoalService
	mockGoalsService := &MockGoalsService{}  // Adjust this line based on your package structure
	controller := &GoalsController{
		GoalService: mockGoalsService,
	}

	// Test when creating a user goal
	expectedResponse := mockGoalJSON
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedMap)
	if err != nil {
		t.Errorf("Error unmarshalling expectedResponse: %s", err)
	}

	mockGoalsService.On("CreateGoal", mock.AnythingOfType("domain.Goal")).Return(`{"InsertedID":"beec7ec1-be03-47e7-820b-1da0bf8b1003"}`, nil)

	reqBody := []byte(`{
		"id": "abcdefg",
		"user_email": "user@example.com",
		"name": "Run Marathon",
		"target_distance": 42.2,
		"start_date": 1677840400,
		"target_date": 1735689599,
		"modality": "Run",
	}`)

	req := httptest.NewRequest(http.MethodPost, "/goals/create", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()

	controller.AddGoal(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}

	mockGoalsService.AssertExpectations(t)
}

func TestUpdateUserGoal(t *testing.T) {
	// Mock the GoalService
	mockGoalsService := &MockGoalsService{}  // Adjust this line based on your package structure
	controller := &GoalsController{
		GoalService: mockGoalsService,
	}

	// Test when updating a user goal
	expectedResponse := mockGoalJSON
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedMap)
	if err != nil {
		t.Errorf("Error unmarshalling expectedResponse: %s", err)
	}

	mockGoalsService.On("UpdateGoal", mock.AnythingOfType("domain.Goal")).Return(nil, nil)

	reqBody := []byte(`{
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
	}`)

	req := httptest.NewRequest(http.MethodPut, "/goals/update", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()

	controller.UpdateGoal(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}

	mockGoalsService.AssertExpectations(t)
}

func TestDeleteUserGoal(t *testing.T) {
	// Mock the GoalService
	mockGoalsService := &MockGoalsService{}  // Adjust this line based on your package structure
	controller := &GoalsController{
		GoalService: mockGoalsService,
	}

	// Test when deleting a user goal
	mockGoalsService.On("DeleteGoal", mock.AnythingOfType("domain.Goal")).Return(nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/goals/delete", nil)
	res := httptest.NewRecorder()

	controller.DeleteGoal(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}

	mockGoalsService.AssertExpectations(t)
}
