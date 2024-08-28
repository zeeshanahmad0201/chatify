package router

import (
	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/chatify/backend/api/controllers"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/signup", controllers.Signup).Methods("POST")

	// Get all messages
	r.HandleFunc("/messages", controllers.GetMessages).Methods("GET")

	// Message actions
	r.HandleFunc("/message/send", controllers.StoreMessage).Methods("POST")
	r.HandleFunc("/message/delete/{messageId}", controllers.DeleteMessage).Methods("DELETE")
	r.HandleFunc("/message/status/{messageId}/read", controllers.MarkAsRead).Methods("PUT")
	r.HandleFunc("/message/status/{messageId}/delivered", controllers.MarkAsRead).Methods("PUT")

	// web socket
	r.HandleFunc("/messages/listen", controllers.MessageWebsocket)

	return r
}
