package auth

import (
	"fmt"

	"github.com/zeeshanahmad0201/chatify/backend/internal/user"
	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/database"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(login *models.Login) (*models.User, error) {
	userCollection := database.GetUsersCollection()

	ctx, cancel := helpers.CreateContext()
	defer cancel()

	filter := bson.M{"email": login.Email}

	var foundUser *models.User

	err := userCollection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		fmt.Printf("user not found %s", err.Error())
		return nil, fmt.Errorf("user not found")
	}

	validPass := helpers.VerifyPassword(&login.Password, &foundUser.Password)
	if !validPass {
		return nil, fmt.Errorf("invalid email/password")
	}

	token, refreshToken, err := helpers.GenerateTokens(foundUser.Name, foundUser.Email, foundUser.UserID)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	err = user.UpdateTokens(token, refreshToken, foundUser.UserID)
	if err != nil {
		return nil, err
	}

	return foundUser, nil

}
