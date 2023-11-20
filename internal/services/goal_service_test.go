package services

import (
	"encoding/json"
	"testing"

	"github.com/akperrine/quik-coach/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockGoalRepository struct {
	mock.Mock
}

func (m *MockGoalRepository) FindGoalsByEmail(email string) ([]domain.GoalDto, error) {
	args := m.Called(email)
	return args.Get(0).([]domain.GoalDto), args.Error(1)
}

func (m *MockGoalRepository) Create(goal domain.Goal) (*mongo.InsertOneResult, error) {
	args := m.Called(goal)
	if result, ok := args.Get(0).(*mongo.InsertOneResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockGoalRepository) Update(goal domain.Goal) (*mongo.UpdateResult, error) {
	args := m.Called(goal)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockGoalRepository) Delete(goalID string) (*mongo.DeleteResult, error) {
	args := m.Called(goalID)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestFindUserGoals(t *testing.T) {
	mockRepository := &MockGoalRepository{}
	service := NewGoalService(mockRepository)

	// Mock the repository to return sample goals
	expectedGoals := []domain.GoalDto{
		{
			ID:              "123456789",
			UserEmail:       "john@example.com",
			Name:            "Run 5K",
			TargetDistance:  5.0,
			CurrentDistance: 2.5,
			StartDate:       1677840400,
			TargetDate:      1735689599,
			Modality:        "Run",
			Workouts:        nil,  // No workouts in this example
		},
	}
	mockRepository.On("FindGoalsByEmail", "john@example.com").Return(expectedGoals, nil)

	// Call the service method
	resultJSON, err := service.FindUserGoals("john@example.com")
	assert.NoError(t, err)

	// Unmarshal the JSON for assertion
	var resultGoals []domain.GoalDto
	err = json.Unmarshal(resultJSON, &resultGoals)
	assert.NoError(t, err)

	// Assert that the goals match the expected goals
	assert.Equal(t, expectedGoals, resultGoals)
	mockRepository.AssertExpectations(t)
}

func TestCreateGoal(t *testing.T) {
	mockRepository := &MockGoalRepository{}
	service := NewGoalService(mockRepository)

	// Mock the repository to return a sample InsertOneResult
	expectedResult := &mongo.InsertOneResult{
		InsertedID: "12345",
	}
	mockRepository.On("Create", mock.AnythingOfType("domain.Goal")).Return(expectedResult, nil)

	// Call the service method
	goal := domain.Goal{
		ID:              "123456789",
		UserEmail:       "john@example.com",
		Name:            "Run 5K",
		TargetDistance:  5.0,
		StartDate:       1677840400,
		TargetDate:      1735689599,
		Modality:        "run",
		Workouts:        nil,  // No workouts in this example
	}
	result, err := service.CreateGoal(goal)
	assert.NoError(t, err)

	// Assert that the result matches the expected result
	assert.NotNil(t, result.InsertedID)	// mockRepository.AssertExpectations(t)
	assert.Equal(t, expectedResult.InsertedID, result.InsertedID)
	
	mockRepository.AssertExpectations(t)
}

func TestUpdateGoal(t *testing.T) {
	mockRepository := &MockGoalRepository{}
	service := NewGoalService(mockRepository)

	// Mock the repository to return a sample UpdateResult
	expectedResult := &mongo.UpdateResult{}
	mockRepository.On("Update", mock.AnythingOfType("domain.Goal")).Return(expectedResult, nil)

	// Call the service method
	goal := domain.Goal{
		ID:              "123456789",
		UserEmail:       "john@example.com",
		Name:            "Run 5K",
		TargetDistance:  5.0,
		StartDate:       1677840400,
		TargetDate:      1735689599,
		Modality:        "run",
		Workouts:        nil,  // No workouts in this example
	}
	result, err := service.UpdateGoal(goal)
	assert.NoError(t, err)

	// Assert that the result matches the expected result
	assert.Equal(t, expectedResult, result)
	mockRepository.AssertExpectations(t)
}

func TestDeleteGoal(t *testing.T) {
	mockRepository := &MockGoalRepository{}
	service := NewGoalService(mockRepository)

	// Mock the repository to return a sample DeleteResult
	expectedResult := &mongo.DeleteResult{}
	mockRepository.On("Delete", "goalID").Return(expectedResult, nil)

	// Call the service method
	result, err := service.DeleteGoal(domain.Goal{ID: "goalID"})
	assert.NoError(t, err)

	// Assert that the result matches the expected result
	assert.Equal(t, expectedResult, result)
	mockRepository.AssertExpectations(t)
}