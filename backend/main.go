package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FuelPrice struct {
	Brand        string `json:"brand"`
	FuelName     string `json:"fuel_name"`
	FuelType     string `json:"fuel_type"`
	Price        int    `json:"price"`
	Location     string `json:"location"`
	EffectiveAt  string `json:"effective_at,omitempty"`
	SourceURL    string `json:"source_url"`
	SourceStatus string `json:"source_status"`
	ScrapedAt    string `json:"scraped_at"`
}

type PriceCache struct {
	GeneratedAt string      `json:"generated_at"`
	Regions     []string    `json:"regions"`
	Prices      []FuelPrice `json:"prices"`
	Errors      []string    `json:"errors"`
}

type Source struct {
	Brand string
	URL   string
}

var jakartaRegions = []string{
	"Jakarta Pusat",
	"Jakarta Selatan",
	"Jakarta Barat",
	"Jakarta Timur",
	"Jakarta Utara",
}

var fuelTypeByName = []struct {
	Needles []string
	Type    string
}{
	{[]string{"dexlite", "v-power diesel", "ultimate diesel"}, "DIESEL_CN_51"},
	{[]string{"pertamina dex"}, "DIESEL_CN_53"},
	{[]string{"pertalite", "revvo 90"}, "RON_90"},
	{[]string{"pertamax green", "revvo 95", "v-power", "bp ultimate"}, "RON_95"},
	{[]string{"pertamax turbo", "nitro"}, "RON_98"},
	{[]string{"pertamax", "shell super", "bp 92", "revvo 92", "regular 92"}, "RON_92"},
}

func main() {
	port := getenv("PORT", "8080")
	cachePath := filepath.Join("data", "fuel_prices.json")
	client := &http.Client{Timeout: scrapeTimeout()}

	cache, err := scrapeAndStore(context.Background(), client, cachePath)
	if err != nil {
		log.Printf("initial scrape completed with warnings: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/v1/prices", func(w http.ResponseWriter, r *http.Request) {
		current, readErr := readCache(cachePath)
		if readErr == nil {
			cache = current
		}
		writeJSON(w, http.StatusOK, cache)
	})
	mux.HandleFunc("POST /api/v1/scrape", func(w http.ResponseWriter, r *http.Request) {
		next, scrapeErr := scrapeAndStore(r.Context(), client, cachePath)
		cache = next
		if scrapeErr != nil {
			writeJSON(w, http.StatusAccepted, cache)
			return
		}
		writeJSON(w, http.StatusOK, cache)
	})

	log.Printf("backend listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, withCORS(mux)))
}

func scrapeAndStore(ctx context.Context, client *http.Client, cachePath string) (PriceCache, error) {
	cache := PriceCache{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Regions:     jakartaRegions,
		Prices:      []FuelPrice{},
		Errors:      []string{},
	}

	for _, source := range sources() {
		prices, err := scrapeSource(ctx, client, source)
		if err != nil {
			cache.Errors = append(cache.Errors, fmt.Sprintf("%s: %v", source.Brand, err))
			continue
		}
		cache.Prices = append(cache.Prices, prices...)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return cache, err
	}
	payload, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return cache, err
	}
	if err := os.WriteFile(cachePath, payload, 0o644); err != nil {
		return cache, err
	}

	if len(cache.Errors) > 0 {
		return cache, errors.New(strings.Join(cache.Errors, "; "))
	}
	return cache, nil
}

func sources() []Source {
	return []Source{
		{
			Brand: "Pertamina",
			URL: getenv(
				"PERTAMINA_PRICE_URL",
				"https://pertaminapatraniaga.com/api/api/v1/post/get-by-slug/page/harga-terbaru-bbm?language=id",
			),
		},
		{
			Brand: "Shell",
			URL: getenv(
				"SHELL_PRICE_URL",
				"https://www.shell.co.id/in_id/pengendara-bermotor/bahan-bakar-shell/harga-bahan-bakar-shell.html",
			),
		},
		{
			Brand: "BP",
			URL: getenv(
				"BP_PRICE_URL",
				"https://www.bp.com/id_id/indonesia/home/produk-dan-layanan/spbu.html",
			),
		},
		{
			Brand: "Vivo",
			URL:   getenv("VIVO_PRICE_URL", "https://www.vivoenergy.co.id/"),
		},
	}
}

func scrapeSource(ctx context.Context, client *http.Client, source Source) ([]FuelPrice, error) {
	urls := []string{source.URL}
	if strings.HasSuffix(source.URL, ".html") {
		urls = append(urls, strings.TrimSuffix(source.URL, ".html")+".model.json")
	}

	var lastErr error
	for _, candidateURL := range urls {
		nextSource := source
		nextSource.URL = candidateURL
		prices, err := fetchAndParseSource(ctx, client, nextSource)
		if err == nil {
			return prices, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func fetchAndParseSource(ctx context.Context, client *http.Client, source Source) ([]FuelPrice, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "JakartaBBMTracker/0.1 (+https://localhost)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("official page returned HTTP %d", res.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, 4<<20))
	if err != nil {
		return nil, err
	}

	if source.Brand == "Pertamina" {
		prices := parsePertaminaPatraNiagaPrices(source, body)
		if len(prices) == 0 {
			return nil, errors.New("no Prov. DKI Jakarta Gasoline/GasOil rows found")
		}
		return prices, nil
	}

	prices := parsePrices(source, string(body))
	if len(prices) == 0 {
		return nil, errors.New("no Jakarta price rows found")
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
			if !ok {
				continue
			}
			title := fmt.Sprint(item["title"])
			if !strings.EqualFold(title, "Gasoline") && !strings.EqualFold(title, "Gasoil") {
				continue
			}
			rawRows, ok := item["data"].([]any)
			if !ok {
				continue
			}
			for _, rawRow := range rawRows {
				row, ok := rawRow.(map[string]any)
				if !ok || fmt.Sprint(row["REGION"]) != "Prov. DKI Jakarta" {
					continue
				}
				for key, rawValue := range row {
					if key == "REGION" {
						continue
					}
					fuelName := pertaminaFuelNameFromAsset(key)
					if fuelName == "" {
						continue
					}
					price := parsePlainPrice(fmt.Sprint(rawValue))
					if price == 0 {
						continue
					}
					prices = append(prices, FuelPrice{
						Brand:        source.Brand,
						FuelName:     fuelName,
						FuelType:     inferFuelType(fuelName),
						Price:        price,
						Location:     "Prov. DKI Jakarta",
						EffectiveAt:  effectiveAt,
						SourceURL:    source.URL,
						SourceStatus: "scraped",
						ScrapedAt:    scrapedAt,
					})
				}
			}
		}
	}
	return dedupePrices(prices)
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

func parsePrices(source Source, raw string) []FuelPrice {
	text := normalizeHTML(raw)
	effectiveAt := parseEffectiveDate(text)
	scrapedAt := time.Now().UTC().Format(time.RFC3339)
	prices := make([]FuelPrice, 0)

	for _, fuelName := range fuelCandidates(source.Brand) {
		price := parsePriceNearFuel(text, fuelName, source.Brand == "Vivo")
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

func normalizeHTML(raw string) string {
	raw = strings.NewReplacer(
		`\\u003C`, "<",
		`\\u003E`, ">",
		`\\u0026`, "&",
		`\\r\\n`, "\n",
		`\\n`, "\n",
		`\\t`, " ",
		`\u003C`, "<",
		`\u003E`, ">",
		`\u0026`, "&",
		`\r\n`, "\n",
		`\n`, "\n",
		`\t`, " ",
		`\u00a0`, " ",
	).Replace(raw)
	reScript := regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	reTag := regexp.MustCompile(`(?s)<[^>]+>`)
	reSpace := regexp.MustCompile(`[ \t\r\f\v]+`)
	cleaned := reScript.ReplaceAllString(raw, " ")
	cleaned = strings.ReplaceAll(cleaned, "</tr>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "</li>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "<br>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "<br/>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "<br />", "\n")
	cleaned = reTag.ReplaceAllString(cleaned, " ")
	cleaned = html.UnescapeString(cleaned)
	cleaned = reSpace.ReplaceAllString(cleaned, " ")
	return cleaned
}

func splitCandidateLines(text string) []string {
	rawLines := strings.Split(text, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(strings.ToLower(line), "rp") || regexp.MustCompile(`\d{1,3}([.]\d{3})+`).MatchString(line) {
			lines = append(lines, line)
		}
	}
	if len(lines) > 0 {
		return lines
	}
	return []string{text}
}

func fuelCandidates(brand string) []string {
	candidates := map[string][]string{
		"Pertamina": {"Pertamax Green 95", "Pertamax Turbo", "Pertamina Dex", "Pertamax", "Dexlite", "Pertalite"},
		"Shell":     {"Shell V-Power Nitro+", "Shell V-Power Diesel", "Shell V-Power", "Shell Super"},
		"BP":        {"BP Ultimate Diesel", "BP Ultimate", "BP 92", "Regular 92"},
		"Vivo":      {"Revvo 95", "Revvo 92", "Revvo 90"},
	}
	return candidates[brand]
}

func parseFuelName(brand, line string) string {
	lower := strings.ToLower(line)
	for _, candidate := range fuelCandidates(brand) {
		if strings.Contains(lower, strings.ToLower(candidate)) {
			return candidate
		}
	}
	return ""
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
	if lowerFuel != "shell v-power" && lowerFuel != "bp ultimate" && lowerFuel != "pertamax" {
		return false
	}
	afterStart := absolute + len(lowerFuel)
	afterEnd := afterStart + 18
	if afterEnd > len(lowerText) {
		afterEnd = len(lowerText)
	}
	after := lowerText[afterStart:afterEnd]
	return strings.HasPrefix(after, " diesel") ||
		strings.HasPrefix(after, " nitro") ||
		strings.HasPrefix(after, " turbo") ||
		strings.HasPrefix(after, " green")
}

func parseRupiah(line string) int {
	re := regexp.MustCompile(`(?i)rp\s*([0-9]{1,3}(?:[.,][0-9]{3})+|[0-9]+)`)
	match := re.FindStringSubmatch(line)
	if len(match) < 2 {
		re = regexp.MustCompile(`\b([0-9]{1,3}(?:[.][0-9]{3})+)\b`)
		match = re.FindStringSubmatch(line)
	}
	if len(match) < 2 {
		return 0
	}
	value := strings.NewReplacer(".", "", ",", "").Replace(match[1])
	price, _ := strconv.Atoi(value)
	return price
}

func parsePlainPrice(value string) int {
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(value, "-") {
		return 0
	}
	value = strings.NewReplacer(".", "", ",", "").Replace(value)
	price, _ := strconv.Atoi(value)
	return price
}

func parseEffectiveDate(text string) string {
	re := regexp.MustCompile(`(?i)\b(\d{1,2}\s+(januari|februari|maret|april|mei|juni|juli|agustus|september|oktober|november|desember)\s+\d{4})\b`)
	match := re.FindStringSubmatch(text)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}

func mentionsJakarta(line string) bool {
	lower := strings.ToLower(line)
	return strings.Contains(lower, "jakarta") || strings.Contains(lower, "dki")
}

func inferFuelType(fuelName string) string {
	lower := strings.ToLower(fuelName)
	for _, mapping := range fuelTypeByName {
		for _, needle := range mapping.Needles {
			if strings.Contains(lower, needle) {
				return mapping.Type
			}
		}
	}
	return "UNKNOWN"
}

func dedupePrices(prices []FuelPrice) []FuelPrice {
	seen := map[string]bool{}
	result := make([]FuelPrice, 0, len(prices))
	for _, price := range prices {
		key := price.Brand + "|" + price.FuelName + "|" + strconv.Itoa(price.Price)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, price)
	}
	return result
}

func readCache(path string) (PriceCache, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return PriceCache{}, err
	}
	var cache PriceCache
	if err := json.Unmarshal(payload, &cache); err != nil {
		return PriceCache{}, err
	}
	return cache, nil
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

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func scrapeTimeout() time.Duration {
	raw := getenv("SCRAPE_TIMEOUT_SECONDS", "20")
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		seconds = 20
	}
	return time.Duration(seconds) * time.Second
}
