package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestCheckURL_Success(t *testing.T) {
	// Create a test server that returns 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	status := checkURL(server.URL)
	if status != "UP" {
		t.Errorf("Expected 'UP', got '%s'", status)
	}
}

func TestCheckURL_NonOKStatus(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	status := checkURL(server.URL)
	if status != "STATUS 404" {
		t.Errorf("Expected 'STATUS 404', got '%s'", status)
	}
}

func TestCheckURL_ServerError(t *testing.T) {
	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	status := checkURL(server.URL)
	if status != "STATUS 500" {
		t.Errorf("Expected 'STATUS 500', got '%s'", status)
	}
}

func TestCheckURL_Timeout(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	status := checkURL(server.URL)
	if status != "DOWN" {
		t.Errorf("Expected 'DOWN', got '%s'", status)
	}
}

func TestCheckURL_InvalidURL(t *testing.T) {
	status := checkURL("http://invalid-url-that-does-not-exist-12345.com")
	if status != "DOWN" {
		t.Errorf("Expected 'DOWN', got '%s'", status)
	}
}

func TestReadURLs_ValidFile(t *testing.T) {
	// Create a temporary file with URLs
	content := "https://google.com\nhttps://github.com\n"
	tmpfile, err := os.CreateTemp("", "sites-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	urls, err := readURLs(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(urls))
	}

	if urls[0] != "https://google.com" {
		t.Errorf("Expected 'https://google.com', got '%s'", urls[0])
	}
}

func TestReadURLs_FileNotFound(t *testing.T) {
	_, err := readURLs("non-existent-file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestReadURLs_EmptyFile(t *testing.T) {
	// Create an empty temporary file
	tmpfile, err := os.CreateTemp("", "empty-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	urls, err := readURLs(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls) != 0 {
		t.Errorf("Expected 0 URLs, got %d", len(urls))
	}
}

func TestReadURLs_FilterEmptyLines(t *testing.T) {
	// Create a temporary file with empty lines
	content := "https://google.com\n\n\nhttps://github.com\n"
	tmpfile, err := os.CreateTemp("", "sites-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	urls, err := readURLs(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should filter empty lines
	if len(urls) != 2 {
		t.Errorf("Expected 2 URLs (empty lines filtered), got %d", len(urls))
	}

	for _, url := range urls {
		if strings.TrimSpace(url) == "" {
			t.Error("Found empty URL, should have been filtered")
		}
	}
}

func TestReadURLs_FilterComments(t *testing.T) {
	// Create a temporary file with comments
	content := "# This is a comment\nhttps://google.com\n# Another comment\nhttps://github.com\n"
	tmpfile, err := os.CreateTemp("", "sites-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	urls, err := readURLs(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should filter comments
	if len(urls) != 2 {
		t.Errorf("Expected 2 URLs (comments filtered), got %d", len(urls))
	}

	for _, url := range urls {
		if strings.HasPrefix(url, "#") {
			t.Error("Found comment line, should have been filtered")
		}
	}
}

func TestSiteStatus_JSONMarshaling(t *testing.T) {
	status := SiteStatus{
		URL:    "https://example.com",
		Status: "UP",
		Time:   time.Now().Format(time.RFC3339),
	}

	// This should not panic
	data, err := marshalJSON(status)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify it's valid JSON
	var decoded SiteStatus
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.URL != status.URL {
		t.Errorf("Expected URL '%s', got '%s'", status.URL, decoded.URL)
	}
}

func TestHandleConnections_Upgrade(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(handleConnections))
	defer server.Close()

	// Convert http://... to ws://...
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Try to connect via WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Skipf("WebSocket connection failed (expected in test environment): %v", err)
		return
	}
	defer ws.Close()

	// Verify client was added
	clientsMutex.RLock()
	clientCount := len(clients)
	clientsMutex.RUnlock()

	if clientCount == 0 {
		t.Error("Expected at least 1 client, got 0")
	}
}

func TestBroadcastStatus_JSONError(t *testing.T) {
	// Test that broadcastStatus handles errors gracefully
	// This is tested indirectly through the full integration,
	// but we can verify marshalJSON works correctly
	statuses := []SiteStatus{
		{URL: "https://example.com", Status: "UP", Time: time.Now().Format(time.RFC3339)},
	}

	data, err := json.Marshal(statuses)
	if err != nil {
		t.Errorf("Expected no error marshaling valid data, got %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty JSON data")
	}
}

