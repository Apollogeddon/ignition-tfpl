package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetResource(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.Header.Get("X-Ignition-API-Token") != "test-token" {
			t.Errorf("Expected token test-token, got %s", r.Header.Get("X-Ignition-API-Token"))
		}
		if r.URL.Path != "/data/api/v1/resources/find/ignition/test-type/test-name" {
			t.Errorf("Expected path /data/api/v1/resources/find/ignition/test-type/test-name, got %s", r.URL.Path)
		}

		response := ResourceResponse[map[string]any]{
			Name:    "test-name",
			Enabled: true,
			Config:  map[string]any{"key": "value"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Client
	c, err := NewClient(server.URL, "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test
	var dest ResourceResponse[map[string]any]
	err = c.GetResource(context.Background(), "test-type", "test-name", &dest)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if dest.Name != "test-name" {
		t.Errorf("Expected name test-name, got %s", dest.Name)
	}
	if dest.Config["key"] != "value" {
		t.Errorf("Expected config key=value, got %v", dest.Config["key"])
	}
}

func TestClient_CreateResource(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/data/api/v1/resources/ignition/test-type" {
			t.Errorf("Expected path /data/api/v1/resources/ignition/test-type, got %s", r.URL.Path)
		}

		// API returns array of created resources
		response := []ResourceResponse[map[string]any]{
			{
				Name:      "test-name",
				Enabled:   true,
				Signature: "sig-123",
				Config:    map[string]any{"key": "value"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Client
	c, err := NewClient(server.URL, "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test
	item := ResourceResponse[map[string]any]{
		Name:    "test-name",
		Enabled: true,
		Config:  map[string]any{"key": "value"},
	}
	var dest ResourceResponse[map[string]any]
	err = c.CreateResource(context.Background(), "test-type", item, &dest)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if dest.Signature != "sig-123" {
		t.Errorf("Expected signature sig-123, got %s", dest.Signature)
	}
}

func TestClient_APIError(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := APIErrorResponse{
			Success: false,
			Problem: &struct {
				Message    string   `json:"message"`
				StackTrace []string `json:"stacktrace,omitempty"`
			}{
				Message: "Resource already exists",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Client
	c, err := NewClient(server.URL, "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test
	var dest ResourceResponse[map[string]any]
	err = c.GetResource(context.Background(), "test-type", "test-name", &dest)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedMsg := "API error: Resource already exists"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}
