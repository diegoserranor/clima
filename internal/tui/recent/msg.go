package recent

import "github.com/diegoserranor/clima/internal/openmeteo"

type dataMsg struct {
	locations []openmeteo.GeocodingResult
}

type errorMsg struct {
	err error
}

type RecentCompleteMsg struct {
	Location openmeteo.GeocodingResult
	OK       bool
}

type NewSearchMsg struct{}
