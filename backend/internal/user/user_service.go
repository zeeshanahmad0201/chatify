package user

import (
	"backend/backend/models"
	"backend/backend/pkg/database"
	"backend/backend/pkg/helpers"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateTokens update the token of the user with provided userId
func UpdateTokens(token string, refreshToken string, userId string) error {

	if token == "" || refreshToken == "" {
		return fmt.Errorf("no valid tokens provided")
	}

	ctx, cancel := helpers.CreateContext()
	defer cancel()

	userCollection := database.GetUsersCollection()
	if userCollection == nil {
		return fmt.Errorf("invalid userCollection")
	}

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: models.UserKeyToken, Value: token})
	updateObj = append(updateObj, bson.E{Key: models.UserKeyRefreshToken, Value: refreshToken})
	updateObj = append(updateObj, bson.E{Key: models.UserKeyUpdatedOn, Value: helpers.GetCurrentTimeInMillis()})

	filter := bson.M{models.UserKeyUserID: userId}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateByID(ctx, filter, bson.D{
		{
			Key:   "$set",
			Value: updateObj,
		},
	}, &opt)

	if err != nil {
		return fmt.Errorf("unable to update the tokens for the user")
	}

	return nil
}
