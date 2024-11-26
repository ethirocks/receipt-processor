// store common structs
package common

// Receipt represents the structure of a receipt.
type Receipt struct {
	ID           string `json:"id"`           // Unique identifier for the receipt
	Retailer     string `json:"retailer"`     // Retailer's name
	PurchaseDate string `json:"purchaseDate"` // Date of purchase
	PurchaseTime string `json:"purchaseTime"` // Time of purchase
	Items        []Item `json:"items"`        // List of purchased items
	Total        string `json:"total"`        // Total purchase amount
}

// Item represents an item within a receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"` // Item description
	Price            string `json:"price"`            // Price of the item
}
