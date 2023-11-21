package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockWorkoutsService struct {
	mock.Mock
}

var mockWorkoutArray string = `{
	"id": "1",
	"workout_id": "007ffa15-5d8e-4ff6-95bb-821e83e20e34",
	"user_email": "a@example.com",
	"distance": 5,
	"date": 1772617599,
	"modality": "run"
}`

func (m *MockWorkoutsService) FindUserWorkouts(userEmail string) ([]byte, error) {
	args := m.Called(userEmail)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockWorkoutsService) FindGoalWorkouts(id string) ([]byte, error) {
	args := m.Called(id)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockWorkoutsService) CreateWorkout(Workout domain.Workout) (*mongo.InsertOneResult, error) {
	args := m.Called(Workout)
	if result, ok := args.Get(0).(*mongo.InsertOneResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockWorkoutsService) UpdateWorkout(Workout domain.Workout) (*mongo.UpdateResult, error) {
	args := m.Called(Workout)
	if result, ok := args.Get(0).(*mongo.UpdateResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockWorkoutsService) DeleteWorkout(Workout domain.Workout) (*mongo.DeleteResult, error) {
	args := m.Called(Workout)
	if result, ok := args.Get(0).(*mongo.DeleteResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func TestGetUserWorkouts(t *testing.T) {
	mockWorkoutsService := &MockWorkoutsService{}

	controller := &WorkoutsController{
		WorkoutService: mockWorkoutsService,
	}

	expectedResponse := mockWorkoutArray
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedMap)
	if err != nil {
		t.Errorf("Error unmarshalling expectedResponse: %s", err)
	}

	mockWorkoutsService.On("FindUserWorkouts", "existing@example.com").Return([]byte(expectedResponse), nil)

	req := httptest.NewRequest(http.MethodGet, "/workouts/user/existing@example.com", nil)
	res := httptest.NewRecorder()

	controller.GetUserWorkouts(res, req)

	actualResponse := res.Body.String()
	if assert.Equal(t, expectedResponse, actualResponse) {
		var actualMap map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &actualMap)
		assert.NoError(t, err) 
	}

	mockWorkoutsService.AssertExpectations(t)
}

func TestGetGoalWorkouts(t *testing.T) {
	mockWorkoutsService := &MockWorkoutsService{}

	controller := &WorkoutsController{
		WorkoutService: mockWorkoutsService,
	}

	expectedResponse := mockWorkoutArray
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedMap)
	if err != nil {
		t.Errorf("Error unmarshalling expectedResponse: %s", err)
	}

	mockWorkoutsService.On("FindGoalWorkouts", "1").Return([]byte(expectedResponse), nil)

	req := httptest.NewRequest(http.MethodGet, "/workouts/goal/1", nil)
	res := httptest.NewRecorder()

	controller.GetGoalWorkouts(res, req)

	actualResponse := res.Body.String()
	if assert.Equal(t, expectedResponse, actualResponse) {
		var actualMap map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &actualMap)
		assert.NoError(t, err) 
	}

	mockWorkoutsService.AssertExpectations(t)
}

func TestCreateWorkout(t *testing.T) {
	mockWorkoutsService := &MockWorkoutsService{}

	controller := &WorkoutsController{
		WorkoutService: mockWorkoutsService,
	}

	expectedResponse := `{
		"id": "1",
		"goal_id": "",
		"user_email": "a@example.com",
		"distance": 5,
		"date": 1772617599,
		"modality": "run"
	}`
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedMap)
	if err != nil {
		t.Errorf("Error unmarshalling expectedResponse: %s", err)
	}
	log.Println(err)
	
	expectedInsertedID := "1"
	mockWorkoutsService.On("CreateWorkout", mock.AnythingOfType("domain.Workout")).Return(&mongo.InsertOneResult{InsertedID: expectedInsertedID}, nil)


	reqBody := []byte(`{
		"id": "1",
		"user_email": "a@example.com",
		"distance": 5,
		"date": 1772617599,
		"modality": "run"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/workouts/create", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()

	controller.Addworkout(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}

	var actualMap map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &actualMap); err != nil {
		t.Errorf("Error decoding actualResponse %s", err)
	}

	fmt.Println("Expected Map:", expectedMap)
	fmt.Println("Actual Map:", actualMap)

	if !reflect.DeepEqual(actualMap, expectedMap) {
        t.Errorf("Actual map is not equal to expected map")
    }
	mockWorkoutsService.AssertExpectations(t)
}

func TestUpdateWorkout(t *testing.T) {
	mockWorkoutsService := &MockWorkoutsService{}

	controller := &WorkoutsController{
		WorkoutService: mockWorkoutsService,
	}

	mockWorkoutsService.On("UpdateWorkout", mock.AnythingOfType("domain.Workout")).Return(&mongo.UpdateResult{}, nil)

	reqBody := []byte(`{
		"id": "1",
		"user_email": "a@example.com",
		"distance": 5,
		"date": 1772617599,
		"modality": "run"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/workouts/update", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()

	controller.Updateworkout(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}
	actual := strings.TrimSpace(res.Body.String())
	assert.Equal(t, actual,  "\"Workout updated succesfuly\"")


	mockWorkoutsService.AssertExpectations(t)
}

func TestDeleteWorkout(t *testing.T) {
	mockWorkoutsService := &MockWorkoutsService{}

	controller := &WorkoutsController{
		WorkoutService: mockWorkoutsService,
	}

	mockWorkoutsService.On("DeleteWorkout", mock.AnythingOfType("domain.Workout")).Return(&mongo.DeleteResult{}, nil)

	reqBody := []byte(`{
		"id": "1",
		"user_email": "a@example.com",
		"distance": 5,
		"date": 1772617599,
		"modality": "run"
	}`)

	req := httptest.NewRequest(http.MethodDelete, "/workouts/delete", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()

	controller.Deleteworkout(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}
	actual := strings.TrimSpace(res.Body.String())
	assert.Equal(t, actual,  "{\"DeletedCount\":0}")


	mockWorkoutsService.AssertExpectations(t)
}