package validation

import (
	"errors"
	"regexp"
)

// Regular expressions for validation
var (
	retailerRegex         = regexp.MustCompile(`^[\w\s\-\&]+$`)       // Matches valid retailer names
	dateRegex             = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // Matches dates in YYYY-MM-DD format
	timeRegex             = regexp.MustCompile(`^\d{2}:\d{2}$`)       // Matches time in HH:mm format (24-hour)
	priceRegex            = regexp.MustCompile(`^\d+\.\d{2}$`)        // Matches prices like 123.45
	shortDescriptionRegex = regexp.MustCompile(`^[\w\s\-\']+$`)       // Matches valid item descriptions
)

// ValidateReceipt validates the fields of a receipt
func ValidateReceipt(retailer string, purchaseDate string, purchaseTime string, total string, items []map[string]string) error {
	// Validate retailer
	if retailer == "" || !retailerRegex.MatchString(retailer) {
		return errors.New("invalid retailer name")
	}

	// Validate purchase date
	if purchaseDate == "" || !dateRegex.MatchString(purchaseDate) {
		return errors.New("invalid purchase date format, expected YYYY-MM-DD")
	}

	// Validate purchase time
	if purchaseTime == "" || !timeRegex.MatchString(purchaseTime) {
		return errors.New("invalid purchase time format, expected HH:mm")
	}

	// Validate total amount
	if total == "" || !priceRegex.MatchString(total) {
		return errors.New("invalid total amount format, expected a decimal with two places")
	}

	// Validate items
	if len(items) == 0 {
		return errors.New("at least one item is required")
	}

	for _, item := range items {
		shortDescription := item["shortDescription"]
		price := item["price"]

		if shortDescription == "" || !shortDescriptionRegex.MatchString(shortDescription) {
			return errors.New("invalid short description for an item")
		}

		if price == "" || !priceRegex.MatchString(price) {
			return errors.New("invalid price for an item, expected a decimal with two places")
		}
	}

	return nil
}
