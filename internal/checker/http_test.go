package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatusCode500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := server.Client()

	code, _, err := GetStatusCode(context.Background(), client, server.URL)

	if code != http.StatusInternalServerError || err != nil {
		t.Errorf("Ожидался HTTP 500, получено: status=%d, err=%v", code, err)
	}
}

func TestGetStatusCodeConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	url := server.URL
	client := server.Client()
	server.Close()

	code, _, err := GetStatusCode(context.Background(), client, url)

	if code != 0 || err == nil {
		t.Errorf("Ожидался status=0 и ошибка, получено: status=%d, err=%v", code, err)
	}
}
