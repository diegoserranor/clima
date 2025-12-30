package search

import "github.com/diegoserranor/clima/internal/openmeteo"

type dataMsg struct {
	locations []openmeteo.GeocodingResult
}

type errorMsg struct {
	err error
}

type SearchCompleteMsg struct {
	Location openmeteo.GeocodingResult
}

type RecentMsg struct{}
