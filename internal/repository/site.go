package repository

import (
	"context"

	"github.com/EgorVarshutin05/website-checking/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllSites(ctx context.Context, pool *pgxpool.Pool) ([]domain.WebSite, error) {
	rows, err := pool.Query(ctx, "SELECT id, name, url FROM sites ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sites := make([]domain.WebSite, 0)

	for rows.Next() {
		var site domain.WebSite

		err := rows.Scan(&site.ID, &site.Name, &site.Url)
		if err != nil {
			return nil, err
		}

		sites = append(sites, site)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

func CreateSite(ctx context.Context, pool *pgxpool.Pool, req domain.CreateSiteRequest) (domain.WebSite, error) {
	row := pool.QueryRow(ctx, "INSERT INTO sites (name, url) VALUES ($1, $2) RETURNING id, name, url", req.Name, req.URL)

	var site domain.WebSite

	err := row.Scan(
		&site.ID,
		&site.Name,
		&site.Url,
	)

	if err != nil {
		return domain.WebSite{}, err
	}

	return site, nil

}
