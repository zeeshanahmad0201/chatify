package router

import (
	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/chatify/backend/api/controllers"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/message/send", controllers.StoreMessage).Methods("POST")
	r.HandleFunc("/messages", controllers.GetMessages).Methods("GET")

	return r
}
