package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type WebSite struct {
	ID             int
	Name           string
	Url            string
	Status         string
	ResponseTimeMs int64
	CheckedAt      string
}

func getStatusCode(url string) (int, int64, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return http.StatusBadGateway, 0, err
	}

	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	duration := time.Since(start)

	statusCode := resp.StatusCode

	return statusCode, duration.Milliseconds(), nil
}

func getStatusUD(statusCode int) string {
	if statusCode > 399 {
		return "down"
	}
	return "up"
}

func main() {
	url := "https://nosuchsite12345.ru"
	statusCode, duration, _ := getStatusCode(url)
	state := getStatusUD(statusCode)

	result := WebSite{
		ID:             1,
		Name:           "VK",
		Url:            url,
		Status:         state,
		ResponseTimeMs: duration,
	}

	fmt.Printf("URL: %s,\nСостояние: %s,\nСтатус-код: %d,\nВремя ответа: %d мс\n", result.Url, result.Status, statusCode, result.ResponseTimeMs)
}
