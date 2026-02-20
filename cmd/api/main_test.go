package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	// create application and response recorder
	app := &application{}
	rr := httptest.NewRecorder()

	// create test data and headers
	data := envelope{
		"message": "hello",
	}
	headers := http.Header{}
	headers.Set("X-Test-Header", "test")

	// assert well-formed JSON
	err := app.writeJSON(rr, http.StatusOK, data, headers)
	if err != nil {
		t.Fatalf("writeJSON error: %v", err)
	}

	// assert 200 OK HTTP status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	// assert JSON Content-Type
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	// assert HTTP response body
	if !strings.Contains(rr.Body.String(), `"message": "hello"`) {
		t.Errorf("Response body does not contain expected JSON")
	}
}

func TestReadJSON(t *testing.T) {
	// create application and response recorder
	app := &application{}
	rr := httptest.NewRecorder()

	// set JSON in a HTTP request
	jsonBody := `{"name":"George"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(jsonBody))

	// define destination to read JSON into
	var input struct {
		Name string `json:"name"`
	}

	// assert readJSON success
	err := app.readJSON(rr, req, &input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// assert that Name is properly read
	if input.Name != "George" {
		t.Errorf("expected Name to be George, got %s", input.Name)
	}
}
