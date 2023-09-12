package rest

import (
	"context"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/internal/rest/accesstokens"
	"heypay-cash-in-server/internal/rest/depositvalidations"
	"heypay-cash-in-server/internal/rest/receivedtransfers"
	"heypay-cash-in-server/internal/rest/reversedtransfers"
	"heypay-cash-in-server/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/{rest:.*}", corsHandler).Methods("OPTIONS")

	r.HandleFunc("/", homeHandler).Methods("GET")

	r.HandleFunc("/access-tokens", accessTokenPostHandler).Methods("POST")
	r.HandleFunc("/deposit-validations", depositValidationsPostHandler).Methods("POST")
	r.HandleFunc("/received-transfers", receivedTransfersPostHandler).Methods("POST")
	r.HandleFunc("/reversed-transfers", reversedTransfersPostHandler).Methods("POST")

	r.HandleFunc("/accessTokens", accessTokenPostHandler).Methods("POST")
	r.HandleFunc("/depositValidations", depositValidationsPostHandler).Methods("POST")
	r.HandleFunc("/receivedTransfers", receivedTransfersPostHandler).Methods("POST")
	r.HandleFunc("/reversedTransfers", reversedTransfersPostHandler).Methods("POST")
	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		ctx = context.WithValue(ctx, constants.ContextKeyIPAddress, utils.GetIP(r))
		ctx = context.WithValue(ctx, constants.ContextKeyLoggerID, id.String())
		ctx = context.WithValue(ctx, constants.ContextKeyStartTime, time.Now().UnixNano())
		r = r.WithContext(ctx)
		utils.HttpLogRequest(r)
		next.ServeHTTP(w, r)
	})
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	home(w, r)
}

func accessTokenPostHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	accesstokens.PostAccessTokens(w, r)
}

func depositValidationsPostHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	depositvalidations.PostDepositValidations(w, r)
}

func receivedTransfersPostHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	receivedtransfers.PostReceivedTransfers(w, r)
}

func reversedTransfersPostHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	reversedtransfers.PostReversedTransfers(w, r)
}
