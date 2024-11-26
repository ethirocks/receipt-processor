package common

import (
	"fmt"
	"sync"
)

// ReceiptStorage holds receipts in memory with fast lookup and insertion order tracking.
type ReceiptStorage struct {
	Receipts map[string]Receipt // Map for fast lookups
	Points   map[string]int64   // Map for storing points associated with receipts
	Order    []string           // Slice to store receipt IDs in order of insertion
	mu       sync.Mutex         // Mutex to handle concurrent access
}

// Global instance of the in-memory receipt storage.
var Storage = ReceiptStorage{
	Receipts: make(map[string]Receipt),
	Points:   make(map[string]int64),
	Order:    []string{},
}

// AddReceipt adds a new receipt to the storage.
func (rs *ReceiptStorage) AddReceipt(receipt Receipt, points int64) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, exists := rs.Receipts[receipt.ID]; exists {
		return fmt.Errorf("receipt with ID %s already exists", receipt.ID)
	}

	rs.Receipts[receipt.ID] = receipt
	rs.Points[receipt.ID] = points
	rs.Order = append(rs.Order, receipt.ID)

	return nil
}

// GetAllReceipts returns all receipts in insertion order.
func (rs *ReceiptStorage) GetAllReceipts() ([]Receipt, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if len(rs.Receipts) == 0 {
		return nil, fmt.Errorf("no receipts found")
	}

	receiptList := make([]Receipt, 0, len(rs.Receipts))
	for _, receiptID := range rs.Order {
		receiptList = append(receiptList, rs.Receipts[receiptID])
	}
	return receiptList, nil
}

// GetReceiptByID retrieves a specific receipt by ID.
func (rs *ReceiptStorage) GetReceiptByID(id string) (Receipt, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	receipt, exists := rs.Receipts[id]
	if !exists {
		return Receipt{}, fmt.Errorf("receipt with ID %s not found", id)
	}
	return receipt, nil
}

// GetReceiptPoints retrieves points for a specific receipt by ID.
func (rs *ReceiptStorage) GetReceiptPoints(id string) (int64, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	points, exists := rs.Points[id]
	if !exists {
		return 0, fmt.Errorf("points for receipt with ID %s not found", id)
	}
	return points, nil
}

// UpdateReceipt updates an existing receipt in the storage.
func (rs *ReceiptStorage) UpdateReceipt(id string, receipt Receipt, points int64) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	_, exists := rs.Receipts[id]
	if !exists {
		return fmt.Errorf("receipt with ID %s not found", id)
	}

	rs.Receipts[id] = receipt
	rs.Points[id] = points
	return nil
}
