package checker

import (
	"context"
	"net/http"
	"time"
)

func GetStatusCode(ctx context.Context, client *http.Client, url string) (int, int64, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	start := time.Now()

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return 0, time.Since(start).Milliseconds(), err
	}
	defer resp.Body.Close()

	return resp.StatusCode, duration.Milliseconds(), nil
}
