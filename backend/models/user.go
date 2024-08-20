package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	UserKeyUserID       = "userId"
	UserKeyToken        = "token"
	UserKeyRefreshToken = "refreshToken"
	UserKeyUpdatedOn    = "updatedOn"
	UserKeyEmail        = "email"
)

var UserValidationErrs = map[string]string{
	"Name":     "name is required",
	"Email":    "email is required",
	"Password": "password is required",
}

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" validate:"required"`
	Email        string             `json:"email" validate:"required"`
	Password     string             `json:"password" validate:"required,min=6"`
	AddedOn      int64              `json:"addedOn" bson:"added_on"`
	UpdatedOn    int64              `json:"updatedOn" bson:"updated_on"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refreshToken" bson:"refresh_token"`
	UserID       string             `json:"userId" bson:"user_id"`
}
