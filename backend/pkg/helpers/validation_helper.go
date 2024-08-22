package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Define custom messages for specific fields and tags
var FieldErrorMessages = map[string]map[string]string{
	"Email": {
		"required": "email is required",
		"email":    "invalid email",
	},
	"Name": {
		"required": "name is required",
	},
	"Password": {
		"required": "password is required",
		"min":      "password must be at least 6 characters long",
	},
}

func GetValidationErrMsg(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		firstErr := validationErrs[0]
		field := firstErr.Field()
		tag := firstErr.Tag()

		if fieldMsgs, exists := FieldErrorMessages[field]; exists {
			if msg, ok := fieldMsgs[tag]; ok {
				return msg
			}
		}

		return fmt.Sprintf("Validation failed on field '%s' with rule '%s'", field, tag)
	}

	return "Invalid Payload"
}
