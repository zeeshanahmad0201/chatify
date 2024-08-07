package main

import (
	database "backend/backend/pkg/db"
	"fmt"
)

func main() {
	fmt.Println("Welcome to Chatify!")

	// init mongodb
	database.InitMongo()
	defer database.CloseMongo()

	// local setup
	port := "8080"
	addr := ":"
}
