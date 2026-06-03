package main

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port          string
	CachePath     string
	ScrapeTimeout time.Duration
}

func loadConfig() Config {
	return Config{
		Port:          getenv("PORT", "8080"),
		CachePath:     filepath.Join("data", "fuel_prices.json"),
		ScrapeTimeout: scrapeTimeout(),
	}
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func backgroundContext() context.Context {
	return context.Background()
}

func scrapeTimeout() time.Duration {
	raw := getenv("SCRAPE_TIMEOUT_SECONDS", "20")
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		seconds = 20
	}
	return time.Duration(seconds) * time.Second
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(getenvRaw(key))
	if value == "" {
		return fallback
	}
	return value
}
