package common

import (
	"testing"
)

// Helper function to create sample receipts
func createSampleReceipt(id, retailer, purchaseDate, purchaseTime, total string, items []Item) Receipt {
	return Receipt{
		ID:           id,
		Retailer:     retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Total:        total,
		Items:        items,
	}
}

func TestAddReceipt(t *testing.T) {
	rs := &ReceiptStorage{
		Receipts: make(map[string]Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	receipt := createSampleReceipt("1", "Retailer A", "2023-11-25", "12:00", "100.00", []Item{
		{"Item A", "50.00"},
		{"Item B", "50.00"},
	})

	// Add receipt
	err := rs.AddReceipt(receipt, 100)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	// Attempt to add the same receipt again
	err = rs.AddReceipt(receipt, 100)
	if err == nil || err.Error() != "receipt with ID 1 already exists" {
		t.Errorf("expected error 'receipt with ID 1 already exists', but got: %v", err)
	}
}

func TestGetAllReceipts(t *testing.T) {
	rs := &ReceiptStorage{
		Receipts: make(map[string]Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	// Test with no receipts
	_, err := rs.GetAllReceipts()
	if err == nil || err.Error() != "no receipts found" {
		t.Errorf("expected error 'no receipts found', but got: %v", err)
	}

	// Add a receipt and test retrieval
	receipt := createSampleReceipt("1", "Retailer A", "2023-11-25", "12:00", "100.00", []Item{
		{"Item A", "50.00"},
		{"Item B", "50.00"},
	})
	rs.AddReceipt(receipt, 100)

	receipts, err := rs.GetAllReceipts()
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if len(receipts) != 1 {
		t.Errorf("expected 1 receipt, but got: %d", len(receipts))
	}
	if receipts[0].ID != "1" {
		t.Errorf("expected receipt ID '1', but got: %s", receipts[0].ID)
	}
}

func TestGetReceiptByID(t *testing.T) {
	rs := &ReceiptStorage{
		Receipts: make(map[string]Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	// Test with no receipt
	_, err := rs.GetReceiptByID("1")
	if err == nil || err.Error() != "receipt with ID 1 not found" {
		t.Errorf("expected error 'receipt with ID 1 not found', but got: %v", err)
	}

	// Add a receipt and test retrieval
	receipt := createSampleReceipt("1", "Retailer A", "2023-11-25", "12:00", "100.00", []Item{
		{"Item A", "50.00"},
		{"Item B", "50.00"},
	})
	rs.AddReceipt(receipt, 100)

	retrievedReceipt, err := rs.GetReceiptByID("1")
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if retrievedReceipt.ID != "1" {
		t.Errorf("expected receipt ID '1', but got: %s", retrievedReceipt.ID)
	}
}

func TestGetReceiptPoints(t *testing.T) {
	rs := &ReceiptStorage{
		Receipts: make(map[string]Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	// Test with no receipt points
	_, err := rs.GetReceiptPoints("1")
	if err == nil || err.Error() != "points for receipt with ID 1 not found" {
		t.Errorf("expected error 'points for receipt with ID 1 not found', but got: %v", err)
	}

	// Add a receipt and test points retrieval
	receipt := createSampleReceipt("1", "Retailer A", "2023-11-25", "12:00", "100.00", []Item{
		{"Item A", "50.00"},
		{"Item B", "50.00"},
	})
	rs.AddReceipt(receipt, 100)

	points, err := rs.GetReceiptPoints("1")
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if points != 100 {
		t.Errorf("expected points 100, but got: %d", points)
	}
}

func TestUpdateReceipt(t *testing.T) {
	rs := &ReceiptStorage{
		Receipts: make(map[string]Receipt),
		Points:   make(map[string]int64),
		Order:    []string{},
	}

	// Add a receipt
	receipt := createSampleReceipt("1", "Retailer A", "2023-11-25", "12:00", "100.00", []Item{
		{"Item A", "50.00"},
		{"Item B", "50.00"},
	})
	rs.AddReceipt(receipt, 100)

	// Update the receipt
	updatedReceipt := createSampleReceipt("1", "Retailer B", "2023-11-26", "13:00", "200.00", []Item{
		{"Item C", "200.00"},
	})
	err := rs.UpdateReceipt("1", updatedReceipt, 200)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	// Verify the update
	retrievedReceipt, err := rs.GetReceiptByID("1")
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if retrievedReceipt.Retailer != "Retailer B" {
		t.Errorf("expected retailer 'Retailer B', but got: %s", retrievedReceipt.Retailer)
	}
	if retrievedReceipt.Total != "200.00" {
		t.Errorf("expected total '200.00', but got: %s", retrievedReceipt.Total)
	}
}
