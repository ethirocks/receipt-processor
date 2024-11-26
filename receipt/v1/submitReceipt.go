package v1

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethirajmudhaliar/GH-risk-api/common"
	"github.com/ethirajmudhaliar/GH-risk-api/logger"
	"github.com/ethirajmudhaliar/GH-risk-api/validation"
	"github.com/google/uuid"
)

// SubmitReceipt handles the submission of a receipt for processing
func SubmitReceipt(w http.ResponseWriter, r *http.Request) {
	var newReceipt common.Receipt

	// Parse the JSON body
	err := json.NewDecoder(r.Body).Decode(&newReceipt)
	if err != nil {
		logger.Error("Error decoding request body: " + err.Error())
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields in the receipt
	if newReceipt.Retailer == "" || newReceipt.PurchaseDate == "" || newReceipt.PurchaseTime == "" || newReceipt.Total == "" || len(newReceipt.Items) == 0 {
		logger.Error("Missing required fields in receipt submission")
		common.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Validate receipt fields using the validation package
	if err := validation.ValidateReceipt(
		newReceipt.Retailer,
		newReceipt.PurchaseDate,
		newReceipt.PurchaseTime,
		newReceipt.Total,
		convertItemsToMap(newReceipt.Items),
	); err != nil {
		logger.Error("Validation error: " + err.Error())
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate a new UUID for the receipt ID
	newReceipt.ID = generateUniqueID()

	// Calculate points (replace with your logic)
	points := calculatePoints(newReceipt)

	// Add the new receipt to the in-memory storage
	if err := common.Storage.AddReceipt(newReceipt, points); err != nil {
		logger.Error("Error adding receipt to storage: " + err.Error())
		common.RespondWithError(w, http.StatusInternalServerError, "Could not store the receipt")
		return
	}

	logger.Info("Receipt submitted successfully with ID: " + newReceipt.ID)

	// Respond with the newly created receipt ID
	response := common.JSONResponse{
		Success: true,
		Data:    map[string]string{"id": newReceipt.ID},
	}
	common.RespondWithJSON(w, http.StatusCreated, response)
}

// calculatePoints is a placeholder function to calculate receipt points

func calculatePoints(receipt common.Receipt) int64 {
	points := int64(0)

	// Rule 1: 1 point for each alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if isAlphanumeric(char) {
			points++
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount
	if isRoundDollar(receipt.Total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if isMultipleOf(receipt.Total, 0.25) {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += int64(len(receipt.Items) / 2 * 5)

	// Rule 5: Points for items with description length as multiple of 3
	for _, item := range receipt.Items {
		descLength := len(strings.TrimSpace(item.ShortDescription))
		if descLength%3 == 0 {
			itemPrice := parsePrice(item.Price)
			points += int64(math.Ceil(itemPrice * 0.2))
		}
	}

	// Rule 6: 6 points if the purchase day is odd
	if isOddDay(receipt.PurchaseDate) {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00 PM and 4:00 PM
	if isAfternoon(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

// Helper function to check if a character is alphanumeric
func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

// Helper function to check if the total is a round dollar amount
func isRoundDollar(total string) bool {
	price := parsePrice(total)
	return price == float64(int(price))
}

// Helper function to check if the total is a multiple of a given factor
func isMultipleOf(total string, factor float64) bool {
	price := parsePrice(total)
	return math.Mod(price, factor) == 0
}

// Helper function to parse a price string into a float
func parsePrice(price string) float64 {
	parsedPrice, _ := strconv.ParseFloat(price, 64)
	return parsedPrice
}

// Helper function to check if the purchase day is odd
func isOddDay(date string) bool {
	parts := strings.Split(date, "-")
	if len(parts) < 3 {
		return false
	}
	day, _ := strconv.Atoi(parts[2])
	return day%2 != 0
}

// Helper function to check if the purchase time is between 2:00 PM and 4:00 PM
func isAfternoon(time string) bool {
	parts := strings.Split(time, ":")
	if len(parts) < 2 {
		return false
	}
	hour, _ := strconv.Atoi(parts[0])
	return hour >= 14 && hour < 16
}

func convertItemsToMap(items []common.Item) []map[string]string {
	result := make([]map[string]string, len(items))
	for i, item := range items {
		result[i] = map[string]string{
			"shortDescription": item.ShortDescription,
			"price":            item.Price,
		}
	}
	return result
}

// Helper function to generate a unique ID
func generateUniqueID() string {
	for {
		id := uuid.New().String()
		if _, exists := common.Storage.Receipts[id]; !exists {
			return id
		}
	}
}
