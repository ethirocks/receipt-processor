package validation

import "testing"

func TestValidateReceiptValid(t *testing.T) {
	err := ValidateReceipt(
		"M&M Corner Market",
		"2023-11-25",
		"13:45",
		"12.34",
		[]map[string]string{
			{"shortDescription": "Apples", "price": "5.00"},
			{"shortDescription": "Bananas", "price": "7.34"},
		},
	)

	if err != nil {
		t.Errorf("expected no error for valid receipt, got: %v", err)
	}
}

func TestValidateReceiptMissingFields(t *testing.T) {
	err := ValidateReceipt("", "2023-11-25", "13:45", "12.34", []map[string]string{})
	if err == nil {
		t.Errorf("expected error for missing retailer and items, but got none")
	}
}

func TestValidateReceiptInvalidDate(t *testing.T) {
	err := ValidateReceipt("M&M Corner Market", "25-11-2023", "13:45", "12.34", []map[string]string{
		{"shortDescription": "Apples", "price": "5.00"},
	})
	if err == nil {
		t.Errorf("expected error for invalid date format, but got none")
	}
}

func TestValidateReceiptInvalidPrice(t *testing.T) {
	err := ValidateReceipt("M&M Corner Market", "2023-11-25", "13:45", "12.34", []map[string]string{
		{"shortDescription": "Apples", "price": "5"},
	})
	if err == nil {
		t.Errorf("expected error for invalid item price format, but got none")
	}
}
