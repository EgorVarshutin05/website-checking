package domain

type WebSite struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Status         string `json:"status"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	CheckedAt      string `json:"checked_at"`
}
