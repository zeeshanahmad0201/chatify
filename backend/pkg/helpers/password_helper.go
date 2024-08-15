package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(providedPass *string, currentPass *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*providedPass), []byte(*currentPass))

	if err != nil {
		fmt.Printf("invalid password %s", err)
	}

	return err != nil

}
