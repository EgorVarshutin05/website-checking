package main

import (
	"errors"
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
