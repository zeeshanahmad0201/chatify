package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeeshanahmad0201/chatify/backend/models"
)

func StoreMessage(w http.ResponseWriter, r *http.Request) {
	// extract payload
	var message *models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// validate struct
	validate := validator.New()
	err = validate.Struct(&message)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		if len(errs) == 0 {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		firstErr := errs[0]
		field := firstErr.Field()
		message, exists := models.MessageValidationErrs[field]
		if !exists {
			message = firstErr.Error()
		}
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// send/store message
	err := 
}
