package weather

import "github.com/diegoserranor/clima/internal/openmeteo"

type dataMsg struct {
	forecast openmeteo.ForecastResponse
}

type errorMsg struct {
	err error
}

type savedMsg struct {
	err error
}

type NewSearchMsg struct{}

type RecentMsg struct{}
