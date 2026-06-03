package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func scrapeBP(ctx context.Context, client *http.Client, source Source) ([]FuelPrice, error) {
	body, err := fetchBody(ctx, client, source.URL)
	if err != nil {
		return nil, err
	}
	prices := parseBPPrices(source, string(body))
	if len(prices) == 0 {
		return nil, fmt.Errorf("no Jabodetabek price rows found")
	}
	return prices, nil
}

func parseBPPrices(source Source, raw string) []FuelPrice {
	effectiveAt := parseEffectiveDate(normalizeHTML(raw))
	scrapedAt := time.Now().UTC().Format(time.RFC3339)
	rows := bpProductRows(raw)
	prices := []FuelPrice{}

	for fuelName, price := range rows {
		if price == 0 {
			continue
		}
		prices = append(prices, FuelPrice{
			Brand:        source.Brand,
			FuelName:     fuelName,
			FuelType:     inferFuelType(fuelName),
			Price:        price,
			Location:     "Jabodetabek",
			EffectiveAt:  effectiveAt,
			SourceURL:    source.URL,
			SourceStatus: "scraped",
			ScrapedAt:    scrapedAt,
		})
	}
	return dedupePrices(prices)
}

func bpProductRows(raw string) map[string]int {
	result := map[string]int{}
	rowPattern := regexp.MustCompile(`(?is)<tr>\s*<td[^>]*>\s*<b>(BP Ultimate Diesel|BP Ultimate|BP 92|Regular 92)</b>\s*</td>\s*<td[^>]*>\s*([^<]+)\s*</td>`)
	for _, match := range rowPattern.FindAllStringSubmatch(raw, -1) {
		if len(match) < 3 {
			continue
		}
		fuelName := strings.TrimSpace(match[1])
		result[fuelName] = parseRupiah(match[2])
	}
	return result
}
