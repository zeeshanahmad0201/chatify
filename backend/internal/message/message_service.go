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

	msgStr, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		log.Printf("Invalid message id: %v", err)
		return fmt.Errorf("invalid message ID")
	}

	filter := bson.M{models.MessageFieldID: msgStr}

	_, err = msgCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("error while deleting the message: %v", err)
		return fmt.Errorf("invalid message id")
	}

	return nil
}

func FetchMessageByUserID(messageId string, senderId string) (*models.Message, error) {
	msgsCollection := database.GetMsgsCollection()

	ctx, cancel := helpers.CreateContext()
	defer cancel()

	msgStr, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		log.Printf("Invalid message id: %v", err)
		return nil, fmt.Errorf("invalid message ID")
	}

	filter := bson.M{
		models.MessageFieldID:       msgStr,
		models.MessageFieldSenderID: senderId,
	}

	var msg *models.Message
	err = msgsCollection.FindOne(ctx, filter).Decode(&msg)
	if err != nil {
		log.Printf("error while fetching the doc: %v", err)
		return nil, fmt.Errorf("unable to find the message")
	}

	return msg, nil
}

func UpdateMessageStatus(messageId string, status models.MessageStatus) error {
	msgCollection := database.GetMsgsCollection()

	ctx, cancel := helpers.CreateContext()
	defer cancel()

	msgId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		log.Printf("error while converting to ObjectID: %v", err)
		return fmt.Errorf("invalid message id")
	}
	filter := bson.M{
		models.MessageFieldID: msgId,
	}

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: models.MessageFieldStatus, Value: models.MessageStatus(status)})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = msgCollection.UpdateOne(ctx, filter, bson.D{
		{
			Key: "$set", Value: updateObj,
		},
	}, &opt)
	if err != nil {
		log.Printf("error while updating the message status: %v", err)
		return fmt.Errorf("unable to update the message status")
	}

	return nil
}
