package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type WebSite struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Status         string `json:"status"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	CheckedAt      string `json:"checked_at"`
}

func getSitesHandler(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	ctx := r.Context()
	urls := []string{
		"https://mail.ru/",
		"https://cs.petrsu.ru/studies/index.php.ru",
		"https://petrsu.ru/",
	}

	sites := make([]WebSite, len(urls))

	for i, site := range urls {

		wg.Add(1)
		go func() {
			defer wg.Done()

			statusCode, duration, err := getStatusCode(ctx, site)
			if err != nil {
				log.Printf("Ошибка функции getStatusCode: %v", err)
			}
			status := getStatusUD(statusCode, err)

			sites[i] = WebSite{
				ID:             i + 1,
				Name:           "unnamed",
				Url:            site,
				Status:         status,
				ResponseTimeMs: duration,
				CheckedAt:      time.Now().UTC().Format(time.RFC3339),
			}
		}()

	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sites); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
	}

}

func getStatusCode(ctx context.Context, url string) (int, int64, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, time.Since(start).Milliseconds(), err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	_, err = io.Copy(io.Discard, resp.Body)
	duration := time.Since(start)
	if err != nil {
		return 0, duration.Milliseconds(), err
	}

	return statusCode, duration.Milliseconds(), nil
}

func getStatusUD(statusCode int, err error) string {
	if statusCode > 399 || err != nil || statusCode < 200 {
		return "down"
	}
	return "up"
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/sites", getSitesHandler)

	log.Println("HTTP-сервер запущен на порту 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
