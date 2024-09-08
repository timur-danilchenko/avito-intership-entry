package api

import (
	"net/http"
	"timur-danilchenko/avito-intership-entry/src/main/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRouter(router *mux.Router) {
	router.HandleFunc("/ping", handlers.PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/users", handlers.GetAllUsersHandler).Methods(http.MethodGet)
	router.HandleFunc("/user", handlers.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", handlers.GetUserByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", handlers.UpdateUserByIDHandler).Methods(http.MethodPut)
}

func RegisterAPIRouter(router *mux.Router) *mux.Router {
	subrouter := router.PathPrefix("/api").Subrouter()
	return subrouter
}
