package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeeshanahmad0201/chatify/backend/internal/auth"
	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var loginReq *models.Login

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "unable to decode payload", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(loginReq); err != nil {
		log.Printf("error while validating loginReq: %s", err.Error())
		http.Error(w, helpers.GetValidationErrMsg(err), http.StatusBadRequest)
		return
	}

	user, err := auth.Login(loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func Signup(w http.ResponseWriter, r *http.Request) {

	var user *models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		log.Println("error:", err.Error())
		http.Error(w, helpers.GetValidationErrMsg(err), http.StatusBadRequest)
		return
	}

	msg, err := auth.Signup(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}
