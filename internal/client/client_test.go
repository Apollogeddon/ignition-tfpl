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
			Enabled: boolPtr(true),
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
				Enabled:   boolPtr(true),
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
		Enabled: boolPtr(true),
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

func TestClient_UpdateResource(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		if r.URL.Path != "/data/api/v1/resources/ignition/test-type" {
			t.Errorf("Expected path /data/api/v1/resources/ignition/test-type, got %s", r.URL.Path)
		}

		response := []ResourceResponse[map[string]any]{
			{
				Name:      "test-name",
				Enabled:   boolPtr(true),
				Signature: "sig-456",
				Config:    map[string]any{"key": "updated"},
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
		Enabled: boolPtr(true),
		Config:  map[string]any{"key": "updated"},
	}
	var dest ResourceResponse[map[string]any]
	err = c.UpdateResource(context.Background(), "test-type", item, &dest)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if dest.Signature != "sig-456" {
		t.Errorf("Expected signature sig-456, got %s", dest.Signature)
	}
}

func TestClient_DeleteResource(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/data/api/v1/resources/ignition/test-type/test-name/sig-123" {
			t.Errorf("Expected path /data/api/v1/resources/ignition/test-type/test-name/sig-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Client
	c, err := NewClient(server.URL, "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test
	err = c.DeleteResource(context.Background(), "test-type", "test-name", "sig-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_ProjectOperations(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/data/api/v1/projects/find/test-project" {
				t.Errorf("Expected path /data/api/v1/projects/find/test-project, got %s", r.URL.Path)
			}
			p := Project{Name: "test-project", Enabled: true}
			json.NewEncoder(w).Encode(p)
		case http.MethodPost:
			if r.URL.Path != "/data/api/v1/projects" {
				t.Errorf("Expected path /data/api/v1/projects, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			if r.URL.Path != "/data/api/v1/projects/test-project" {
				t.Errorf("Expected path /data/api/v1/projects/test-project, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			if r.URL.Path != "/data/api/v1/projects/test-project" {
				t.Errorf("Expected path /data/api/v1/projects/test-project, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// Client
	c, err := NewClient(server.URL, "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test Create
	p := Project{Name: "test-project", Enabled: true}
	created, err := c.CreateProject(ctx, p)
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}
	if created.Name != "test-project" {
		t.Errorf("Expected name test-project, got %s", created.Name)
	}

	// Test Update
	updated, err := c.UpdateProject(ctx, p)
	if err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}
	if updated.Name != "test-project" {
		t.Errorf("Expected name test-project, got %s", updated.Name)
	}

	// Test Delete
	err = c.DeleteProject(ctx, "test-project")
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
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

func boolPtr(b bool) *bool {
	return &b
}

