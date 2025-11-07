package main

import (
	"flag"
	"fmt"

	"github.com/esferadigital/clima/internal/openmeteo"
)

func main() {
	code := flag.Int("code", -1, "WMO weather code to render (prints all when omitted)")
	flag.Parse()

	if *code >= 0 {
		printIcon(*code)
		return
	}

	for _, c := range openmeteo.WeatherIconCodes() {
		printIcon(c)
		fmt.Println()
	}
}

func printIcon(code int) {
	fmt.Printf("Code %d\n", code)
	fmt.Println(openmeteo.MapWeatherIcon(float64(code)))
}
