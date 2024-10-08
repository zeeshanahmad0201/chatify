package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/zeeshanahmad0201/chatify/backend/internal/message"
	"github.com/zeeshanahmad0201/chatify/backend/internal/user"
	"github.com/zeeshanahmad0201/chatify/backend/models"
	"github.com/zeeshanahmad0201/chatify/backend/pkg/helpers"
)

func StoreMessage(w http.ResponseWriter, r *http.Request) {
	// get token from headers
	authHeaders := r.Header.Get("Authorization")
	if authHeaders == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}
	token := helpers.ExtractToken(authHeaders)

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

	// fetch user based on the token
	sender := user.FetchUserByToken(token)
	if sender == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	messageObj.SenderID = sender.UserID

	if messageObj.SenderID == messageObj.ReceiverID {
		http.Error(w, "You can't send messages to yourself", http.StatusBadRequest)
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

func GetMessages(w http.ResponseWriter, r *http.Request) {
	// get receiver id from query params
	receiverID := r.URL.Query().Get("receiverId")
	if receiverID == "" {
		http.Error(w, "receiverId is required", http.StatusBadRequest)
		return
	}

	// get token from headers
	token := helpers.ExtractTokenFromRequest(r)
	if token == "" {
		log.Printf("no token found %v", r.Header)
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// fetch user based on the token
	sender := user.FetchUserByToken(token)
	if sender == nil {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// fetch messages based on the token of the sender and id of the receiver
	messages := message.FetchMessages(sender.UserID, receiverID)

	if len(messages) == 0 {
		http.Error(w, "no messages found!", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

func MessageWebsocket(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			// Allow connections only from localhost
			return origin == "http://localhost" || origin == "http://127.0.0.1"
		},
	}

	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", con)
		http.Error(w, "Failed to upgrade websocket", http.StatusInternalServerError)
		return
	}
	defer con.Close()

	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			log.Printf("Websocket read error: %v", err)
			return
		}

		log.Printf("Message received: %v", message)

		err = con.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Websocket write error: %v", err)
			break
		}
	}
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// extract token
	token := helpers.ExtractTokenFromRequest(r)
	if token == "" {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// validate token
	userInfo := user.FetchUserByToken(token)
	if userInfo == nil {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// extract message id
	messageId := mux.Vars(r)["messageId"]
	if messageId == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// validate message id with sender id
	_, err := message.FetchMessageByUserID(messageId, userInfo.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// delete message
	err = message.DeleteMessage(messageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("message deleted successfully")
}

func MarkAsRead(w http.ResponseWriter, r *http.Request) {
	updateStatus(w, r, models.Read)
}

func MarkAsDelivered(w http.ResponseWriter, r *http.Request) {
	updateStatus(w, r, models.Delivered)
}

func updateStatus(w http.ResponseWriter, r *http.Request, newStatus models.MessageStatus) {
	// extract token
	token := helpers.ExtractTokenFromRequest(r)
	if token == "" {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// validate token
	userInfo := user.FetchUserByToken(token)
	if userInfo == nil {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// extract the message id
	msgId := r.URL.Query().Get("messageId")
	if msgId == "" {
		http.Error(w, "message id is missing from the route", http.StatusBadRequest)
		return
	}

	msg, err := message.FetchMessageByUserID(msgId, userInfo.UserID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusBadRequest)
		return
	}

	if newStatus <= msg.Status {
		var statusMessage string
		switch msg.Status {
		case models.Sent:
			statusMessage = "message already sent"
		case models.Delivered:
			statusMessage = "message already delivered"
		case models.Read:
			statusMessage = "message already read"
		}
		http.Error(w, statusMessage, http.StatusBadRequest)
		return
	}

	err = message.UpdateMessageStatus(msgId, newStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("message status updated")
}
