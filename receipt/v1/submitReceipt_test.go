package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethirajmudhaliar/GH-risk-api/common"
	"github.com/gorilla/mux"
)

func TestSubmitReceiptSuccess(t *testing.T) {
	// Reset the global storage
	common.Storage = common.ReceiptStorage{
		Receipts: make(map[string]common.Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	payload := `{
		"retailer": "Retailer A",
		"purchaseDate": "2023-11-25",
		"purchaseTime": "12:00",
		"total": "100.00",
		"items": [
			{"shortDescription": "Item A", "price": "50.00"},
			{"shortDescription": "Item B", "price": "50.00"}
		]
	}`

	req, err := http.NewRequest("POST", "/v1/receipts/process", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitReceipt)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, status)
	}

	var response common.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	if response.Success != true {
		t.Errorf("expected Success true, got %v", response.Success)
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected Data to be a map, got %T", response.Data)
	}

	if _, exists := data["id"]; !exists {
		t.Errorf("expected 'id' in response data, but not found")
	}
}

func TestSubmitReceiptMissingFields(t *testing.T) {
	common.Storage = common.ReceiptStorage{
		Receipts: make(map[string]common.Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	payload := `{"retailer": "Retailer A"}`
	req, err := http.NewRequest("POST", "/v1/receipts/process", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitReceipt)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	var response common.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	if response.Success != false {
		t.Errorf("expected Success false, got %v", response.Success)
	}

	if response.Error != "Missing required fields" {
		t.Errorf("expected error 'Missing required fields', got '%s'", response.Error)
	}
}

func TestSubmitReceiptInvalidPrice(t *testing.T) {
	common.Storage = common.ReceiptStorage{
		Receipts: make(map[string]common.Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	payload := `{
		"retailer": "Retailer A",
		"purchaseDate": "2023-11-25",
		"purchaseTime": "12:00",
		"total": "100.00",
		"items": [
			{"shortDescription": "Item A", "price": "50"},  // Invalid price
			{"shortDescription": "Item B", "price": "50.00"}
		]
	}`

	req, err := http.NewRequest("POST", "/v1/receipts/process", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitReceipt)

	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	// Check the response body for an error message
	expectedError := "Invalid request payload"
	var response common.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	if response.Error != expectedError {
		t.Errorf("expected validation error '%s', got '%s'", expectedError, response.Error)
	}
}

func TestSubmitReceiptInvalidPayload(t *testing.T) {
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte("invalid-json")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", SubmitReceipt)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expected := `{"success":false,"error":"Invalid request payload"}`
	if rr.Body.String() != expected {
		t.Errorf("expected response '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestSubmitReceiptLargePayload(t *testing.T) {
	// Create a receipt with 1000 items
	items := make([]map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = map[string]string{
			"shortDescription": fmt.Sprintf("Item %d", i),
			"price":            "1.00",
		}
	}

	payload := map[string]interface{}{
		"retailer":     "Retailer A",
		"purchaseDate": "2023-11-25",
		"purchaseTime": "12:00",
		"total":        "1000.00",
		"items":        items,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", SubmitReceipt)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, status)
	}
}

func TestSubmitReceiptMissingFields1(t *testing.T) {
	payload := `{"retailer": "Retailer A"}`
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", SubmitReceipt)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expected := `{"success":false,"error":"Missing required fields"}`
	if rr.Body.String() != expected {
		t.Errorf("expected response '%s', got '%s'", expected, rr.Body.String())
	}
}
