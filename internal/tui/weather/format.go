package weather

import (
	"fmt"
	"time"

	"github.com/esferadigital/clima/internal/openmeteo"
)

func formatMeasurement(measurement openmeteo.FloatMeasurement) string {
	return formatValueWithUnit(measurement.Value, measurement.Unit)
}

func formatValueWithUnit(value float64, unit string) string {
	if unit == "" {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

func formatHourlyTime(raw string, offsetSeconds int, timezone string) string {
	if raw == "" {
		return "-"
	}
	layouts := []string{time.RFC3339, "2006-01-02T15:04"}
	var parsed time.Time
	var err error
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, raw)
		if err == nil {
			break
		}
	}
	if err != nil {
		return raw
	}
	if timezone == "" {
		timezone = "UTC"
	}
	loc := time.FixedZone(timezone, offsetSeconds)
	return parsed.In(loc).Format("15:04")
}

func formatDailyDate(raw string, offsetSeconds int, timezone string) string {
	if raw == "" {
		return "-"
	}
	layouts := []string{"2006-01-02", time.RFC3339}
	var parsed time.Time
	var err error
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, raw)
		if err == nil {
			break
		}
	}
	if err != nil {
		return raw
	}
	if timezone == "" {
		timezone = "UTC"
	}
	loc := time.FixedZone(timezone, offsetSeconds)
	return parsed.In(loc).Format("Mon 02")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
