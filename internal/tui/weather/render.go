package weather

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/theme"
)

var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("13")).
			Foreground(lipgloss.Color("0")).
			MarginBottom(1).
			PaddingLeft(1).
			PaddingRight(1).
			Italic(true)

	dividerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderBottomForeground(lipgloss.Color("13"))

	columnWidthStyle = lipgloss.NewStyle().Width(22)

	columnBorderStyle = lipgloss.NewStyle().
				BorderRight(true).
				BorderStyle(lipgloss.MarkdownBorder()).
				BorderBottomForeground(lipgloss.Color("13"))
)

func renderDefault() string {
	return theme.OuterFrameStyle.Render("Unknown state (weather forecast screen).")
}

func renderError(errStr string) string {
	content := fmt.Sprintf("An error has occurred:\n%s\n\nPress 'q' to quit.", errStr)
	return theme.OuterFrameStyle.Render(content)
}

func renderLoading(ellipsis spinner.Model) string {
	return theme.OuterFrameStyle.Render(fmt.Sprintf("Loading forecast%s", ellipsis.View()))
}

func renderHeader(location openmeteo.GeocodingResult) string {
	header := location.Name
	parts := []string{}

	if location.Admin1 != "" {
		parts = append(parts, location.Admin1)
	}
	if location.Country != "" {
		parts = append(parts, location.Country)
	}
	if len(parts) > 0 {
		header += "\n" + theme.SubtleStyle.Render(strings.Join(parts, ", "))
	}
	return lipgloss.NewStyle().MarginBottom(1).Render(header)
}

func renderCurrent(forecast openmeteo.ForecastResponse) string {
	var current string
	if weatherCode, ok := forecast.CurrentMeasurement(openmeteo.CurrentWeatherCode); ok {
		icon := openmeteo.MapWeatherIcon(weatherCode.Value)
		icon = theme.AccentStyle.Render(icon)
		conditions := openmeteo.MapWeatherCode(weatherCode.Value)
		conditions = theme.AccentStyle.Render(conditions)
		current = icon + "\n" + conditions
	}

	if currentTemp, ok := forecast.CurrentMeasurement(openmeteo.CurrentTemperature2m); ok {
		temperature := formatMeasurement(currentTemp)
		current += "\n" + temperature
	}

	if currentApparentTemp, ok := forecast.CurrentMeasurement(openmeteo.CurrentApparentTemperature); ok {
		feelsLikeVal := formatMeasurement(currentApparentTemp)
		current += "\n" + theme.SubtleStyle.Render(fmt.Sprintf("(feels like %s)", feelsLikeVal))
	}

	return current + "\n"
}

func renderCurrentDetails(forecast openmeteo.ForecastResponse) string {
	var currentdetails string
	var col1 string
	var col2 string
	var col3 string

	minTempLabel := theme.LabelStyle.Render("Min")
	var minTempValue string
	if minSeries, ok := forecast.DailySeries(openmeteo.DailyTemperature2mMin); ok && len(minSeries.Values) > 0 {
		minTempValue = formatValueWithUnit(minSeries.Values[0], minSeries.Unit)
	} else {
		minTempValue = "-"
	}
	col1 += fmt.Sprintf("%s%s\n", minTempLabel, minTempValue)

	maxTempLabel := theme.LabelStyle.Render("Max")
	var maxTempValue string
	if maxSeries, ok := forecast.DailySeries(openmeteo.DailyTemperature2mMax); ok && len(maxSeries.Values) > 0 {
		maxTempValue = formatValueWithUnit(maxSeries.Values[0], maxSeries.Unit)
	} else {
		maxTempValue = "-"
	}
	col1 += fmt.Sprintf("%s%s\n", maxTempLabel, maxTempValue)

	uvLabel := theme.LabelStyle.Render("UV index")
	var uvValue string
	if uvSeries, ok := forecast.DailySeries(openmeteo.DailyUVIndexMax); ok && len(uvSeries.Values) > 0 {
		uvValue = fmt.Sprintf("%.1f", uvSeries.Values[0])
	} else {
		uvValue = "-"
	}
	col1 += fmt.Sprintf("%s%s", uvLabel, uvValue)
	col1 = columnWidthStyle.Inherit(columnBorderStyle).MarginRight(2).Render(col1)

	windLabel := theme.LabelStyle.Render("Wind")
	windValue := "-"
	if windSpeed, ok := forecast.CurrentMeasurement(openmeteo.CurrentWindSpeed10m); ok {
		windValue = formatMeasurement(windSpeed)
	}
	col2 += fmt.Sprintf("%s%s\n", windLabel, windValue)

	gustsLabel := theme.LabelStyle.Render("Gusts")
	gustsValue := "-"
	if windGusts, ok := forecast.CurrentMeasurement(openmeteo.CurrentWindGusts10m); ok {
		gustsValue = formatMeasurement(windGusts)
	}
	currentdetails += gustsLabel + gustsValue
	col2 += fmt.Sprintf("%s%s\n", gustsLabel, gustsValue)

	directionLabel := theme.LabelStyle.Render("Direction")
	directionValue := "-"
	if windDirection, ok := forecast.CurrentMeasurement(openmeteo.CurrentWindDirection10m); ok {
		directionValue = formatMeasurement(windDirection)
	}
	col2 += fmt.Sprintf("%s%s", directionLabel, directionValue)
	col2 = columnWidthStyle.Inherit(columnBorderStyle).MarginRight(2).Render(col2)

	humidityLabel := theme.LabelStyle.Render("Humidity")
	humidityValue := "-"
	if humidity, ok := forecast.CurrentMeasurement(openmeteo.CurrentRelativeHumidity2m); ok {
		humidityValue = formatMeasurement(humidity)
	}
	col3 += fmt.Sprintf("%s%s\n", humidityLabel, humidityValue)

	precipitationLabel := theme.LabelStyle.Render("Precip")
	precipitationValue := "-"
	if precipitation, ok := forecast.CurrentMeasurement(openmeteo.CurrentPrecipitation); ok {
		precipitationValue = formatMeasurement(precipitation)
	}
	col3 += fmt.Sprintf("%s%s\n", precipitationLabel, precipitationValue)

	pressureLabel := theme.LabelStyle.Render("Pressure")
	pressureValue := "-"
	if pressure, ok := forecast.CurrentMeasurement(openmeteo.CurrentSeaLevelPressure); ok {
		pressureValue = formatMeasurement(pressure)
	}
	currentdetails += pressureLabel + pressureValue
	col3 += fmt.Sprintf("%s%s", pressureLabel, pressureValue)
	col3 = columnWidthStyle.Render(col3)

	currentdetails = lipgloss.JoinHorizontal(lipgloss.Top, col1, col2, col3)

	return lipgloss.NewStyle().PaddingBottom(1).Render(currentdetails)
}

// Renders the forecast for the next few hours except for the current hour. The number of rendered hours depends on the total width available.
func renderHourly(width int, forecast openmeteo.ForecastResponse) string {
	// cw -> the width of the column without right margin
	// mr -> the right margin for every column except the last
	// width -> total available width
	cw := columnWidthStyle.GetWidth()
	mr := 2
	maxAllowed := (width + mr) / (cw + mr)
	if maxAllowed < 1 {
		return theme.SubtleStyle.Render("The terminal window is too small")
	}

	// We need at least "now" + 1 future hour to do anything useful.
	if len(forecast.HourlyTimes) < 2 {
		return theme.SubtleStyle.Render("Hourly forecast unavailable")
	}
	hourlySeries := forecast.HourlyTimes[1:]
	hourCount := len(hourlySeries)

	// Clamp max allowed to the available hourly data series from the API response
	maxAvailable := hourCount
	if maxAllowed > maxAvailable {
		maxAllowed = maxAvailable
	}

	var wmoSeries []float64
	weatherCodes, hasWeatherCodes := forecast.HourlySeries(openmeteo.HourlyWeatherCode)
	if hasWeatherCodes {
		wmoSeries = weatherCodes.Values[1:]
	}

	var tempSeries []float64
	temperatures, hasTemperatures := forecast.HourlySeries(openmeteo.HourlyTemperature2m)
	if hasTemperatures {
		tempSeries = temperatures.Values[1:]
	}

	cols := make([]string, 0, maxAllowed)
	for i := range maxAllowed {
		timeStr := theme.SubtleStyle.Render(formatHourlyTime(hourlySeries[i]))

		wmoStr := "-"
		if hasWeatherCodes {
			if mapped := openmeteo.MapWeatherCode(wmoSeries[i]); mapped != "" {
				wmoStr = theme.AccentStyle.Render(mapped)
			}
		}

		tempStr := "-"
		if hasTemperatures {
			tempStr = formatValueWithUnit(tempSeries[i], temperatures.Unit)
		}

		column := lipgloss.JoinVertical(lipgloss.Left, timeStr, wmoStr, tempStr)
		style := columnWidthStyle
		if i != maxAllowed-1 {
			style = style.Inherit(columnBorderStyle).MarginRight(mr)
		}
		column = style.Render(column)
		cols = append(cols, column)
	}

	hourColumns := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
	hourly := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("Next few hours"), hourColumns)
	return lipgloss.NewStyle().PaddingBottom(1).Render(hourly)
}

// Renders the forecast for the next few days except for the current day. The number of days rendered depends on the available width.
func renderDaily(width int, forecast openmeteo.ForecastResponse) string {
	// cw -> the width of the column without right margin
	// mr -> the right margin for every column except the last
	// width -> total available width
	cw := columnWidthStyle.GetWidth()
	mr := 2
	maxAllowed := (width + mr) / (cw + mr)
	if maxAllowed < 1 {
		return theme.SubtleStyle.Render("The terminal window is too small")
	}

	// We need at least "today" + 1 future day to do anything useful.
	if len(forecast.DailyTimes) < 2 {
		return theme.SubtleStyle.Render("Daily forecast unavailable")
	}
	dailySeries := forecast.DailyTimes[1:]
	dayCount := len(dailySeries)

	// Clamp max allowed to the available daily data series from the API response
	maxAvailable := dayCount
	if maxAllowed > maxAvailable {
		maxAllowed = maxAvailable
	}

	var wmoSeries []float64
	weatherCodes, hasCodes := forecast.DailySeries(openmeteo.DailyWeatherCode)
	if hasCodes {
		wmoSeries = weatherCodes.Values[1:]
	}

	var minSeries []float64
	minTemps, hasMin := forecast.DailySeries(openmeteo.DailyTemperature2mMin)
	if hasMin {
		minSeries = minTemps.Values[1:]
	}

	var maxSeries []float64
	maxTemps, hasMax := forecast.DailySeries(openmeteo.DailyTemperature2mMax)
	if hasMax {
		maxSeries = maxTemps.Values[1:]
	}

	cols := make([]string, 0, maxAllowed)
	for i := range maxAllowed {
		dayStr := theme.SubtleStyle.Render(formatDailyDate(forecast.DailyTimes[i]))

		wmoStr := "-"
		if hasCodes {
			if mapped := openmeteo.MapWeatherCode(wmoSeries[i]); mapped != "" {
				wmoStr = theme.AccentStyle.Render(mapped)
			}
		}

		minStr := "-"
		if hasMin {
			minStr = formatValueWithUnit(minSeries[i], minTemps.Unit)
		}

		maxStr := "-"
		if hasMax {
			maxStr = formatValueWithUnit(maxSeries[i], maxTemps.Unit)
		}

		minLabel := theme.LabelStyle.Render("Min")
		maxLabel := theme.LabelStyle.Render("Max")

		column := lipgloss.JoinVertical(
			lipgloss.Left,
			dayStr,
			wmoStr,
			fmt.Sprintf("%s%s", minLabel, minStr),
			fmt.Sprintf("%s%s", maxLabel, maxStr),
		)
		style := columnWidthStyle
		if i != maxAllowed-1 {
			style = style.Inherit(columnBorderStyle).MarginRight(mr)
		}
		column = style.Render(column)
		cols = append(cols, column)
	}

	dailyColumns := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
	daily := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("Next few days"), dailyColumns)
	return lipgloss.NewStyle().Render(daily)
}

func renderBody(width int, header, current, currentDetails, hourly, daily string) string {
	header = renderSection(width, header, false)
	current = renderSection(width, current, false)
	currentDetails = renderSection(width, currentDetails, true)
	hourly = renderSection(width, hourly, true)
	daily = renderSection(width, daily, false)
	return lipgloss.JoinVertical(lipgloss.Left, header, current, currentDetails, hourly, daily)
}

func renderSection(width int, content string, withDivider bool) string {
	if width < 0 {
		width = 0
	}

	style := lipgloss.NewStyle()
	if withDivider {
		divider := dividerStyle
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
