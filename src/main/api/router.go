package api

import (
	"timur-danilchenko/avito-intership-entry/src/main/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRouter(router *mux.Router) {
	router.HandleFunc("/ping", handlers.PingHandler)
}
