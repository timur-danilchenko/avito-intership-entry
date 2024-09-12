package main

import (
	"fmt"
	"net/http"
	"timur-danilchenko/avito-intership-entry/api"
	"timur-danilchenko/avito-intership-entry/utilities"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := api.RegisterAPIRouter(mux.NewRouter())

	api.RegisterRouter(router)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	serverUrl := utilities.GetEnv("SERVER_ADDRESS", "localhost:8080")

	log.Info(fmt.Sprintf("Server started on %s", serverUrl))

	http.ListenAndServe(serverUrl, router)
}
