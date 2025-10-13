package weather

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/theme"
)

var (
	blockSpacing    = lipgloss.NewStyle().MarginBottom(1)
	sectionDivider  = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderBottomForeground(lipgloss.Color("13"))
	hourColumnStyle = lipgloss.NewStyle().Width(16).Padding(0, 1)
	dayColumnStyle  = lipgloss.NewStyle().Width(16).Padding(0, 1)
)

type forecastColumn struct {
	heading     string
	description string
	details     []string
}

func renderError() string {
	return "error"
}

func renderLoading(ellipsis spinner.Model) string {
	return fmt.Sprintf("Loading forecast%s\n", ellipsis.View())
}

func renderHeader(location openmeteo.GeocodingResult, forecast openmeteo.ForecastResponse) string {
	place := location.Name + ", " + location.Country
	if weatherCode, ok := forecast.CurrentMeasurement(openmeteo.CurrentWeatherCode); ok {
		conditions := openmeteo.MapWeatherCode(weatherCode.Value)
		conditions = theme.Accent.Render(conditions)
		return lipgloss.JoinVertical(lipgloss.Left, place, conditions)
	}
	return place
}

func renderCurrent(forecast openmeteo.ForecastResponse) string {
	weather := forecast
	s := ""

	if currentTemp, ok := weather.CurrentMeasurement(openmeteo.CurrentTemperature2m); ok {
		temperature := formatMeasurement(currentTemp)
		if currentApparentTemp, ok := weather.CurrentMeasurement(openmeteo.CurrentApparentTemperature); ok {
			feelsLike := fmt.Sprintf(" (feels like %s)", formatMeasurement(currentApparentTemp))
			temperature += feelsLike
		}
		s += temperature
	}

	minTempLabel := theme.Label.Render("\nMin")
	var minTempValue string
	if minSeries, ok := weather.DailySeries(openmeteo.DailyTemperature2mMin); ok && len(minSeries.Values) > 0 {
		minTempValue = formatValueWithUnit(minSeries.Values[0], minSeries.Unit)
	} else {
		minTempValue = "-"
	}
	s += minTempLabel + minTempValue

	maxTempLabel := theme.Label.Render("\nMax")
	var maxTempValue string
	if maxSeries, ok := weather.DailySeries(openmeteo.DailyTemperature2mMax); ok && len(maxSeries.Values) > 0 {
		maxTempValue = formatValueWithUnit(maxSeries.Values[0], maxSeries.Unit)
	} else {
		maxTempValue = "-"
	}
	s += maxTempLabel + maxTempValue

	windLabel := theme.Label.Render("\n\nWind")
	windValue := "-"
	if windSpeed, ok := weather.CurrentMeasurement(openmeteo.CurrentWindSpeed10m); ok {
		windValue = formatMeasurement(windSpeed)
		if windDirection, ok := weather.CurrentMeasurement(openmeteo.CurrentWindDirection10m); ok {
			windValue = fmt.Sprintf("%.1f %s @ %.1f %s", windSpeed.Value, windSpeed.Unit, windDirection.Value, windDirection.Unit)
		}
	}
	s += windLabel + windValue

	windGustsLabel := theme.Label.Render("\nWind gusts")
	windGustsValue := "-"
	if windGusts, ok := weather.CurrentMeasurement(openmeteo.CurrentWindGusts10m); ok {
		windGustsValue = formatMeasurement(windGusts)
	}
	s += windGustsLabel + windGustsValue

	humidityLabel := theme.Label.Render("\nHumidity")
	humidityValue := "-"
	if humidity, ok := weather.CurrentMeasurement(openmeteo.CurrentRelativeHumidity2m); ok {
		humidityValue = formatMeasurement(humidity)
	}
	s += humidityLabel + humidityValue

	precipitationLabel := theme.Label.Render("\nPrecipitation")
	precipitationValue := "-"
	if precipitation, ok := weather.CurrentMeasurement(openmeteo.CurrentPrecipitation); ok {
		precipitationValue = formatMeasurement(precipitation)
	}
	s += precipitationLabel + precipitationValue

	pressureLabel := theme.Label.Render("\nPressure")
	pressureValue := "-"
	if pressure, ok := weather.CurrentMeasurement(openmeteo.CurrentSeaLevelPressure); ok {
		pressureValue = formatMeasurement(pressure)
	}
	s += pressureLabel + pressureValue

	uvLabel := theme.Label.Render("\nUV index")
	var uvValue string
	if uvSeries, ok := weather.DailySeries(openmeteo.DailyUVIndexMax); ok && len(uvSeries.Values) > 0 {
		uvValue = fmt.Sprintf("%.1f", uvSeries.Values[0])
	} else {
		uvValue = "-"
	}
	s += uvLabel + uvValue
	return s + "\n"
}

func renderHourly(forecast openmeteo.ForecastResponse) string {
	columns := buildHourlyColumns(forecast)
	if len(columns) == 0 {
		return theme.Subtle.Render("Hourly forecast unavailable")
	}
	return renderColumns(hourColumnStyle, columns) + "\n"
}

func buildHourlyColumns(forecast openmeteo.ForecastResponse) []forecastColumn {
	if len(forecast.HourlyTimes) == 0 {
		return nil
	}

	weatherCodes, hasWeatherCodes := forecast.HourlySeries(openmeteo.HourlyWeatherCode)
	temperatures, hasTemperatures := forecast.HourlySeries(openmeteo.HourlyTemperature2m)
	precipitation, hasPrecipitation := forecast.HourlySeries(openmeteo.HourlyPrecipitation)

	hours := min(5, len(forecast.HourlyTimes))
	columns := make([]forecastColumn, 0, hours)

	for i := range hours {
		label := theme.Subtle.Render(formatHourlyTime(forecast.HourlyTimes[i], forecast.UTCOffsetSeconds, forecast.Timezone))

		description := "-"
		if hasWeatherCodes && len(weatherCodes.Values) > i {
			if mapped := openmeteo.MapWeatherCode(weatherCodes.Values[i]); mapped != "" {
				description = theme.Accent.Render(mapped)
			}
		}

		forecastLine := "-"
		if hasTemperatures && len(temperatures.Values) > i {
			forecastLine = formatValueWithUnit(temperatures.Values[i], temperatures.Unit)
		}
		if hasPrecipitation && len(precipitation.Values) > i && precipitation.Values[i] > 0 {
			precip := formatValueWithUnit(precipitation.Values[i], precipitation.Unit)
			if forecastLine == "-" {
				forecastLine = precip
			} else {
				forecastLine = fmt.Sprintf("%s, %s", forecastLine, precip)
			}
		}

		columns = append(columns, forecastColumn{
			heading:     label,
			description: description,
			details:     []string{forecastLine},
		})
	}

	return columns
}

func renderDaily(forecast openmeteo.ForecastResponse) string {
	columns := buildDailyColumns(forecast)
	if len(columns) == 0 {
		return theme.Subtle.Render("Daily outlook unavailable")
	}
	return renderColumns(dayColumnStyle, columns)
}

func buildDailyColumns(forecast openmeteo.ForecastResponse) []forecastColumn {
	if len(forecast.DailyTimes) == 0 {
		return nil
	}

	minTemps, hasMin := forecast.DailySeries(openmeteo.DailyTemperature2mMin)
	maxTemps, hasMax := forecast.DailySeries(openmeteo.DailyTemperature2mMax)
	weatherCodes, hasCodes := forecast.DailySeries(openmeteo.DailyWeatherCode)

	days := min(5, len(forecast.DailyTimes))
	columns := make([]forecastColumn, 0, days)

	for i := range days {
		label := theme.Subtle.Render(formatDailyDate(forecast.DailyTimes[i], forecast.UTCOffsetSeconds, forecast.Timezone))

		description := "-"
		if hasCodes && len(weatherCodes.Values) > i {
			if mapped := openmeteo.MapWeatherCode(weatherCodes.Values[i]); mapped != "" {
				description = theme.Accent.Render(mapped)
			}
		}

		minValue := "-"
		if hasMin && len(minTemps.Values) > i {
			minValue = formatValueWithUnit(minTemps.Values[i], minTemps.Unit)
		}

		maxValue := "-"
		if hasMax && len(maxTemps.Values) > i {
			maxValue = formatValueWithUnit(maxTemps.Values[i], maxTemps.Unit)
		}

		columns = append(columns, forecastColumn{
			heading:     label,
			description: description,
			details: []string{
				fmt.Sprintf("Min %s", minValue),
				fmt.Sprintf("Max %s", maxValue),
			},
		})
	}

	return columns
}

func renderColumns(style lipgloss.Style, columns []forecastColumn) string {
	rendered := make([]string, 0, len(columns))
	for _, column := range columns {
		lines := make([]string, 0, 2+len(column.details))
		if column.heading != "" {
			lines = append(lines, column.heading)
		}
		if column.description != "" {
			lines = append(lines, column.description)
		}
		lines = append(lines, column.details...)
		rendered = append(rendered, style.Render(lipgloss.JoinVertical(lipgloss.Left, lines...)))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func renderBody(width int, header, current, hourly, daily string) string {
	header = renderSection(width, header, false)
	current = renderSection(width, current, true)
	hourly = renderSection(width, hourly, true)
	daily = renderSection(width, daily, false)
	return lipgloss.JoinVertical(lipgloss.Left, header, current, hourly, daily)
}

func renderSection(width int, content string, withDivider bool) string {
	if width < 0 {
		width = 0
	}

	style := blockSpacing
	if withDivider {
		divider := sectionDivider
		if width > 0 {
			divider = divider.Width(width)
		}
		style = style.Inherit(divider)
	}
	if width > 0 {
		style = style.Width(width)
	}
	return style.Render(content)
}
