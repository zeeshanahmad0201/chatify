package user

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

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{
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

func FetchUserByEmail(email string) *models.User {
	filter := bson.M{models.UserKeyEmail: email}
	return FetchUserByFilter(filter)
}

func FetchUserByToken(token string) *models.User {
	filter := bson.M{models.UserKeyToken: token}
	return FetchUserByFilter(filter)
}

func FetchUserByFilter(filter primitive.M) *models.User {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	usersCollection := database.GetUsersCollection()

	var foundUser *models.User
	err := usersCollection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error fetching user: %v", err)
	}

	return foundUser
}

func AddUser(userModel *models.User) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	usersCollection := database.GetUsersCollection()

	_, err := usersCollection.InsertOne(ctx, userModel)
	if err != nil {
		fmt.Printf("InsertOne: %s", err.Error())
		return err
	}

	return nil
}
