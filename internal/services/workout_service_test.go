package services

import (
	"encoding/json"
	"testing"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockWorkoutsRepository struct {
	mock.Mock
}

func (m *MockWorkoutsRepository) FindByEmail(email string) ([]domain.Workout, error) {
	args := m.Called(email)
	return args.Get(0).([]domain.Workout), args.Error(1)
}

func (m *MockWorkoutsRepository) FindByGoal(email string) ([]domain.Workout, error) {
	args := m.Called(email)
	return args.Get(0).([]domain.Workout), args.Error(1)
}

func (m *MockWorkoutsRepository) Create(worokout domain.Workout) (*mongo.InsertOneResult, error) {
	args := m.Called(worokout)
	if result, ok := args.Get(0).(*mongo.InsertOneResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

// Delete implements domain.WorkoutRepository.
func (m *MockWorkoutsRepository) Delete(id string) (*mongo.DeleteResult, error) {
	args := m.Called(id)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// Update implements domain.WorkoutRepository.
func (m *MockWorkoutsRepository) Update(workout domain.Workout) (*mongo.UpdateResult, error) {
	args := m.Called(workout)
	if result, ok := args.Get(0).(*mongo.UpdateResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func TestFindUserWorkouts(t *testing.T) {
	mockRepository := &MockWorkoutsRepository{}
	service := NewWorkoutService(mockRepository)

	expectedWorkouts := []domain.Workout{
		{ID:        "1",
		GoalID:    "goal123",
		UserEmail: "user@example.com",
		Distance:  10.5,
		Date:      1678000000,
		Modality:  "run",},
	}

	mockRepository.On("FindByEmail", "user@example.com").Return(expectedWorkouts,nil)

	var resultWorkouts []domain.Workout
	resultJSON, err := service.FindUserWorkouts("user@example.com")
	json.Unmarshal(resultJSON, &resultWorkouts)

	assert.NoError(t, err)
	assert.Equal(t,expectedWorkouts, resultWorkouts)
	mockRepository.AssertExpectations(t)

}

func TestFindGoalWorkouts(t *testing.T) {
	mockRepository := &MockWorkoutsRepository{}
	service := NewWorkoutService(mockRepository)

	expectedWorkouts := []domain.Workout{
		{ID:        "1",
		GoalID:    "goal123",
		UserEmail: "user@example.com",
		Distance:  10.5,
		Date:      1678000000,
		Modality:  "run",},
	}

	mockRepository.On("FindByGoal", "1").Return(expectedWorkouts,nil)

	var resultWorkouts []domain.Workout
	resultJSON, err := service.FindGoalWorkouts("1")
	json.Unmarshal(resultJSON, &resultWorkouts)

	assert.NoError(t, err)
	assert.Equal(t,expectedWorkouts, resultWorkouts)
	mockRepository.AssertExpectations(t)
}
