package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethirajmudhaliar/GH-risk-api/common"
	"github.com/gorilla/mux"
)

func TestGetReceiptPointsSuccess(t *testing.T) {
	// Reset the global storage
	common.Storage = common.ReceiptStorage{
		Receipts: make(map[string]common.Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	// Add a receipt to the storage
	receipt := common.Receipt{ID: "1", Retailer: "Retailer A", PurchaseDate: "2023-11-25", PurchaseTime: "12:00", Total: "100.00"}
	common.Storage.AddReceipt(receipt, 150) // Assign 150 points to this receipt

	// Create a request to the endpoint
	req, err := http.NewRequest("GET", "/v1/receipts/1/points", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Set up the router and serve the request
	router := mux.NewRouter()
	router.HandleFunc("/v1/receipts/{id}/points", GetReceiptPoints)
	router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the response body for points
	expected := `{"success":true,"data":{"points":150}}`
	if rr.Body.String() != expected {
		t.Errorf("expected response body '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestGetReceiptPointsNotFound(t *testing.T) {
	// Reset the global storage
	common.Storage = common.ReceiptStorage{
		Receipts: make(map[string]common.Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	// Create a request to a non-existent receipt
	req, err := http.NewRequest("GET", "/v1/receipts/non-existent-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Set up the router and serve the request
	router := mux.NewRouter()
	router.HandleFunc("/v1/receipts/{id}/points", GetReceiptPoints)
	router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
	}

	// Check the response body for an error message
	expected := `{"success":false,"error":"Receipt not found"}`
	if rr.Body.String() != expected {
		t.Errorf("expected response body '%s', got '%s'", expected, rr.Body.String())
	}
}
