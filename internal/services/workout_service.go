package services

import (
	"encoding/json"
	"fmt"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/google/uuid"
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

func (s *workoutService) CreateWorkout(workout domain.Workout) (*mongo.InsertOneResult, error) {
	if _, ok := domain.ModalitySet[workout.Modality]; !ok {
		return nil, fmt.Errorf("invalid modality")
	}

	workout.ID = uuid.NewString()

	result, err := s.workoutRepository.Create(workout)

	if err != nil {
		fmt.Println(err)
		return nil, err
	} 
	return result, nil
}

func (s *workoutService) UpdateWorkout(workout domain.Workout) (*mongo.UpdateResult, error) {
	if _, ok := domain.ModalitySet[workout.Modality]; !ok {
		return nil, fmt.Errorf("invalid modality")
	}

	result, err := s.workoutRepository.Update(workout)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}

func (s *workoutService) DeleteWorkout(workout domain.Workout) (*mongo.DeleteResult, error) {
	return s.workoutRepository.Delete(workout.ID)
}



