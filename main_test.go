package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethirajmudhaliar/GH-risk-api/common"
	v1 "github.com/ethirajmudhaliar/GH-risk-api/receipt/v1"
	"github.com/gorilla/mux"
)

func TestLoggingMiddleware(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/v1/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	router.Use(LoggingMiddleware)

	req, err := http.NewRequest("GET", "/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestRoutes(t *testing.T) {
	router := mux.NewRouter()

	router.HandleFunc("/receipts/process", v1.SubmitReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", v1.GetReceiptPoints).Methods("GET")

	req, err := http.NewRequest("GET", "/receipts/non-existent-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
	}
}

func TestCreateRiskRoute(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/v1/risks", v1.SubmitReceipt).Methods("POST")

	req, err := http.NewRequest("POST", "/v1/risks", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()

	common.Storage = common.ReceiptStorage{
		Receipts: make(map[string]common.Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	tests := []struct {
		method      string
		url         string
		body        []byte
		statusCode  int
		description string
	}{
		{"GET", "/receipts/non-existent-id/points", nil, http.StatusNotFound, "Get receipt points for non-existent receipt"},
		{"POST", "/receipts/process", []byte(`{"retailer": "Retailer A", "purchaseDate": "2023-11-25", "purchaseTime": "12:00", "total": "100.00", "items": [{"shortDescription": "Item A", "price": "50.00"}]}`), http.StatusCreated, "Submit a receipt"},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.body))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.statusCode {
			t.Errorf("%s: expected status code %d, got %d", tt.description, tt.statusCode, status)
		}
	}
}
