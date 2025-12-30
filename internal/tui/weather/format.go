package weather

import (
	"fmt"
	"time"

	"github.com/diegoserranor/clima/internal/openmeteo"
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

func formatHourlyTime(raw string) string {
	if raw == "" {
		return "-"
	}
	const inputLayout = "2006-01-02T15:04"
	t, err := time.Parse(inputLayout, raw)
	if err != nil {
		return raw
	}
	return t.Format("3 PM")
}

func formatDailyDate(raw string) string {
	if raw == "" {
		return "-"
	}
	const inputLayout = "2006-01-02"
	t, err := time.Parse(inputLayout, raw)
	if err != nil {
		return raw
	}
	return t.Format("Mon 2")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
