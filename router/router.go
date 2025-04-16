package router

import (
	"FetchChallenge/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/receipts/process").Methods(http.MethodPost).HandlerFunc(handlers.PostReceiptHandler)
	r.Path("/receipts/{id}/points").Methods(http.MethodGet).HandlerFunc(handlers.GetPointsHandler)
	r.PathPrefix("/").HandlerFunc(handlers.NotFoundHandler)

	return r
}
