package checker

import (
	"context"
	"testing"
)

func TestSecureHTTPClientBlocksLocalhost(t *testing.T) {
	client := NewSecureHTTPClient()

	code, _, err := GetStatusCode(
		context.Background(),
		client,
		"http://127.0.0.1:8080",
	)

	if err == nil {
		t.Fatal("ожидалась ошибка для локального адреса")
	}

	if code != 0 {
		t.Errorf("ожидался статус 0, получено %d", code)
	}
}

func TestSecureHTTPClientBlocksPrivateIP(t *testing.T) {
	client := NewSecureHTTPClient()

	code, _, err := GetStatusCode(
		context.Background(),
		client,
		"http://192.168.1.1",
	)

	if err == nil {
		t.Fatal("ожидалась ошибка для частного IP")
	}

	if code != 0 {
		t.Errorf("ожидался статус 0, получено %d", code)
	}
}

func TestSecureHTTPClientBlocksPort(t *testing.T) {
	client := NewSecureHTTPClient()

	_, _, err := GetStatusCode(
		context.Background(),
		client,
		"http://1.1.1.1:5432",
	)

	if err == nil {
		t.Fatal("ожидалась ошибка для запрещённого порта")
	}
}
