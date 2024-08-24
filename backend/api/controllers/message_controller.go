package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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

		// TODO: complete the websocket
		log.Printf("Message received: %v", message)

		err = con.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Websocket write error: %v", err)
			break
		}
	}
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	token := helpers.ExtractTokenFromRequest(r)
	if token == "" {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

}
