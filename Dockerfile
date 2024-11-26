# Step 1: Use an official Go image as the builder
FROM golang:1.20 AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum, and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the source code to the working directory
COPY . .

# Step 5: Build the application
RUN go build -o receipt-processor .

# Step 6: Use Ubuntu 22.04 for running the application
FROM ubuntu:22.04

# Install dependencies (glibc)
RUN apt-get update && apt-get install -y libc6

# Set up a working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/receipt-processor .

# Expose the port the app listens on
EXPOSE 8080

# Command to run the application
CMD ["./receipt-processor"]
