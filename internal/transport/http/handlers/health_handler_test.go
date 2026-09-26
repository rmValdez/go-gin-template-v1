package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/handlers"
)

func TestHealthHandler(t *testing.T) {
	// Set Gin to test mode so it doesn't print out debug info
	gin.SetMode(gin.TestMode)

	// Create a new router
	r := gin.New()
	
	// Register the handler exactly like the real app does, but directly
	r.GET("/health", handlers.Health)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	// Parse the response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v\n", err)
	}

	// Check if the success boolean is true (based on our standard JSON response format)
	if success, ok := response["success"].(bool); !ok || !success {
		t.Fatalf("Expected success to be true, got %v", response["success"])
	}

	// Check if data contains message "OK"
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object in response")
	}

	if msg, ok := data["status"].(string); !ok || msg != "ok" {
		t.Fatalf("Expected status 'ok', got '%v'", data["status"])
	}
}
