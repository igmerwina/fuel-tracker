package scraper

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

var fuelTypeByName = []struct {
	Needles []string
	Type    string
}{
	{[]string{"dexlite", "v-power diesel", "ultimate diesel", "bio solar"}, "DIESEL_CN_51"},
	{[]string{"pertamina dex"}, "DIESEL_CN_53"},
	{[]string{"pertalite", "revvo 90"}, "RON_90"},
	{[]string{"pertamax green", "revvo 95", "v-power", "bp ultimate"}, "RON_95"},
	{[]string{"pertamax turbo", "nitro"}, "RON_98"},
	{[]string{"pertamax", "shell super", "bp 92", "revvo 92", "regular 92"}, "RON_92"},
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
	cleaned = strings.ReplaceAll(cleaned, "</td>", " | ")
	cleaned = strings.ReplaceAll(cleaned, "</li>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "<br>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "<br/>", "\n")
	cleaned = strings.ReplaceAll(cleaned, "<br />", "\n")
	cleaned = reTag.ReplaceAllString(cleaned, " ")
	cleaned = html.UnescapeString(cleaned)
	cleaned = reSpace.ReplaceAllString(cleaned, " ")
	return cleaned
}

func parseRupiah(line string) int {
	re := regexp.MustCompile(`(?i)(?:rp|idr)\s*([0-9]{1,3}(?:[.,][0-9]{3})+|[0-9]+)`)
	match := re.FindStringSubmatch(line)
	if len(match) < 2 {
		re = regexp.MustCompile(`\b([0-9]{1,3}(?:[.,][0-9]{3})+)\b`)
		match = re.FindStringSubmatch(line)
	}
	if len(match) < 2 {
		return 0
	}
	return parsePlainPrice(match[1])
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
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(\d{1,2}\s+(januari|februari|maret|april|mei|juni|juli|agustus|september|oktober|november|desember)\s+\d{4})\b`),
		regexp.MustCompile(`(?i)(harga berlaku efektif\s+\d{1,2}\s+\w+\s+\d{4})`),
	}
	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(text)
		if len(match) > 1 {
			return strings.TrimSpace(match[1])
		}
	}
	return ""
}

func mentionsJakarta(line string) bool {
	lower := strings.ToLower(line)
	return strings.Contains(lower, "jakarta") ||
		strings.Contains(lower, "dki") ||
		strings.Contains(lower, "jabodetabek")
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
