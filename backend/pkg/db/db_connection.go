package database

import (
	"backend/backend/pkg/helpers"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	mongoOnce      sync.Once
)

func InitMongo() *mongo.Client {
	mongoOnce.Do(func() {
		connectionURL := os.Getenv("DB_URI")
		if connectionURL == "" {
			log.Fatal("DB_URI environment variable is not set")
		}

		clientOptions := options.Client().ApplyURI(connectionURL)

		ctx, cancel := helpers.CreateContext()
		defer cancel()

		// connect with mongo
		clientInstance, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = clientInstance.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	return clientInstance
}

func GetDBName() string {
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		log.Fatal("DB_NAME environment variable not set")
	}

	return dbname
}

func GetCollection(collEnv string) *mongo.Collection {
	if collEnv == "" {
		log.Fatalln("collEnv can't be empty")
	}

	collName := os.Getenv(collEnv)

	if collName == "" {
		log.Fatalf("%s environment variable is not set", collEnv)
	}

	client := InitMongo()
	dbName := GetDBName()

	return client.Database(dbName).Collection(collName)
}

func GetUsersCollection() *mongo.Collection {
	return GetCollection("DB_USER_COLLECTION")
}

func CloseMongo() {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	clientInstance.Disconnect(ctx)
}
