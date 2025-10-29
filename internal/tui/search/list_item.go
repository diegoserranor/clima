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
	place := i.Name
	if i.Admin1 != "" {
		place = place + ", " + i.Admin1
	}
	if i.Country != "" {
		place = place + ", " + i.Country
	}
	return place
}

func (i searchListItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", i.Latitude, i.Longitude)
}
