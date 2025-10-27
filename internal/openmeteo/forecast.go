package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Parameters for the Open-Meteo Forecast V1 API.
// These are not exclusive. Check the docs for additional ones.
// https://open-meteo.com/en/docs
type ForecastParams struct {
	Latitude      float64
	Longitude     float64
	Timezone      string
	ForecastHours int
	ForecastDays  int
	Current       []CurrentVariables
	Daily         []DailyVariables
	Hourly        []HourlyVariables
}

// FloatMeasurement pairs a numeric value with the unit reported by the API.
type FloatMeasurement struct {
	Value float64
	Unit  string
}

// FloatSeries holds a time series of numeric values and the shared unit.
type FloatSeries struct {
	Values []float64
	Unit   string
}

// ForecastResponse is a typed view of the Forecast V1 payload. The JSON is
// decoded into a private representation first, then mapped onto the strongly
// typed measurement structures below so callers never have to assert from any.
// Keys are trimmed to the known variable enums so downstream code can rely on
// compile-time checks.
type ForecastResponse struct {
	Latitude         float64
	Longitude        float64
	Elevation        float64
	GenerationTimeMs float64
	UTCOffsetSeconds int
	Timezone         string
	TimezoneAbbrev   string
	CurrentTime      string
	HourlyTimes      []string
	DailyTimes       []string
	Current          map[CurrentVariables]FloatMeasurement
	Daily            map[DailyVariables]FloatSeries
	Hourly           map[HourlyVariables]FloatSeries
}

// CurrentMeasurement retrieves a single current measurement if it was requested.
func (f ForecastResponse) CurrentMeasurement(variable CurrentVariables) (FloatMeasurement, bool) {
	if f.Current == nil {
		return FloatMeasurement{}, false
	}
	m, ok := f.Current[variable]
	return m, ok
}

// DailySeries retrieves a daily time series if it was requested.
func (f ForecastResponse) DailySeries(variable DailyVariables) (FloatSeries, bool) {
	if f.Daily == nil {
		return FloatSeries{}, false
	}
	series, ok := f.Daily[variable]
	return series, ok
}

// HourlySeries retrieves an hourly time series if it was requested.
func (f ForecastResponse) HourlySeries(variable HourlyVariables) (FloatSeries, bool) {
	if f.Hourly == nil {
		return FloatSeries{}, false
	}
	series, ok := f.Hourly[variable]
	return series, ok
}

// Variables available to request from the Open-Meteo Forecast V1 API for current weather.
type CurrentVariables string

const (
	CurrentTemperature2m       CurrentVariables = "temperature_2m"
	CurrentRelativeHumidity2m  CurrentVariables = "relative_humidity_2m"
	CurrentApparentTemperature CurrentVariables = "apparent_temperature"
	CurrentIsDay               CurrentVariables = "is_day"
	CurrentWeatherCode         CurrentVariables = "weather_code"
	CurrentCloudCover          CurrentVariables = "cloud_cover"
	CurrentSeaLevelPressure    CurrentVariables = "pressure_msl"
	CurrentSurfacePressure     CurrentVariables = "surface_pressure"
	CurrentPrecipitation       CurrentVariables = "precipitation"
	CurrentRain                CurrentVariables = "rain"
	CurrentShowers             CurrentVariables = "showers"
	CurrentSnowfall            CurrentVariables = "snowfall"
	CurrentWindSpeed10m        CurrentVariables = "wind_speed_10m"
	CurrentWindDirection10m    CurrentVariables = "wind_direction_10m"
	CurrentWindGusts10m        CurrentVariables = "wind_gusts_10m"
)

// Variables available to request from the Open-Meteo Forecast V1 API for daily weather.
type DailyVariables string

const (
	DailyTemperature2mMin DailyVariables = "temperature_2m_min"
	DailyTemperature2mMax DailyVariables = "temperature_2m_max"
	DailyWeatherCode      DailyVariables = "weathercode"
	DailyUVIndexMax       DailyVariables = "uv_index_max"
)

// Variables available to request from the Open-Meteo Forecast V1 API for hourly weather.
type HourlyVariables string

const (
	HourlyTemperature2m HourlyVariables = "temperature_2m"
	HourlyWeatherCode   HourlyVariables = "weathercode"
	HourlyPrecipitation HourlyVariables = "precipitation"
)

var (
	currentVariableLookup = map[string]CurrentVariables{
		string(CurrentTemperature2m):       CurrentTemperature2m,
		string(CurrentRelativeHumidity2m):  CurrentRelativeHumidity2m,
		string(CurrentApparentTemperature): CurrentApparentTemperature,
		string(CurrentIsDay):               CurrentIsDay,
		string(CurrentWeatherCode):         CurrentWeatherCode,
		string(CurrentCloudCover):          CurrentCloudCover,
		string(CurrentSeaLevelPressure):    CurrentSeaLevelPressure,
		string(CurrentSurfacePressure):     CurrentSurfacePressure,
		string(CurrentPrecipitation):       CurrentPrecipitation,
		string(CurrentRain):                CurrentRain,
		string(CurrentShowers):             CurrentShowers,
		string(CurrentSnowfall):            CurrentSnowfall,
		string(CurrentWindSpeed10m):        CurrentWindSpeed10m,
		string(CurrentWindDirection10m):    CurrentWindDirection10m,
		string(CurrentWindGusts10m):        CurrentWindGusts10m,
	}
	dailyVariableLookup = map[string]DailyVariables{
		string(DailyTemperature2mMin): DailyTemperature2mMin,
		string(DailyTemperature2mMax): DailyTemperature2mMax,
		string(DailyWeatherCode):      DailyWeatherCode,
		string(DailyUVIndexMax):       DailyUVIndexMax,
	}
	hourlyVariableLookup = map[string]HourlyVariables{
		string(HourlyTemperature2m): HourlyTemperature2m,
		string(HourlyWeatherCode):   HourlyWeatherCode,
		string(HourlyPrecipitation): HourlyPrecipitation,
	}
)

const FORECAST_API_URL = "https://api.open-meteo.com/v1/forecast"

type forecastResponseRaw struct {
	Latitude         float64           `json:"latitude"`
	Longitude        float64           `json:"longitude"`
	Elevation        float64           `json:"elevation"`
	GenerationTimeMs float64           `json:"generation_time_ms"`
	UTCOffsetSeconds int               `json:"utc_offset_seconds"`
	Timezone         string            `json:"timezone"`
	TimezoneAbbrev   string            `json:"timezone_abbreviation"`
	CurrentUnits     map[string]string `json:"current_units"`
	Current          map[string]any    `json:"current"`
	HourlyUnits      map[string]string `json:"hourly_units"`
	Hourly           map[string][]any  `json:"hourly"`
	DailyUnits       map[string]string `json:"daily_units"`
	Daily            map[string][]any  `json:"daily"`
}

// Retrieve the current forecast data for a given location and parameters.
// Data is provided by the Open-Meteo API. The payload is normalised into
// typed measurement structs so callers never have to downcast from any.
func GetForecast(params ForecastParams) (ForecastResponse, error) {
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f", FORECAST_API_URL, params.Latitude, params.Longitude)
	if params.Timezone != "" {
		url += fmt.Sprintf("&timezone=%s", params.Timezone)
	}
	if params.ForecastHours > 0 {
		url += fmt.Sprintf("&forecast_hours=%d", params.ForecastHours)
	}
	if params.ForecastDays > 0 {
		url += fmt.Sprintf("&forecast_days=%d", params.ForecastDays)
	}
	if len(params.Current) > 0 {
		currentVars := writeVariableCSV(params.Current)
		url += fmt.Sprintf("&current=%s", currentVars)
	}
	if len(params.Daily) > 0 {
		dailyVars := writeVariableCSV(params.Daily)
		url += fmt.Sprintf("&daily=%s", dailyVars)
	}
	if len(params.Hourly) > 0 {
		hourlyVars := writeVariableCSV(params.Hourly)
		url += fmt.Sprintf("&hourly=%s", hourlyVars)
	}

	resp, err := http.Get(url)
	if err != nil {
		return ForecastResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ForecastResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var raw forecastResponseRaw
	if err := decoder.Decode(&raw); err != nil {
		return ForecastResponse{}, err
	}

	return raw.toForecastResponse(), nil
}

// Compose a comma-separated string of variable names.
// This is used to send the variables as a query parameter.
func writeVariableCSV[T ~string](variables []T) string {
	variableNames := make([]string, len(variables))
	for i, variable := range variables {
		variableNames[i] = string(variable)
	}
	return strings.Join(variableNames, ",")
}

func (raw forecastResponseRaw) toForecastResponse() ForecastResponse {
	response := ForecastResponse{
		Latitude:         raw.Latitude,
		Longitude:        raw.Longitude,
		Elevation:        raw.Elevation,
		GenerationTimeMs: raw.GenerationTimeMs,
		UTCOffsetSeconds: raw.UTCOffsetSeconds,
		Timezone:         raw.Timezone,
		TimezoneAbbrev:   raw.TimezoneAbbrev,
		Current:          make(map[CurrentVariables]FloatMeasurement),
		Daily:            make(map[DailyVariables]FloatSeries),
		Hourly:           make(map[HourlyVariables]FloatSeries),
	}

	if raw.Current != nil {
		if currentTime, ok := raw.Current["time"].(string); ok {
			response.CurrentTime = currentTime
		}
	}
	if raw.Daily != nil {
		if times, ok := toStringSlice(raw.Daily["time"]); ok {
			response.DailyTimes = times
		}
	}
	if raw.Hourly != nil {
		if times, ok := toStringSlice(raw.Hourly["time"]); ok {
			response.HourlyTimes = times
		}
	}

	if raw.Current != nil && raw.CurrentUnits != nil {
		for key, unit := range raw.CurrentUnits {
			variable, ok := currentVariableLookup[key]
			if !ok {
				continue
			}
			value, ok := toFloat64(raw.Current[key])
			if !ok {
				continue
			}
			response.Current[variable] = FloatMeasurement{
				Value: value,
				Unit:  unit,
			}
		}
	}

	if raw.Daily != nil && raw.DailyUnits != nil {
		for key, unit := range raw.DailyUnits {
			variable, ok := dailyVariableLookup[key]
			if !ok {
				continue
			}
			values, ok := toFloatSlice(raw.Daily[key])
			if !ok {
				continue
			}
			response.Daily[variable] = FloatSeries{
				Values: values,
				Unit:   unit,
			}
		}
	}

	if raw.Hourly != nil && raw.HourlyUnits != nil {
		for key, unit := range raw.HourlyUnits {
			variable, ok := hourlyVariableLookup[key]
			if !ok {
				continue
			}
			values, ok := toFloatSlice(raw.Hourly[key])
			if !ok {
				continue
			}
			response.Hourly[variable] = FloatSeries{
				Values: values,
				Unit:   unit,
			}
		}
	}

	return response
}

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func toFloatSlice(values []any) ([]float64, bool) {
	if values == nil {
		return nil, false
	}
	result := make([]float64, 0, len(values))
	for _, value := range values {
		floatVal, ok := toFloat64(value)
		if !ok {
			return nil, false
		}
		result = append(result, floatVal)
	}
	return result, true
}

func toStringSlice(values []any) ([]string, bool) {
	if values == nil {
		return nil, false
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if str, ok := value.(string); ok {
			result = append(result, str)
			continue
		}
		return nil, false
	}
	return result, true
}
