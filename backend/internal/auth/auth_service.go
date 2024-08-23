package auth

import (
	"fmt"

	"github.com/zeeshanahmad0201/chatify/backend/internal/user"
	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(login *models.Login) (*models.User, error) {
	foundUser := user.FetchUserByEmail(login.Email)
	if foundUser == nil {
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

	foundUser.Token = token
	foundUser.RefreshToken = refreshToken

	return foundUser, nil
}

func Signup(signupReq *models.User) (string, error) {
	foundUser := user.FetchUserByEmail(signupReq.Email)
	if foundUser != nil {
		return "", fmt.Errorf("user already exists")
	}

	msg := "unable to add user to the database. please try again later!"

	hashedPassword, err := helpers.HashPassword(signupReq.Password)
	if err != nil {
		fmt.Printf("HashPassword: %s", err.Error())
		return "", fmt.Errorf(msg)
	}
	signupReq.Password = hashedPassword

	signupReq.ID = primitive.NewObjectID()
	signupReq.UserID = signupReq.ID.Hex()
	signupReq.AddedOn = helpers.GetCurrentTimeInMillis()

	token, refreshToken, err := helpers.GenerateTokens(signupReq.Name, signupReq.Email, signupReq.UserID)
	if err != nil {
		fmt.Printf("GenerateTokens: %s", err.Error())
		return "", fmt.Errorf(msg)
	}

	signupReq.Token = token
	signupReq.RefreshToken = refreshToken

	err = user.AddUser(signupReq)
	if err != nil {
		return "", fmt.Errorf(msg)
	}

	return "user has been created successfully", nil

}
