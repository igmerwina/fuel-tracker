package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := loadConfig()
	client := newHTTPClient(cfg.ScrapeTimeout)
	service := NewScrapeService(client, cfg.CachePath)

	if _, err := service.ScrapeAndStore(backgroundContext()); err != nil {
		log.Printf("initial scrape completed with warnings: %v", err)
	}

	log.Printf("backend listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, newRouter(service)))
}
