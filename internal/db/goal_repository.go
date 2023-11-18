package db

import (
	"context"
	"log"
	"time"

	domain "github.com/akperrine/quik-coach/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type goalRepository struct {
	goalsCollection *mongo.Collection
	workoutCollection *mongo.Collection
}

func NewGoalRepository(goalsCollection , workoutCollection *mongo.Collection) domain.GoalRepository {
	return &goalRepository{
		goalsCollection: goalsCollection,
		workoutCollection: workoutCollection,
	}
}

func (r *goalRepository) FindGoalsByEmail(email string) ([]domain.GoalDto, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "user_email",Value: email}}}}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "workouts"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "goal_id"},
			{Key: "as", Value: "workouts"},
		}},
	}
	
	pipeline := mongo.Pipeline{matchStage, lookupStage}


	cursor, err := r.goalsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var goals []domain.GoalDto
	for cursor.Next(ctx) {
		var goal domain.GoalDto
		if err := cursor.Decode(&goal); err != nil {
			log.Println("Error decoding goal:", err)
			continue
		}
		log.Println(goal)
		var totDistance float64
		for _, workout := range goal.Workouts {
			log.Println("wod ",workout.Distance)
			totDistance += float64(workout.Distance)
		}
		goal.CurrentDistance = totDistance
		goals = append(goals, goal)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return goals, nil
}

func (r *goalRepository) Create(goal domain.Goal) (*mongo.InsertOneResult, error) {

	return r.goalsCollection.InsertOne(context.TODO(), goal)
}

func (*goalRepository) Update(id string) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}


func (*goalRepository) Delete(id string) (*mongo.DeleteResult, error) {
	panic("unimplemented")
}



