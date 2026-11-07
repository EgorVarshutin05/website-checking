package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/EgorVarshutin05/website-checking/internal/checker"
	"github.com/EgorVarshutin05/website-checking/internal/handler"
	"github.com/EgorVarshutin05/website-checking/internal/repository"
	"github.com/EgorVarshutin05/website-checking/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL не задан")
	}

	pool, err := repository.NewPostrgresPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer pool.Close()

	log.Println("Подключение к PostgreSQL установлено")

	client := checker.NewSecureHTTPClient()

	mux := http.NewServeMux()

	monitorService := service.NewMonitorService(pool, client)
	siteHandler := handler.NewSiteHandler(monitorService)

	mux.HandleFunc("GET /api/v1/sites", siteHandler.GetSites)
	mux.HandleFunc("POST /api/v1/sites", siteHandler.CreateSite)

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web"))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
