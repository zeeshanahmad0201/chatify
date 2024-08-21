package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeeshanahmad0201/chatify/backend/internal/message"
	"github.com/zeeshanahmad0201/chatify/backend/models"
)

func StoreMessage(w http.ResponseWriter, r *http.Request) {
	// extract payload
	var messageObj *models.Message
	err := json.NewDecoder(r.Body).Decode(&messageObj)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// validate struct
	validate := validator.New()
	err = validate.Struct(&messageObj)
	if err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
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
	err = message.StoreMessage(messageObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("message sent successfully!")
}
