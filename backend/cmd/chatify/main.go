package main

import (
	"backend/backend/api/router"
	database "backend/backend/pkg/db"
	"backend/backend/pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	server       *http.Server
	shutdownOnce sync.Once
)

func main() {
	fmt.Println("Welcome to Chatify!")

	// init mongodb
	database.InitMongo()
	defer database.CloseMongo()

	// local setup
	initLocalServer()

}

func initLocalServer() {
	port := "8080"
	addr := ":" + port

	r := router.InitRouter()

	server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// channel for interupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server %v", err)
		}
	}()

	<-signalChan

	shutdownOnce.Do(func() {
		// Perform graceful shutdown
		fmt.Println("Shutting down server...")

		ctx, cancel := helpers.CreateContext()
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}

	})
}
