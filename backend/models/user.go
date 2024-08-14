package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id"`
	Name         string             `json:"name" validator:"required"`
	Email        string             `json:"email" validator:"required"`
	Password     string             `validate:"required,min=6"`
	AddedOn      int64              `json:"addedOn"`
	UpdatedOn    int64              `json:"updatedOn"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshToken"`
	UserID       string             `json:"userId`
}
