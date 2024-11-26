// store common structs
package common

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a standard API response format.
type JSONResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Data payload, optional
	Error   string      `json:"error,omitempty"`   // Error message, optional
	Message string      `json:"message,omitempty"` // Additional message, optional
}

// RespondWithJSON sends a JSON response.
func RespondWithJSON(w http.ResponseWriter, status int, payload JSONResponse) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// RespondWithError sends an error response.
func RespondWithError(w http.ResponseWriter, status int, errorMessage string) {
	RespondWithJSON(w, status, JSONResponse{
		Success: false,
		Error:   errorMessage,
	})
}

// RespondWithSuccess sends a success response.
func RespondWithSuccess(w http.ResponseWriter, status int, data interface{}, message string) {
	RespondWithJSON(w, status, JSONResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}
