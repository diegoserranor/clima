package recent

import (
	"fmt"

	"github.com/esferadigital/clima/internal/openmeteo"
)

// Implements list.Item interface and wraps location.Location
type recentLocationItem struct {
	openmeteo.GeocodingResult
}

func (i recentLocationItem) FilterValue() string {
	return i.Name
}

func (i recentLocationItem) Title() string {
	return fmt.Sprintf("%s, %s", i.Name, i.Country)
}

func (i recentLocationItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", i.Latitude, i.Longitude)
}
