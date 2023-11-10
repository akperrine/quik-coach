package main

import (
	"github.com/akperrine/quik-coach/internal/controllers"
	"github.com/akperrine/quik-coach/internal/db"
)



func main() {
	config.Connect()
	http.HandleRequests()
	
}