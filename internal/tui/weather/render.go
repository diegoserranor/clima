package weather

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/theme"
)

var (
	title          = lipgloss.NewStyle().Background(lipgloss.Color("13")).Foreground(lipgloss.Color("0")).MarginBottom(1).PaddingLeft(1).PaddingRight(1).Italic(true)
	sectionDivider = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderBottomForeground(lipgloss.Color("13"))
	columnWidth    = lipgloss.NewStyle().Width(18)
	columnBorder   = lipgloss.NewStyle().BorderRight(true).BorderStyle(lipgloss.MarkdownBorder()).BorderBottomForeground(lipgloss.Color("13"))
)

type forecastColumn struct {
	heading     string
	description string
	details     []string
}

type hourlyColumn struct{}

type dailyColumn struct{}

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
	return lipgloss.NewStyle().PaddingBottom(1).Render(s)
}

func renderHourly(forecast openmeteo.ForecastResponse) string {
	hourCount := len(forecast.HourlyTimes)
	if hourCount == 0 {
		return theme.Subtle.Render("Hourly forecast unavailable")
	}

	weatherCodes, hasWeatherCodes := forecast.HourlySeries(openmeteo.HourlyWeatherCode)
	temperatures, hasTemperatures := forecast.HourlySeries(openmeteo.HourlyTemperature2m)

	// Skip the first item since that is the current hour
	hourColumns := ""
	for i := 1; i < hourCount; i++ {
		time := theme.Subtle.Render(formatHourlyTime(forecast.HourlyTimes[i]))

		conditions := "-"
		if hasWeatherCodes {
			if mapped := openmeteo.MapWeatherCode(weatherCodes.Values[i]); mapped != "" {
				conditions = theme.Accent.Render(mapped)
			}
		}

		temp := "-"
		if hasTemperatures {
			temp = formatValueWithUnit(temperatures.Values[i], temperatures.Unit)
		}

		column := lipgloss.JoinVertical(lipgloss.Left, time, conditions, temp)
		style := columnWidth
		if i != hourCount-1 {
			style = style.Inherit(columnBorder).MarginRight(2)
		}
		column = style.Render(column)
		hourColumns = lipgloss.JoinHorizontal(lipgloss.Top, hourColumns, column)
	}

	hourly := lipgloss.JoinVertical(lipgloss.Left, title.Render("Next few hours"), hourColumns)
	return lipgloss.NewStyle().PaddingBottom(1).Render(hourly)
}

func renderDaily(forecast openmeteo.ForecastResponse) string {
	dayCount := len(forecast.DailyTimes)
	if dayCount == 0 {
		return theme.Subtle.Render("Daily outlook unavailable")
	}

	minTemps, hasMin := forecast.DailySeries(openmeteo.DailyTemperature2mMin)
	maxTemps, hasMax := forecast.DailySeries(openmeteo.DailyTemperature2mMax)
	weatherCodes, hasCodes := forecast.DailySeries(openmeteo.DailyWeatherCode)

	// Skip the first item since that is the current day
	dailyColumns := ""
	for i := 1; i < dayCount; i++ {
		day := theme.Subtle.Render(formatDailyDate(forecast.DailyTimes[i]))

		conditions := "-"
		if hasCodes && len(weatherCodes.Values) > i {
			if mapped := openmeteo.MapWeatherCode(weatherCodes.Values[i]); mapped != "" {
				conditions = theme.Accent.Render(mapped)
			}
		}

		// Min temp
		minValue := "-"
		if hasMin && len(minTemps.Values) > i {
			minValue = formatValueWithUnit(minTemps.Values[i], minTemps.Unit)
		}

		// Max temp
		maxValue := "-"
		if hasMax && len(maxTemps.Values) > i {
			maxValue = formatValueWithUnit(maxTemps.Values[i], maxTemps.Unit)
		}

		column := lipgloss.JoinVertical(
			lipgloss.Left,
			day,
			conditions,
			fmt.Sprintf("Min %s", minValue),
			fmt.Sprintf("Max %s", maxValue),
		)
		style := columnWidth
		if i != dayCount-1 {
			style = style.Inherit(columnBorder).MarginRight(2)
		}

		column = style.Render(column)
		dailyColumns = lipgloss.JoinHorizontal(lipgloss.Top, dailyColumns, column)
	}

	daily := lipgloss.JoinVertical(lipgloss.Left, title.Render("Next few days"), dailyColumns)
	return lipgloss.NewStyle().Render(daily)
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

	style := lipgloss.NewStyle()
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
