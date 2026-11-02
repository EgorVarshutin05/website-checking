package server

import (
	"log"
	"net/http"
	"time"

	"github.com/EgorVarshutin05/website-checking/internal/handler"
)

func main() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/sites", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSitesHandler(w, r, client)
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
