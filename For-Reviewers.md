# Notes For Reviewers

## Core Features:
1. **HTTP Server**:
   - Listening on port 8080 for HTTP traffic.

2. **Endpoints**:
   - **POST /receipts/process**: Submits a new receipt for processing and points calculation.
   - **GET /receipts/{id}/points**: Retrieves the points awarded for a specific receipt by its ID.

3. **Receipt Data Structure**:
   - **Receipt ID**: UUID auto-generated upon submission.
   - **Retailer**: Name of the retailer issuing the receipt.
   - **PurchaseDate**: Date the receipt was issued (format: `YYYY-MM-DD`).
   - **PurchaseTime**: Time the receipt was issued (format: `HH:mm`, 24-hour).
   - **Total**: Total amount paid (format: `123.45`).
   - **Items**: List of items purchased, each with:
     - **shortDescription**: Description of the item.
     - **price**: Price of the item (format: `123.45`).

---

## Improvements Beyond Requirements:
1. **Modular Architecture**:
   - Core components separated into the following packages:
     - `v1`: API version-specific logic for receipts.
     - `common`: Shared utilities, response helpers, and in-memory storage.
     - `logger`: Logging functionality for requests and errors.
     - `validation`: Custom validation logic for receipt fields.

2. **Logging**:
   - Middleware logs HTTP requests with method, path, and processing time.
   - Logs errors and exceptions for debugging and monitoring.

3. **In-Memory Storage**:
   - Thread-safe storage using `sync.Mutex` to manage concurrent access.
   - Stores receipts and their associated points in separate maps for fast lookup.

4. **Validation**:
   - Ensures required fields like `retailer`, `purchaseDate`, `purchaseTime`, and `items` are present.
   - Validates field formats using regular expressions (e.g., `YYYY-MM-DD` for dates).

---

## Testing:
1. **Unit Tests**:
   - Thorough tests for core API handlers:
     - `SubmitReceipt`: Tests for valid receipts, missing fields, and invalid inputs.
     - `GetReceiptPoints`: Tests for successful lookups and missing receipts.
   - Includes validation tests for individual receipt fields.

---

## Extra Enhancements:
1. **Detailed Logging**:
   - Comprehensive logging for successful operations and errors.

2. **README Documentation**:
   - Includes clear instructions, usage examples, architecture details, and API references.

3. **Future Enhancements Suggested**:
   - Persistent database for receipt storage.
   - Token-based authentication for enhanced security.
   - Integration with external logging services for distributed monitoring.
