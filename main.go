package main

import (
	"heypay-cash-in-server/internal/rest"
	"heypay-cash-in-server/utils"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	utils.Info(nil, "Starting the service...")
	router := rest.Router()
	utils.Info(nil, "The service is ready to listen and serve.")
	log.Fatal(http.ListenAndServe(":8000", router))
}
