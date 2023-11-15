package main

import (
	"log"

	"github.com/akperrine/quik-coach/internal/controllers"
	// "github.com/akperrine/quik-coach/internal/db"
)



func main() {
	// var db = db.Connect()

	controllers.HandleRequests()



	log.Fatal("Something went wrong...")

}