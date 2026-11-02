package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatusUD(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		expected   string
	}{
		{"HTTP_200", 200, nil, "up"},
		{"HTTP_500", 500, nil, "down"},
		{"TIMEOUT", 200, errors.New("timeout"), "down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusUD(tt.statusCode, tt.err)

			if result != tt.expected {
				t.Errorf("Ожидалось: %s, получено: %s", tt.expected, result)
			}
		})
	}
}

func TestGetStatusCode500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	code, _, err := getStatusCode(server.URL)
	if code != 500 || err != nil {
		t.Errorf("ожидался HTTP 500, получено: status=%d, err=%v", code, err)
	}
}

func TestGetStatusCodeConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	url := server.URL
	server.Close()

	code, _, err := getStatusCode(url)

	if code != 0 || err == nil {
		t.Errorf("Ожидился 0, поулчено: status=%d, err=%v", code, err)
	}
}
