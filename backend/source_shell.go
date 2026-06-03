package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func scrapeShell(ctx context.Context, client *http.Client, source Source) ([]FuelPrice, error) {
	urls := []string{source.URL}
	if strings.HasSuffix(source.URL, ".html") {
		urls = append(urls, strings.TrimSuffix(source.URL, ".html")+".model.json")
	}

	var lastErr error
	for _, url := range urls {
		nextSource := source
		nextSource.URL = url
		body, err := fetchBody(ctx, client, url)
		if err != nil {
			lastErr = err
			continue
		}
		prices := parseGenericFuelPrices(nextSource, string(body), false)
		if len(prices) > 0 {
			return prices, nil
		}
		lastErr = fmt.Errorf("no Jakarta price rows found")
	}
	return nil, lastErr
}

func parseGenericFuelPrices(source Source, raw string, allowWithoutJakarta bool) []FuelPrice {
	text := normalizeHTML(raw)
	effectiveAt := parseEffectiveDate(text)
	scrapedAt := time.Now().UTC().Format(time.RFC3339)
	prices := make([]FuelPrice, 0)

	for _, fuelName := range fuelCandidates(source.Brand) {
		price := parsePriceNearFuel(text, fuelName, allowWithoutJakarta)
		if price == 0 {
			continue
		}
		prices = append(prices, FuelPrice{
			Brand:        source.Brand,
			FuelName:     fuelName,
			FuelType:     inferFuelType(fuelName),
			Price:        price,
			Location:     "DKI Jakarta",
			EffectiveAt:  effectiveAt,
			SourceURL:    source.URL,
			SourceStatus: "scraped",
			ScrapedAt:    scrapedAt,
		})
	}

	return dedupePrices(prices)
}

func fuelCandidates(brand string) []string {
	candidates := map[string][]string{
		"Shell": {"Shell V-Power Nitro+", "Shell V-Power Diesel", "Shell V-Power", "Shell Super"},
		"BP":    {"BP Ultimate Diesel", "BP Ultimate", "BP 92", "Regular 92"},
		"Vivo":  {"Revvo 95", "Revvo 92", "Revvo 90"},
	}
	return candidates[brand]
}

func parsePriceNearFuel(text, fuelName string, allowWithoutJakarta bool) int {
	lowerText := strings.ToLower(text)
	lowerFuel := strings.ToLower(fuelName)
	start := 0
	for {
		index := strings.Index(lowerText[start:], lowerFuel)
		if index < 0 {
			return 0
		}
		absolute := start + index
		if isAmbiguousFuelMatch(lowerText, absolute, lowerFuel) {
			start = absolute + len(lowerFuel)
			continue
		}
		end := absolute + 420
		if end > len(text) {
			end = len(text)
		}
		window := text[absolute:end]
		if allowWithoutJakarta || mentionsJakarta(window) {
			price := parseRupiah(window)
			if price > 0 {
				return price
			}
		}
		start = absolute + len(lowerFuel)
	}
}

func isAmbiguousFuelMatch(lowerText string, absolute int, lowerFuel string) bool {
	if lowerFuel != "shell v-power" && lowerFuel != "bp ultimate" {
		return false
	}
	afterStart := absolute + len(lowerFuel)
	afterEnd := afterStart + 18
	if afterEnd > len(lowerText) {
		afterEnd = len(lowerText)
	}
	after := lowerText[afterStart:afterEnd]
	return strings.HasPrefix(after, " diesel") || strings.HasPrefix(after, " nitro")
}
