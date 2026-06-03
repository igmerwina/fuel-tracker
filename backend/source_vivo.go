package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func scrapeVivo(ctx context.Context, client *http.Client, source Source) (SourceResult, error) {
	body, err := fetchBody(ctx, client, source.URL)
	if err != nil {
		return SourceResult{}, err
	}
	stations := parseVivoStations(string(body))
	if len(stations) == 0 {
		return SourceResult{}, fmt.Errorf("no DKI Jakarta Vivo stock rows found")
	}

	return SourceResult{Stations: vivoStationLocations(source, stations)}, nil
}

type VivoStation struct {
	Name          string `json:"name"`
	Area          string `json:"area"`
	LastUpdated   string `json:"last_updated"`
	RON92Status   string `json:"ron_92_status"`
	RON95Status   string `json:"ron_95_status"`
	DieselStatus  string `json:"diesel_status"`
	GoogleMapsURL string `json:"google_maps_url"`
}

func parseVivoStations(raw string) []VivoStation {
	rowPattern := regexp.MustCompile(`(?is)<tr\s+data-area="DKI Jakarta">(.*?)</tr>`)
	cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
	linkPattern := regexp.MustCompile(`(?is)<a\s+href="([^"]+)"`)
	stations := []VivoStation{}

	for _, rowMatch := range rowPattern.FindAllStringSubmatch(raw, -1) {
		cells := cellPattern.FindAllStringSubmatch(rowMatch[1], -1)
		if len(cells) < 7 {
			continue
		}
		link := ""
		if linkMatch := linkPattern.FindStringSubmatch(cells[6][1]); len(linkMatch) > 1 {
			link = strings.TrimSpace(linkMatch[1])
		}
		stations = append(stations, VivoStation{
			Name:          cleanCell(cells[0][1]),
			Area:          cleanCell(cells[1][1]),
			LastUpdated:   cleanCell(cells[2][1]),
			RON92Status:   cleanCell(cells[3][1]),
			RON95Status:   cleanCell(cells[4][1]),
			DieselStatus:  cleanCell(cells[5][1]),
			GoogleMapsURL: link,
		})
	}
	return stations
}

func vivoStationLocations(source Source, stations []VivoStation) []StationLocation {
	scrapedAt := time.Now().UTC().Format(time.RFC3339)
	locations := make([]StationLocation, 0, len(stations))
	for _, station := range stations {
		locations = append(locations, StationLocation{
			Brand:         source.Brand,
			Name:          station.Name,
			Area:          station.Area,
			GoogleMapsURL: station.GoogleMapsURL,
			LastUpdated:   station.LastUpdated,
			Status: map[string]string{
				"RON_92": station.RON92Status,
				"RON_95": station.RON95Status,
				"DIESEL": station.DieselStatus,
			},
			SourceURL: source.URL,
			ScrapedAt: scrapedAt,
		})
	}
	return locations
}

func cleanCell(value string) string {
	return strings.TrimSpace(normalizeHTML(value))
}
