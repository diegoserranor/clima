package weather

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/store"
)

func getForecastCmd(lat float64, long float64) tea.Cmd {
	return func() tea.Msg {
		params := openmeteo.ForecastParams{
			Latitude:      lat,
			Longitude:     long,
			Timezone:      "auto",
			ForecastHours: 6,
			ForecastDays:  6,
			Current: []openmeteo.CurrentVariables{
				openmeteo.CurrentTemperature2m,
				openmeteo.CurrentApparentTemperature,
				openmeteo.CurrentRelativeHumidity2m,
				openmeteo.CurrentIsDay,
				openmeteo.CurrentWeatherCode,
				openmeteo.CurrentWindSpeed10m,
				openmeteo.CurrentWindDirection10m,
				openmeteo.CurrentWindGusts10m,
				openmeteo.CurrentPrecipitation,
				openmeteo.CurrentSeaLevelPressure,
			},
			Daily: []openmeteo.DailyVariables{
				openmeteo.DailyTemperature2mMin,
				openmeteo.DailyTemperature2mMax,
				openmeteo.DailyWeatherCode,
				openmeteo.DailyUVIndexMax,
			},
			Hourly: []openmeteo.HourlyVariables{
				openmeteo.HourlyTemperature2m,
				openmeteo.HourlyWeatherCode,
			},
		}
		res, err := openmeteo.GetForecast(params)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}
		return dataMsg{
			forecast: res,
		}
	}
}

func requestNewSearchCmd() tea.Cmd {
	return func() tea.Msg {
		return NewSearchMsg{}
	}
}

func requestRecentCmd() tea.Cmd {
	return func() tea.Msg {
		return RecentMsg{}
	}
}

func saveRecentLocationCmd(location openmeteo.GeocodingResult) tea.Cmd {
	return func() tea.Msg {
		err := store.AddRecentLocation(location)
		return savedMsg{err: err}
	}
}
