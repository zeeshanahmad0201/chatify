package router

import (
	controller "backend/backend/api/controllers"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", controller.Login)

	return r
}
