// logger
package logger

import (
	"log"
	"time"
)

// Info logs informational messages
func Info(message string) {
	log.Printf("[INFO] %s", message)
}

// Error logs error messages
func Error(message string) {
	log.Printf("[ERROR] %s", message)
}

// LogRequest logs the HTTP method, URI, and the time taken to process the request
func LogRequest(method, uri string, start time.Time) {
	log.Printf("[INFO] %s %s completed in %v", method, uri, time.Since(start))
}
