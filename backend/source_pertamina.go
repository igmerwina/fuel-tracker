package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func scrapePertamina(ctx context.Context, client *http.Client, source Source) ([]FuelPrice, error) {
	body, err := fetchBody(ctx, client, source.URL)
	if err != nil {
		return nil, err
	}
	prices := parsePertaminaPatraNiagaPrices(source, body)
	if len(prices) == 0 {
		return nil, fmt.Errorf("no Prov. DKI Jakarta Gasoline/GasOil rows found")
	}
	return prices, nil
}

func parsePertaminaPatraNiagaPrices(source Source, body []byte) []FuelPrice {
	var payload struct {
		Data struct {
			Title   string `json:"title"`
			Content map[string]struct {
				Props map[string]any `json:"props"`
			} `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil
	}

	effectiveAt := payload.Data.Title
	if heading := findPertaminaHeading(payload.Data.Content); heading != "" {
		effectiveAt = heading
	}

	scrapedAt := time.Now().UTC().Format(time.RFC3339)
	prices := []FuelPrice{}
	for _, node := range payload.Data.Content {
		rawItems, ok := node.Props["items"].([]any)
		if !ok {
			continue
		}
		for _, rawItem := range rawItems {
			item, ok := rawItem.(map[string]any)
			if !ok || !isPertaminaFuelTable(fmt.Sprint(item["title"])) {
				continue
			}
			prices = append(prices, pertaminaPricesFromRows(source, item, effectiveAt, scrapedAt)...)
		}
	}
	return dedupePrices(prices)
}

func isPertaminaFuelTable(title string) bool {
	return strings.EqualFold(title, "Gasoline") || strings.EqualFold(title, "Gasoil")
}

func pertaminaPricesFromRows(source Source, item map[string]any, effectiveAt string, scrapedAt string) []FuelPrice {
	rawRows, ok := item["data"].([]any)
	if !ok {
		return nil
	}

	prices := []FuelPrice{}
	for _, rawRow := range rawRows {
		row, ok := rawRow.(map[string]any)
		if !ok || fmt.Sprint(row["REGION"]) != "Prov. DKI Jakarta" {
			continue
		}
		for _, product := range pertaminaProductOrder() {
			rawValue, ok := valueByAssetNeedle(row, product.AssetNeedle)
			price := parsePlainPrice(fmt.Sprint(rawValue))
			if !ok || price == 0 {
				continue
			}
			prices = append(prices, FuelPrice{
				Brand:        source.Brand,
				FuelName:     product.Name,
				FuelType:     inferFuelType(product.Name),
				Price:        price,
				Location:     "Prov. DKI Jakarta",
				EffectiveAt:  effectiveAt,
				SourceURL:    source.URL,
				SourceStatus: "scraped",
				ScrapedAt:    scrapedAt,
			})
		}
	}
	return prices
}

type pertaminaProduct struct {
	Name        string
	AssetNeedle string
}

func pertaminaProductOrder() []pertaminaProduct {
	return []pertaminaProduct{
		{Name: "Pertalite", AssetNeedle: "pertalite"},
		{Name: "Pertamax", AssetNeedle: "product-table-pertamax.png"},
		{Name: "Pertamax Green 95", AssetNeedle: "pertamax-green-95"},
		{Name: "Pertamax Turbo", AssetNeedle: "pertamax-turbo"},
		{Name: "Pertamina Dex", AssetNeedle: "pertamina-dex"},
		{Name: "Dexlite", AssetNeedle: "dexlite"},
		{Name: "Bio Solar Subsidi", AssetNeedle: "bio-solar-subsidi"},
		{Name: "Pertamax Pertashop", AssetNeedle: "pertamax-pertashop"},
	}
}

func valueByAssetNeedle(row map[string]any, needle string) (any, bool) {
	for key, value := range row {
		if strings.Contains(strings.ToLower(key), needle) {
			return value, true
		}
	}
	return nil, false
}

func findPertaminaHeading(content map[string]struct {
	Props map[string]any `json:"props"`
}) string {
	for _, node := range content {
		text, ok := node.Props["text"].(string)
		if ok && strings.Contains(strings.ToLower(text), "update as of") {
			return text
		}
	}
	return ""
}

func pertaminaFuelNameFromAsset(assetURL string) string {
	lower := strings.ToLower(assetURL)
	switch {
	case strings.Contains(lower, "pertamax-turbo"):
		return "Pertamax Turbo"
	case strings.Contains(lower, "pertamax-green-95"):
		return "Pertamax Green 95"
	case strings.Contains(lower, "pertamax-pertashop"):
		return "Pertamax Pertashop"
	case strings.Contains(lower, "product-table-pertamax"):
		return "Pertamax"
	case strings.Contains(lower, "pertalite"):
		return "Pertalite"
	case strings.Contains(lower, "pertamina-dex"):
		return "Pertamina Dex"
	case strings.Contains(lower, "dexlite"):
		return "Dexlite"
	case strings.Contains(lower, "bio-solar-non-subsidi"):
		return "Bio Solar Non Subsidi"
	case strings.Contains(lower, "bio-solar-subsidi"):
		return "Bio Solar Subsidi"
	default:
		return ""
	}
}
