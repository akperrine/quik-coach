package db

import (
	"context"
	"fmt"
	"log"

	domain "github.com/akperrine/quik-coach/internal"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	users := []domain.User{}

	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) FindOneByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(user domain.User) (*mongo.InsertOneResult, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.ID = uuid.NewString()
	user.Password = string(encryptedPass)

	createdUser, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return createdUser, nil
}