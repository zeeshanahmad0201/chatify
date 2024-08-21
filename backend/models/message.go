package models

import "go.mongodb.org/mongo-driver/bson/primitive"

var (
	MessageFieldSenderID   = "senderId"
	MessageFieldReceiverID = "receiverId"
	MessageFieldMessage    = "message"
)

var MessageValidationErrs = map[string]string{
	MessageFieldSenderID:   "sender id is required",
	MessageFieldReceiverID: "receiver id is required",
	MessageFieldMessage:    "message is required",
}

type Message struct {
	ID         primitive.ObjectID `bson:"_id"`
	SenderID   string             `json:"senderId" bson:"senderId" validate:"required"`
	ReceiverID string             `json:"receiverId" bson:"receiverId" validate:"required"`
	Message    string             `json:"message" bson:"message" validate:"required"`
	Timestamp  int64              `json:"timestamp" bson:"timestamp" validate:"required"`
	Status     MessageStatus      `json:"status" bson:"status" validate:"required"`
}

type MessageStatus string

const (
	Sent      MessageStatus = "sent"
	Delivered MessageStatus = "delivered"
	Read      MessageStatus = "read"
)
