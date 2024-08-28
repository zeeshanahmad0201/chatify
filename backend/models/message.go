package models

import "go.mongodb.org/mongo-driver/bson/primitive"

var (
	MessageFieldID         = "_id"
	MessageFieldSenderID   = "senderId"
	MessageFieldReceiverID = "receiverId"
	MessageFieldMessage    = "message"
	MessageFieldTimestamp  = "timestamp"
	MessageFieldStatus     = "status"
)

var MessageValidationErrs = map[string]string{
	MessageFieldSenderID:   "sender id is required",
	MessageFieldReceiverID: "receiver id is required",
	MessageFieldMessage:    "message is required",
}

type Message struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	SenderID   string             `json:"senderId" bson:"senderId"`
	ReceiverID string             `json:"receiverId" bson:"receiverId" validate:"required"`
	Message    string             `json:"message" bson:"message" validate:"required"`
	Timestamp  int64              `json:"timestamp" bson:"timestamp"`
	Status     MessageStatus      `json:"status" bson:"status"`
}

type MessageStatus int

const (
	Sent MessageStatus = iota
	Delivered
	Read
)
