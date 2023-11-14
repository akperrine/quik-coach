package services

import (
	"context"
	"fmt"
	"log"

	"github.com/akperrine/quik-coach/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}


func FindOne(email, password string) map[string]interface{}{
	fmt.Println(password)
	user := &models.User{}


	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(user)

	if err == mongo.ErrNoDocuments {
		// No matching user found
		fmt.Println("User not found")
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	} else if err != nil {
		log.Fatal(err)
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	} else {
		fmt.Printf("Found user: %+v\n", user)
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if passwordErr != nil {
        // Passwords don't match
        fmt.Println("Invalid password")
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
    }

	tokenString, tokenError := CreateToken(*user) 

	if tokenError != nil {
		var resp = map[string]interface{}{"status": false, "message": "Error creating token", "error": tokenError}
		return resp
	}

	user.Password = ""
	
	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}