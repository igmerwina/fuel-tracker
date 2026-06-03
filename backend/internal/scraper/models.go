package scraper

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
	GeneratedAt string            `json:"generated_at"`
	Regions     []string          `json:"regions"`
	Prices      []FuelPrice       `json:"prices"`
	Stations    []StationLocation `json:"stations"`
	Errors      []string          `json:"errors"`
}

type Source struct {
	Brand string
	URL   string
}

type StationLocation struct {
	Brand         string            `json:"brand"`
	Name          string            `json:"name"`
	Area          string            `json:"area"`
	Address       string            `json:"address,omitempty"`
	Latitude      float64           `json:"latitude,omitempty"`
	Longitude     float64           `json:"longitude,omitempty"`
	GoogleMapsURL string            `json:"google_maps_url"`
	LastUpdated   string            `json:"last_updated,omitempty"`
	Status        map[string]string `json:"status,omitempty"`
	SourceURL     string            `json:"source_url"`
	ScrapedAt     string            `json:"scraped_at"`
}

var jakartaRegions = []string{
	"Jakarta Pusat",
	"Jakarta Selatan",
	"Jakarta Barat",
	"Jakarta Timur",
	"Jakarta Utara",
}
