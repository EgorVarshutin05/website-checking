package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/EgorVarshutin05/website-checking/internal/checker"
	"github.com/EgorVarshutin05/website-checking/internal/domain"
	"github.com/EgorVarshutin05/website-checking/internal/service"
)

func GetSitesHandler(w http.ResponseWriter, r *http.Request, client *http.Client) {

	var wg sync.WaitGroup
	ctx := r.Context()
	urls := []string{
		"https://mail.ru/",
		"https://cs.petrsu.ru/studies/index.php.ru",
		"https://petrsu.ru/",
	}

	sites := make([]domain.WebSite, len(urls))

	for i, site := range urls {

		wg.Add(1)
		go func() {
			defer wg.Done()

			statusCode, duration, err := checker.GetStatusCode(ctx, client, site)
			if err != nil {
				log.Printf("Ошибка функции getStatusCode: %v", err)
			}
			status := service.GetStatusUD(statusCode, err)

			sites[i] = domain.WebSite{
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
