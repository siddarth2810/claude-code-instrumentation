package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHooksEndpointRecordsValidHookInput(t *testing.T) {
	body := `{
		"session_id": "session-123",
		"transcript_path": "/tmp/transcript.jsonl",
		"cwd": "/home/sid/project",
		"hook_event_name": "UserPromptSubmit",
		"prompt": "write a test"
	}`

	request := httptest.NewRequest(http.MethodPost, "/hooks", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	NewHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var got map[string]bool
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatalf("response was not valid JSON: %v", err)
	}

	if !got["recorded"] {
		t.Fatalf("recorded = %v, want true", got["recorded"])
	}
}

func TestHooksEndpointRejectsGet(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/hooks", nil)
	response := httptest.NewRecorder()

	NewHandler().ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
}

func TestHooksEndpointRejectsInvalidHookInput(t *testing.T) {
	body := `{"hook_event_name": "SomethingElse"}`
	request := httptest.NewRequest(http.MethodPost, "/hooks", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	NewHandler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
