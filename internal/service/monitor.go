package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/EgorVarshutin05/website-checking/internal/checker"
	"github.com/EgorVarshutin05/website-checking/internal/domain"
	"github.com/EgorVarshutin05/website-checking/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MonitorService struct {
	pool   *pgxpool.Pool
	client *http.Client
}

func NewMonitorService(pool *pgxpool.Pool, client *http.Client) *MonitorService {
	return &MonitorService{
		pool:   pool,
		client: client,
	}
}

func (s *MonitorService) GetSites(ctx context.Context) ([]domain.WebSite, error) {

	var wg sync.WaitGroup

	sites, err := repository.GetAllSites(ctx, s.pool)
	if err != nil {
		return nil, err
	}

	for i, site := range sites {
		wg.Go(
			func() {
				statusCode, duration, err := checker.GetStatusCode(ctx, s.client, site.Url)
				// log.Printf(
				// 	"Проверка сайта: url=%s, statusCode=%d, duration=%dms, err=%v",
				// 	site.Url,
				// 	statusCode,
				// 	duration,
				// 	err,
				// )

				sites[i].Status = GetStatusUD(statusCode, err)
				sites[i].ResponseTimeMs = duration
				sites[i].CheckedAt = time.Now().UTC().Format(time.RFC3339)
			},
		)
	}

	wg.Wait()

	return sites, nil
}

func (s *MonitorService) CreateSite(ctx context.Context, req domain.CreateSiteRequest) (domain.WebSite, error) {
	name := strings.TrimSpace(req.Name)
	if len(name) == 0 {
		return domain.WebSite{}, errors.New("название сайта не может быть пустым")
	}
	req.Name = name

	url := strings.TrimSpace(req.URL)
	if len(url) == 0 {
		return domain.WebSite{}, errors.New("ссылка не может быть пустой")
	}
	req.URL = url

	newSite, err := repository.CreateSite(ctx, s.pool, req)
	if err != nil {
		return domain.WebSite{}, err
	}

	return newSite, nil

}

func GetStatusUD(statusCode int, err error) string {
	if statusCode > 399 || err != nil || statusCode < 200 {
		return "down"
	}
	return "up"
}
