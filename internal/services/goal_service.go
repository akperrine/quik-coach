package services

import (
	"encoding/json"
	"fmt"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type goalService struct {
	goalRepository domain.GoalRepository
}

func NewGoalService(goalRepository domain.GoalRepository) domain.GoalService {
	return &goalService{
		goalRepository: goalRepository,
	}
}

func (s *goalService) FindUserGoals(email string) ([]byte, error) {
	goals, err := s.goalRepository.FindGoalsByEmail(email)
	if err != nil {
		return nil, err
	}

	responseJSON, err := json.Marshal(goals)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}


func (s *goalService) CreateGoal(goal domain.Goal) (*mongo.InsertOneResult, error) {
	if _, ok := domain.ModalitySet[goal.Modality]; !ok {
		// http.Error(w, "Invalid modality chosen", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid modality")
	}

	goal.ID = uuid.NewString()
	
	result, err := s.goalRepository.Create(goal)

	if err != nil {
		fmt.Println(err)
		// http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
		return nil, err
	} 
	return result, nil
}

func (s *goalService) UpdateGoal(goal domain.Goal) (*mongo.UpdateResult, error) {
	if _, ok := domain.ModalitySet[goal.Modality]; !ok {
		// http.Error(w, "Invalid modality chosen", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid modality")
	}

	result, err := s.goalRepository.Update(goal)

	if err != nil {
		fmt.Println(err)
		// http.Error(w, fmt.Sprintf("Error creating new user: %s", err), http.StatusInternalServerError)
		return nil, err
	}
	return result, err
}

func (*goalService) DeleteGoal(goal domain.Goal) (*mongo.DeleteResult, error) {
	panic("unimplemented")
}

