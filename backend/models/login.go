package models

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var LoginValidationErrs = map[string]string{
	"Email":    "Invalid email",
	"Password": "Invalid password",
}
