package v1

import (
	"net/http"

	"github.com/ethirajmudhaliar/GH-risk-api/common"
	"github.com/ethirajmudhaliar/GH-risk-api/logger"
	"github.com/gorilla/mux"
)

// GetReceiptPoints retrieves the points awarded for a specific receipt by its ID
func GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	vars := mux.Vars(r)
	receiptID := vars["id"]

	// Retrieve the points for the given receipt ID from the storage
	points, err := common.Storage.GetReceiptPoints(receiptID)
	if err != nil {
		logger.Info("Receipt with ID not found: " + receiptID)
		logger.Error("Error: " + err.Error())
		common.RespondWithError(w, http.StatusNotFound, "Receipt not found")
		return
	}

	// Log the successful retrieval
	logger.Info("Returning points for receipt ID: " + receiptID)

	// Wrap the points in a JSONResponse and respond
	response := common.JSONResponse{
		Success: true,
		Data:    map[string]int64{"points": points},
	}
	common.RespondWithJSON(w, http.StatusOK, response)
}
