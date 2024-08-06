package pkg

import (
	"backend/backend/pkg/helpers"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

func InitMongo() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		connectionURL := os.Getenv("DB_URI")
		if connectionURL == "" {
			log.Fatal("DB_URI environment variable is not set")
		}

		clientOptions := options.Client().ApplyURI(connectionURL)

		ctx, cancel := helpers.CreateContext()
		defer cancel()

		// connect with mongo
		clientInstance, clientInstanceError := mongo.Connect(ctx, clientOptions)
		if clientInstanceError != nil {
			log.Fatal(clientInstanceError.Error())
		}

		err := clientInstance.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	return clientInstance, nil
}
