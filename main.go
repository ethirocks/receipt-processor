package main

import (
	"net/http"
	"time"

	"github.com/ethirajmudhaliar/GH-risk-api/logger"
	v1 "github.com/ethirajmudhaliar/GH-risk-api/receipt/v1"
	"github.com/gorilla/mux"
)

// LoggingMiddleware logs details about incoming HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Info("Started " + r.Method + " " + r.RequestURI)

		next.ServeHTTP(w, r)

		logger.LogRequest(r.Method, r.RequestURI, start)
	})
}

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Define the routes for the Receipt Processor API
	router.HandleFunc("/receipts/process", v1.SubmitReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", v1.GetReceiptPoints).Methods("GET")

	// Add the logging middleware
	router.Use(LoggingMiddleware)

	return router
}

func main() {
	router := SetupRouter()

	logger.Info("Starting server on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Error("Error starting server: " + err.Error())
	}
}
