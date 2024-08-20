package router

import (
	"github.com/gorilla/mux"
	"github.com/zeeshanahmad0201/chatify/backend/api/controllers"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/signup", controllers.Signup)

	return r
}
