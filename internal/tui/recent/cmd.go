package recent

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/diegoserranor/clima/internal/openmeteo"
	"github.com/diegoserranor/clima/internal/store"
)

func getRecentLocationsCmd() tea.Cmd {
	return func() tea.Msg {
		locations, err := store.LoadRecentLocations()
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		return dataMsg{
			locations: locations,
		}
	}
}

func pickCmd(location openmeteo.GeocodingResult, ok bool) tea.Cmd {
	return func() tea.Msg {
		return RecentCompleteMsg{
			Location: location,
			OK:       ok,
		}
	}
}

func requestNewSearchCmd() tea.Cmd {
	return func() tea.Msg {
		return NewSearchMsg{}
	}
}
