package message

import (
	"fmt"
	"log"

	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/database"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

// StoreMessage stores a message in the database.
func StoreMessage(message *models.Message) error {
	// Get database collections
	messageCollection := database.GetMsgsCollection()
	userCollection := database.GetUsersCollection()

	// Create a context with timeout
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	// Validate sender ID
	filter := bson.M{models.UserKeyUserID: message.SenderID}
	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error validating sender ID: %v", err)
		return fmt.Errorf("error validating sender ID")
	}
	if count == 0 {
		return fmt.Errorf("invalid sender ID")
	}

	// Validate receiver ID
	filter = bson.M{models.UserKeyID: message.ReceiverID}
	count, err = userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error validating receiver ID: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("invalid receiver ID")
	}

	// Set message timestamp and status
	message.Timestamp = helpers.GetCurrentTimeInMillis()
	message.Status = models.Sent

	// Store message in the database
	_, err = messageCollection.InsertOne(ctx, message)
	if err != nil {
		log.Printf("InsertOne failed: %s", err.Error())
		return fmt.Errorf("something went wrong while storing the message. Please try again later")
	}

	return nil
}
