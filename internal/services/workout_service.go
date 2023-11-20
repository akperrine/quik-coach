package services

import (
	"encoding/json"

	domain "github.com/akperrine/quik-coach/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type workoutService struct {
	workoutRepository domain.WorkoutRepository
}

func NewWorkoutService(workoutRepository domain.WorkoutRepository) domain.WorkoutService {
	return &workoutService{
		workoutRepository: workoutRepository,
	}
}


func (s *workoutService) FindGoalWorkouts(id string) ([]byte, error) {
	workouts, err := s.workoutRepository.FindByGoal(id)
	if err != nil {
		return nil, err
	}

	responseJSON, err := json.Marshal(workouts)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}

func (s *workoutService) FindUserWorkouts(email string) ([]byte, error) {
	workouts, err := s.workoutRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	responseJSON, err := json.Marshal(workouts)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}

func (*workoutService) CreateWorkout(Wwrkout domain.Workout) (*mongo.InsertOneResult, error) {
	panic("unimplemented")
}

// DeleteWorkout implements domain.WorkoutService.
func (*workoutService) DeleteWorkout(workout domain.Workout) (*mongo.DeleteResult, error) {
	panic("unimplemented")
}

// UpdateWorkout implements domain.WorkoutService.
func (*workoutService) UpdateWorkout(workout domain.Workout) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}

