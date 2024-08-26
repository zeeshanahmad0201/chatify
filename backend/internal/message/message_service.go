package message

import (
	"fmt"
	"log"

	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/database"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	filter = bson.M{models.UserKeyUserID: message.ReceiverID}
	count, err = userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error validating receiver ID: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("invalid receiver ID")
	}

	message.ID = primitive.NewObjectID()
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

func FetchMessages(senderId string, receiverId string) []models.Message {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	filter := bson.M{
		models.MessageFieldReceiverID: receiverId,
		models.MessageFieldSenderID:   senderId,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: models.MessageFieldTimestamp, Value: 1}})

	messagesCollection := database.GetMsgsCollection()
	messages := make([]models.Message, 0)

	cursor, err := messagesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("error finding messages: %v", err)
		return messages
	}

	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			log.Printf("error decoding messages: %v", err)
			continue
		}
		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("cursor error: %v", err)
		return messages
	}

	return messages
}

func DeleteMessage(messageId string) error {
	msgCollection := database.GetMsgsCollection()

	ctx, cancel := helpers.CreateContext()
	defer cancel()

	filter := bson.M{models.MessageFieldID: messageId}

	_, err := msgCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("error while deleting the message: %v", err)
		return fmt.Errorf("invalid message id")
	}

	return nil
}
