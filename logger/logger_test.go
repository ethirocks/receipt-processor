package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

func TestInfoLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	Info("This is an info message")

	if !strings.Contains(buf.String(), "[INFO] This is an info message") {
		t.Errorf("expected log to contain '[INFO] This is an info message', got %s", buf.String())
	}
}

func TestErrorLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	Error("This is an error message")

	if !strings.Contains(buf.String(), "[ERROR] This is an error message") {
		t.Errorf("expected log to contain '[ERROR] This is an error message', got %s", buf.String())
	}
}

func TestLogRequest(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	startTime := time.Now()
	LogRequest("GET", "/test", startTime)

	if !strings.Contains(buf.String(), "[INFO] GET /test completed") {
		t.Errorf("expected log to contain '[INFO] GET /test completed', got %s", buf.String())
	}
}
