package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ScrapeService struct {
	client    *http.Client
	cachePath string
	lastCache PriceCache
}

func NewScrapeService(client *http.Client, cachePath string) *ScrapeService {
	return &ScrapeService{
		client:    client,
		cachePath: cachePath,
		lastCache: emptyCache(),
	}
}

func (s *ScrapeService) LastCache() PriceCache {
	return s.lastCache
}

func (s *ScrapeService) ReadCache() (PriceCache, error) {
	payload, err := os.ReadFile(s.cachePath)
	if err != nil {
		return PriceCache{}, err
	}
	var cache PriceCache
	if err := json.Unmarshal(payload, &cache); err != nil {
		return PriceCache{}, err
	}
	s.lastCache = cache
	return cache, nil
}

func (s *ScrapeService) ScrapeAndStore(ctx context.Context) (PriceCache, error) {
	cache := emptyCache()

	for _, source := range sources() {
		result, err := s.scrapeSource(ctx, source)
		if err != nil {
			cache.Errors = append(cache.Errors, fmt.Sprintf("%s: %v", source.Brand, err))
			continue
		}
		cache.Prices = append(cache.Prices, result.Prices...)
		cache.Stations = append(cache.Stations, result.Stations...)
	}

	if err := s.writeCache(cache); err != nil {
		return cache, err
	}
	s.lastCache = cache

	if len(cache.Errors) > 0 {
		return cache, errors.New(strings.Join(cache.Errors, "; "))
	}
	return cache, nil
}

type SourceResult struct {
	Prices   []FuelPrice
	Stations []StationLocation
}

func (s *ScrapeService) scrapeSource(ctx context.Context, source Source) (SourceResult, error) {
	switch source.Brand {
	case "Pertamina":
		return pricesOnly(scrapePertamina(ctx, s.client, source))
	case "Shell":
		return pricesOnly(scrapeShell(ctx, s.client, source))
	case "BP":
		return pricesOnly(scrapeBP(ctx, s.client, source))
	case "Vivo":
		return scrapeVivo(ctx, s.client, source)
	default:
		return SourceResult{}, fmt.Errorf("unsupported source brand %q", source.Brand)
	}
}

func pricesOnly(prices []FuelPrice, err error) (SourceResult, error) {
	if err != nil {
		return SourceResult{}, err
	}
	return SourceResult{Prices: prices}, nil
}

func (s *ScrapeService) writeCache(cache PriceCache) error {
	if err := os.MkdirAll(filepath.Dir(s.cachePath), 0o755); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.cachePath, payload, 0o644)
}

func emptyCache() PriceCache {
	return PriceCache{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Regions:     jakartaRegions,
		Prices:      []FuelPrice{},
		Stations:    []StationLocation{},
		Errors:      []string{},
	}
}
