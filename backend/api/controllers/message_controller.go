package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeeshanahmad0201/chatify/backend/internal/message"
	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
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
	err = validate.Struct(messageObj)
	if err != nil {
		log.Printf("error while validating messageObj: %s", err)
		http.Error(w, helpers.GetValidationErrMsg(err), http.StatusBadRequest)
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
