package openmeteo

func MapWeatherCode(code float64) string {
	wmoCodes := map[float64]string{
		0:  "Clear",
		1:  "Mostly clear",
		2:  "Partly cloudy",
		3:  "Overcast",
		45: "Fog",
		48: "Icy fog",
		51: "Light drizzle",
		53: "Drizzle",
		55: "Heavy drizzle",
		56: "Light freezing drizzle",
		57: "Heavy freezing drizzle",
		61: "Light rain",
		63: "Rain",
		65: "Heavy rain",
		66: "Light freezing rain",
		67: "Heavy freezing rain",
		71: "Light snowfall",
		73: "Snowfall",
		75: "Heavy snowfall",
		77: "Snow grains",
		80: "Light rain showers",
		81: "Rain showers",
		82: "Heavy rain showers",
		85: "Light snow showers",
		86: "Heavy snow showers",
		95: "Thunderstorm",
		96: "Thunderstorm with light hail",
		99: "Thunderstorm with heavy hail",
	}

	return wmoCodes[code]
}
