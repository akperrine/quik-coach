package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type goalRepository struct {
	collection *mongo.Collection
}

// func NewGoalRepository(collection *mongo.Collection) models.GoalRepository {
// 	return &goalRepository{
// 		collection: collection,
// 	}
// }