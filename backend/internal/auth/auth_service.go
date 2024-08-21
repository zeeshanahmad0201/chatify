package auth

import (
	"fmt"

	"github.com/zeeshanahmad0201/chatify/backend/internal/user"
	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/database"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(login *models.Login) (*models.User, error) {

	ctx, cancel := helpers.CreateContext()
	defer cancel()

	userCollection := database.GetUsersCollection()

	filter := bson.M{models.UserKeyEmail: login.Email}

	count, err := userCollection.CountDocuments(ctx, filter)

	if err != nil || count == 0 {
		return nil, fmt.Errorf("user not found")
	}

	hashedPassword, err := helpers.HashPassword(login.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid email/password")
	}
	login.Password = hashedPassword

	var foundUser *models.User

	err = userCollection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		fmt.Printf("FindOne: %s", err.Error())
		return nil, fmt.Errorf("invalid email/password")
	}

	validPass := helpers.VerifyPassword(&login.Password, &foundUser.Password)
	if !validPass {
		return nil, fmt.Errorf("invalid email/password")
	}

	token, refreshToken, err := helpers.GenerateTokens(foundUser.Name, foundUser.Email, foundUser.UserID)

	if err != nil {
		return nil, err
	}

	err = user.UpdateTokens(token, refreshToken, foundUser.UserID)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func Signup(user *models.User) (string, error) {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	userCollection := database.GetUsersCollection()

	filter := bson.M{models.UserKeyEmail: user.Email}

	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil || count > 0 {
		if err != nil {
			fmt.Printf("CountDocuments: %s", err.Error())
		}
		return "", fmt.Errorf("user already exists")
	}

	msg := "unable to add user to the database. please try again later!"

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		fmt.Printf("HashPassword: %s", err.Error())
		return "", fmt.Errorf(msg)
	}
	user.Password = hashedPassword

	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()

	token, refreshToken, err := helpers.GenerateTokens(user.Name, user.Email, user.UserID)
	if err != nil {
		fmt.Printf("GenerateTokens: %s", err.Error())
		return "", fmt.Errorf(msg)
	}

	user.Token = token
	user.RefreshToken = refreshToken

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		fmt.Printf("InsertOne: %s", err.Error())
		return "", fmt.Errorf(msg)
	}

	return "user has been created successfully", nil

}
