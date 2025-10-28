package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/internal/openmeteo"
)

func searchLocationsCmd(name string) tea.Cmd {
	return func() tea.Msg {
		params := openmeteo.GeocodingParams{
			Name:  name,
			Count: DEFAULT_SEARCH_COUNT,
		}
		res, err := openmeteo.SearchLocation(params)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}
		return dataMsg{
			locations: res.Results,
		}
	}
}

func pickCmd(location openmeteo.GeocodingResult) tea.Cmd {
	return func() tea.Msg {
		return SearchCompleteMsg{
			Location: location,
		}
	}
}

func requestRecentCmd() tea.Cmd {
	return func() tea.Msg {
		return RecentMsg{}
	}
}
