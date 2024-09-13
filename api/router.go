package api

import (
	"net/http"
	"timur-danilchenko/avito-intership-entry/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRouter(router *mux.Router) {
	router.HandleFunc("/ping", handlers.PingHandler).Methods(http.MethodGet)

	router.HandleFunc("/tenders", handlers.GetAllTendersHandler).Methods(http.MethodGet)
	router.HandleFunc("/tenders/new", handlers.CreateTenderHandler).Methods(http.MethodPost)
	router.HandleFunc("/tenders/my", handlers.GetUserTendersHandler).Methods(http.MethodGet)
	router.HandleFunc("/tenders/{tenderId}/status", handlers.GetTenderStatusHandler).Methods(http.MethodGet)
	router.HandleFunc("/tenders/{tenderId}/status", handlers.UpdateTenderStatusHandler).Methods(http.MethodPut)
	router.HandleFunc("/tenders/{tenderId}/edit", handlers.EditTenderHandler).Methods(http.MethodPatch)
	router.HandleFunc("/tenders/{tenderId}/rollback/{version}", handlers.RollbackTenderHandler).Methods(http.MethodPut)

	router.HandleFunc("/bids/new", handlers.CreateBidHandler).Methods(http.MethodPost)
	router.HandleFunc("/bids/my", handlers.GetUserBidsHandler).Methods(http.MethodGet)
	router.HandleFunc("/bids/{tenderId}/list", handlers.GetBidsForTenderHandler).Methods(http.MethodGet)
	router.HandleFunc("/bids/{bidId}/status", handlers.GetBidStatusHandler).Methods(http.MethodGet)
	router.HandleFunc("/bids/{bidId}/status", handlers.UpdateBidStatusHandler).Methods(http.MethodPut)
	router.HandleFunc("/bids/{bidId}/edit", handlers.EditBidHandler).Methods(http.MethodPatch)
	router.HandleFunc("/bids/{bidId}/submit_decision", handlers.SubmitBidDecisionHandler).Methods(http.MethodPut)
	router.HandleFunc("/bids/{bidId}/feedback", handlers.SubmitBidFeedbackHandler).Methods(http.MethodPut)
	router.HandleFunc("/bids/{bidId}/rollback/{version}", handlers.RollbackBidHandler).Methods(http.MethodPut)
	router.HandleFunc("/bids/{tenderId}/reviews", handlers.GetBidReviewsHandler).Methods(http.MethodGet)
}

func RegisterAPIRouter(router *mux.Router) *mux.Router {
	subrouter := router.PathPrefix("/api").Subrouter()
	return subrouter
}
