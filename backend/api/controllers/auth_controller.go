package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeeshanahmad0201/chatify/backend/internal/auth"
	"github.com/zeeshanahmad0201/chatify/backend/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.Login

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "unable to decode payload", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(&loginReq); err != nil {
		for _, value := range err.(validator.ValidationErrors) {
			field := value.Field()

			message, exists := models.LoginValidationErrs[field]
			if !exists {
				message = value.Error()
			}
			http.Error(w, message, http.StatusBadRequest)
			break
		}
	}

	user, err := auth.Login(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
