package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	payload := JSONResponse{
		Success: true,
		Data:    map[string]string{"message": "test successful"},
	}
	RespondWithJSON(rr, http.StatusOK, payload)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	// Parse the response body
	var response JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Verify the JSON response payload
	if response.Success != true {
		t.Errorf("expected Success true, got %v", response.Success)
	}

	if data, ok := response.Data.(map[string]interface{}); ok {
		if data["message"] != "test successful" {
			t.Errorf("expected response message 'test successful', got %v", data["message"])
		}
	} else {
		t.Errorf("unexpected data format: %v", response.Data)
	}
}

func TestRespondWithError(t *testing.T) {
	rr := httptest.NewRecorder()

	errorMessage := "something went wrong"
	RespondWithError(rr, http.StatusBadRequest, errorMessage)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	// Parse the response body
	var response JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Verify the JSON response payload
	if response.Success != false {
		t.Errorf("expected Success false, got %v", response.Success)
	}

	if response.Error != errorMessage {
		t.Errorf("expected error message '%s', got '%s'", errorMessage, response.Error)
	}
}
