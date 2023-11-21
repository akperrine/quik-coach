package db

import (
	"context"
	"log"
	"time"

	domain "github.com/akperrine/quik-coach/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type workoutRepository struct {
	workoutCollection *mongo.Collection
}

func NewWorkoutRepository(workoutCollection *mongo.Collection) domain.WorkoutRepository {
	return &workoutRepository{
		workoutCollection: workoutCollection,
	}
}

func (r *workoutRepository) FindByGoal(id string) ([]domain.Workout, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	filter := bson.M{"goal_id": id}

	cursor, err := r.workoutCollection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

	var workouts []domain.Workout
    for cursor.Next(ctx) {
        var workout domain.Workout
        if err := cursor.Decode(&workout); err != nil {
            return nil, err
        }
        workouts = append(workouts, workout)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return workouts, nil
}

func (r *workoutRepository) FindByEmail(email string) ([]domain.Workout, error) {
	
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_email": email}

	cursor, err := r.workoutCollection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

	var workouts []domain.Workout
    for cursor.Next(ctx) {
        var workout domain.Workout
        if err := cursor.Decode(&workout); err != nil {
            return nil, err
        }
        workouts = append(workouts, workout)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return workouts, nil
}

func (r *workoutRepository) Create(worokout domain.Workout) (*mongo.InsertOneResult, error) {
	return r.workoutCollection.InsertOne(context.TODO(), worokout)
}

func (r *workoutRepository) Update(workout domain.Workout) (*mongo.UpdateResult, error) {
	updateData := bson.M{
		"$set": bson.M{
			"distance": workout.Distance,
			"date": workout.Date,
			"modality": workout.Modality,
		},
	}

	return r.workoutCollection.UpdateByID(context.TODO(), workout.ID, updateData)
}

func (r *workoutRepository) Delete(id string) (*mongo.DeleteResult, error) {
	result, err := r.workoutCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}