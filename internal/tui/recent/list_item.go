package recent

import (
	"fmt"

	"github.com/diegoserranor/clima/internal/openmeteo"
)

// Implements list.Item interface and wraps location.Location
type recentLocationItem struct {
	openmeteo.GeocodingResult
}

func (i recentLocationItem) FilterValue() string {
	return i.Name
}

func (i recentLocationItem) Title() string {
	place := i.Name
	if i.Admin1 != "" {
		place = place + ", " + i.Admin1
	}
	if i.Country != "" {
		place = place + ", " + i.Country
	}
	return place
}

func (i recentLocationItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", i.Latitude, i.Longitude)
}
