package weather

import (
	"fmt"

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

	columnWidthStyle = lipgloss.NewStyle().Width(18)

	columnBorderStyle = lipgloss.NewStyle().
				BorderRight(true).
				BorderStyle(lipgloss.MarkdownBorder()).
				BorderBottomForeground(lipgloss.Color("13"))
)

func renderDefault() string {
	return theme.OuterFrameStyle.Render("Unknown state (weather forecast screen).")
}

func renderError() string {
	return theme.OuterFrameStyle.Render("An error has occurred. Press 'q' to quit.")
}

func renderLoading(ellipsis spinner.Model) string {
	return theme.OuterFrameStyle.Render(fmt.Sprintf("Loading forecast%s", ellipsis.View()))
}

func renderHeader(location openmeteo.GeocodingResult, forecast openmeteo.ForecastResponse) string {
	place := location.Name
	if location.Admin1 != "" {
		place = place + ", " + location.Admin1
	}
	if location.Country != "" {
		place = place + ", " + location.Country
	}
	if weatherCode, ok := forecast.CurrentMeasurement(openmeteo.CurrentWeatherCode); ok {
		conditions := openmeteo.MapWeatherCode(weatherCode.Value)
		conditions = theme.AccentStyle.Render(conditions)
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

	minTempLabel := theme.LabelStyle.Render("\nMin")
	var minTempValue string
	if minSeries, ok := weather.DailySeries(openmeteo.DailyTemperature2mMin); ok && len(minSeries.Values) > 0 {
		minTempValue = formatValueWithUnit(minSeries.Values[0], minSeries.Unit)
	} else {
		minTempValue = "-"
	}
	s += minTempLabel + minTempValue

	maxTempLabel := theme.LabelStyle.Render("\nMax")
	var maxTempValue string
	if maxSeries, ok := weather.DailySeries(openmeteo.DailyTemperature2mMax); ok && len(maxSeries.Values) > 0 {
		maxTempValue = formatValueWithUnit(maxSeries.Values[0], maxSeries.Unit)
	} else {
		maxTempValue = "-"
	}
	s += maxTempLabel + maxTempValue

	windLabel := theme.LabelStyle.Render("\n\nWind")
	windValue := "-"
	if windSpeed, ok := weather.CurrentMeasurement(openmeteo.CurrentWindSpeed10m); ok {
		windValue = formatMeasurement(windSpeed)
		if windDirection, ok := weather.CurrentMeasurement(openmeteo.CurrentWindDirection10m); ok {
			windValue = fmt.Sprintf("%.1f %s @ %.1f %s", windSpeed.Value, windSpeed.Unit, windDirection.Value, windDirection.Unit)
		}
	}
	s += windLabel + windValue

	windGustsLabel := theme.LabelStyle.Render("\nWind gusts")
	windGustsValue := "-"
	if windGusts, ok := weather.CurrentMeasurement(openmeteo.CurrentWindGusts10m); ok {
		windGustsValue = formatMeasurement(windGusts)
	}
	s += windGustsLabel + windGustsValue

	humidityLabel := theme.LabelStyle.Render("\nHumidity")
	humidityValue := "-"
	if humidity, ok := weather.CurrentMeasurement(openmeteo.CurrentRelativeHumidity2m); ok {
		humidityValue = formatMeasurement(humidity)
	}
	s += humidityLabel + humidityValue

	precipitationLabel := theme.LabelStyle.Render("\nPrecipitation")
	precipitationValue := "-"
	if precipitation, ok := weather.CurrentMeasurement(openmeteo.CurrentPrecipitation); ok {
		precipitationValue = formatMeasurement(precipitation)
	}
	s += precipitationLabel + precipitationValue

	pressureLabel := theme.LabelStyle.Render("\nPressure")
	pressureValue := "-"
	if pressure, ok := weather.CurrentMeasurement(openmeteo.CurrentSeaLevelPressure); ok {
		pressureValue = formatMeasurement(pressure)
	}
	s += pressureLabel + pressureValue

	uvLabel := theme.LabelStyle.Render("\nUV index")
	var uvValue string
	if uvSeries, ok := weather.DailySeries(openmeteo.DailyUVIndexMax); ok && len(uvSeries.Values) > 0 {
		uvValue = fmt.Sprintf("%.1f", uvSeries.Values[0])
	} else {
		uvValue = "-"
	}
	s += uvLabel + uvValue
	return lipgloss.NewStyle().PaddingBottom(1).Render(s)
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

		column := lipgloss.JoinVertical(
			lipgloss.Left,
			dayStr,
			wmoStr,
			fmt.Sprintf("Min %s", minStr),
			fmt.Sprintf("Max %s", maxStr),
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
