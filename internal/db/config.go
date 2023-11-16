package db

import (
	"context"
	"fmt"
	"time"

	// "github.com/akperrine/quik-coach/internal/controllers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Database{
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:letsgetrusty@cluster0.p80xcty.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverApi)

	ctx, cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database("o2_shark").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); 
	err != nil {
		panic(err)
	  }

	db := client.Database("o2_shark")  
	
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return db
}