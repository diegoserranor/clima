package search

import (
	"fmt"

	"github.com/esferadigital/clima/internal/openmeteo"
)

// Implements list.Item interface and wraps openmeteo.GeocodingResult
type searchListItem struct {
	openmeteo.GeocodingResult
}

func (i searchListItem) FilterValue() string {
	return i.Name
}

func (i searchListItem) Title() string {
	return fmt.Sprintf("%s, %s", i.Name, i.Country)
}

func (i searchListItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", i.Latitude, i.Longitude)
}
