package httpapi

import (
	"encoding/json"
	"jakarta-bbm-tracker/backend/internal/scraper"
	"log"
	"net/http"
)

func NewRouter(service *scraper.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/v1/prices", func(w http.ResponseWriter, r *http.Request) {
		cache, err := service.ReadCache()
		if err != nil {
			writeJSON(w, http.StatusOK, service.LastCache())
			return
		}
		writeJSON(w, http.StatusOK, cache)
	})
	mux.HandleFunc("POST /api/v1/scrape", func(w http.ResponseWriter, r *http.Request) {
		cache, err := service.ScrapeAndStore(r.Context())
		if err != nil {
			writeJSON(w, http.StatusAccepted, cache)
			return
		}
		writeJSON(w, http.StatusOK, cache)
	})
	return withCORS(mux)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
