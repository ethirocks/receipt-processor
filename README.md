# Receipt Processing API

The **Receipt Processing API** is a web-based service designed to process and manage receipts. It provides an HTTP-based interface to submit receipts, retrieve receipt points, and manage receipt-related data, with data persisted in memory for the purposes of this assignment.

---

## Project Overview

The API supports the following operations on receipts:

- **Submit a new receipt**: `POST /receipts/process`
- **Retrieve points for a receipt**: `GET /receipts/{id}/points`

---

## Receipt Data Structure

Each receipt contains the following fields:

| Field             | Type   | Description                                                   |
|-------------------|--------|---------------------------------------------------------------|
| `ID`              | UUID   | Unique identifier for the receipt (generated on submission).  |
| `Retailer`        | string | The name of the retailer who issued the receipt.              |
| `PurchaseDate`    | string | The date the receipt was issued (format: `YYYY-MM-DD`).       |
| `PurchaseTime`    | string | The time the receipt was issued (format: `HH:mm`, 24-hour).   |
| `Total`           | string | The total amount paid on the receipt (format: `123.45`).      |
| `Items`           | array  | List of items purchased (each item has `shortDescription` and `price`). |

---

## Features

- **Receipt Operations**:
  - Submit receipts for processing and points calculation.
  - Retrieve the points awarded for a specific receipt.

- **Logging**:
  - Logs each HTTP request and its processing time.
  - Logs any errors that occur during processing.

- **In-Memory Data Storage**:
  - All receipts are stored in memory (`map[string]Receipt`).
  - Points are calculated and stored in a separate `map[string]int64`.

---

## Endpoints

### 1. `POST /receipts/process`

**Description**: Submit a new receipt for processing and points calculation.

**Request Body**:

```json
{
  "retailer": "Retailer A",
  "purchaseDate": "2023-11-25",
  "purchaseTime": "12:00",
  "total": "100.00",
  "items": [
    {"shortDescription": "Item A", "price": "50.00"},
    {"shortDescription": "Item B", "price": "50.00"}
  ]
}
```

**Response**:

- `201 Created`: Returns the ID of the processed receipt.
- `400 Bad Request`: If the input data is invalid.

---

### 2. `GET /receipts/{id}/points`

**Description**: Retrieve the points awarded for a specific receipt by its ID.

**Response**:

- `200 OK`: Returns the points awarded for the receipt.
- `404 Not Found`: If the receipt is not found.

---

## Running the Project

### Prerequisites

- **Go**: Ensure Go is installed on your machine.
- **Docker**: For containerized setup.

---

### Steps to Run Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/ethirajmudhaliar/GH-receipt-api.git
   cd GH-receipt-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

4. Test the API:
   Use the provided `curl` examples or Postman.

---

### Steps to Run with Docker

1. Build the Docker image:
   ```bash
   docker build -t receipt-processor .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 receipt-processor
   ```

3. Test the API:
   Use the provided `curl` examples or Postman.

4. Stop the container:
   Press `Ctrl+C` or use:
   ```bash
   docker ps
   docker stop <container_id>
   ```

---

## Testing

### Run Tests
Run the unit tests for the project:
```bash
go test ./... -v
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

---

## Architecture

### 1. **main.go**

- Initializes the HTTP server and sets up the routes and middleware.
- **Router Setup**: Defines API routes (`/receipts/process`, `/receipts/{id}/points`) using the Gorilla Mux router.
- **Logging Middleware**: Logs details of incoming requests and their processing time.

### 2. **v1 Package**

Contains the business logic for handling receipt operations, including:
- `SubmitReceipt`: Handles the submission of a new receipt.
- `GetReceiptPoints`: Retrieves the points for a specific receipt by its ID.

### 3. **common Package**

Includes in-memory storage and response helper functions:
- `ReceiptStorage`: In-memory storage using a map and slice for receipts.
- `RespondWithJSON` & `RespondWithError`: Functions to standardize JSON responses and error handling.

### 4. **logger Package**

Provides logging capabilities for the application:
- `Info`: Logs informational messages.
- `Error`: Logs error messages.
- `LogRequest`: Logs HTTP requests, including method, URL, and processing time.

### 5. **validation Package**

Contains validation logic for the API:
- **Receipt Validation**: Ensures that the receipt fields are valid (`retailer`, `purchaseDate`, `purchaseTime`, `total`, `items`).

---

## API Example Usage

### Submit a New Receipt

```bash
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{"retailer": "Retailer A", "purchaseDate": "2023-11-25", "purchaseTime": "12:00", "total": "100.00", "items": [{"shortDescription": "Item A", "price": "50.00"}, {"shortDescription": "Item B", "price": "50.00"}]}'
```

### Get Points for a Receipt

```bash
curl http://localhost:8080/receipts/{id}/points
```

Replace `{id}` with the ID returned from the `/receipts/process` endpoint.

---

## Rules for Point Calculation

1. 1 point for each alphanumeric character in the retailer name.
2. 50 points if the total is a round dollar amount.
3. 25 points if the total is a multiple of `0.25`.
4. 5 points for every two items on the receipt.
5. If an item’s description length is a multiple of 3, the item’s price earns `price * 0.2` points (rounded up).
6. 6 points if the purchase day is odd.
7. 10 points if the purchase time is between 2:00 PM and 4:00 PM.

---
