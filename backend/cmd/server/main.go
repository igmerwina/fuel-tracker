package main

import (
	"context"
	"jakarta-bbm-tracker/backend/internal/config"
	"jakarta-bbm-tracker/backend/internal/httpapi"
	"jakarta-bbm-tracker/backend/internal/scraper"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	client := config.NewHTTPClient(cfg.ScrapeTimeout)
	service := scraper.NewService(client, cfg.CachePath)

	if _, err := service.ScrapeAndStore(context.Background()); err != nil {
		log.Printf("initial scrape completed with warnings: %v", err)
	}

	log.Printf("backend listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, httpapi.NewRouter(service)))
}
