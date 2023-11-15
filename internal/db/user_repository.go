package db

import (
	"context"
	"fmt"
	"log"

	"github.com/akperrine/quik-coach/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}


type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) models.UserRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	users := []models.User{}

	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Failed to decode user with error: %s", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error during cursor iteration with error: %s", err)
		return nil, err
	}

	return users, nil
}

func (r *userRepository) FindOneByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(user models.User) (*mongo.InsertOneResult, error) {
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