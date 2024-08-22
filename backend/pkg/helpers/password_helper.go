package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(providedPass *string, currentPass *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*currentPass), []byte(*providedPass))

	if err != nil {
		fmt.Printf("invalid password %s", err)
	}

	return err == nil

}

// HashPassword hashes the password and return
func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(hashedPass), nil

}
