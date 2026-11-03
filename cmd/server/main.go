package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EgorVarshutin05/website-checking/internal/handler"
	"github.com/EgorVarshutin05/website-checking/internal/repository"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL не задан")
	}

	pool, err := repository.NewPostrgrePool(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer pool.Close()

	log.Println("Подключение к PostgreSQL установлено")

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
