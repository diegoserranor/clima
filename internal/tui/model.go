package tui

import (
	"fmt"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/recent"
	"github.com/esferadigital/clima/internal/tui/search"
	"github.com/esferadigital/clima/internal/tui/weather"
)

// ---- model ----

type route int

const (
	routeRecent = iota
	routeSearch
	routeWeather
)

type Model struct {
	sink    io.Writer
	route   route
	recent  recent.Model
	search  search.Model
	weather weather.Model
}

func (m Model) Init() tea.Cmd {
	return m.recent.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.sink != nil {
		now := time.Now()
		nowStr := now.Format(time.DateTime)
		fmt.Fprintf(m.sink, "%s [msg] %T: %+v\n", nowStr, msg, msg)
	}

	switch msg := msg.(type) {

	// terminal size changes
	case tea.WindowSizeMsg:
		m.recent, _ = m.recent.Update(msg)
		m.search, _ = m.search.Update(msg)
		m.weather, _ = m.weather.Update(msg)
		return m, nil

	// recent
	case recent.RecentCompleteMsg:
		if !msg.OK {
			m.route = routeSearch
			return m, m.search.Init()
		}
		m.route = routeWeather
		m.weather = m.weather.Reset(msg.Location)
		return m, m.weather.Init()
	case recent.NewSearchMsg:
		m.route = routeSearch
		return m, m.search.Init()

	// search
	case search.SearchCompleteMsg:
		m.route = routeWeather
		m.weather = m.weather.Reset(msg.Location)
		return m, m.weather.Init()
	case search.RecentMsg:
		m.route = routeRecent
		m.recent = m.recent.Reset()
		return m, m.recent.Init()

	// weather
	case weather.NewSearchMsg:
		m.route = routeSearch
		return m, m.search.Init()
	case weather.RecentMsg:
		m.route = routeRecent
		m.recent = m.recent.Reset()
		return m, m.recent.Init()
	}

	// Forward updates to sub-components
	var cmd tea.Cmd
	switch m.route {
	case routeRecent:
		m.recent, cmd = m.recent.Update(msg)
		return m, cmd
	case routeSearch:
		m.search, cmd = m.search.Update(msg)
		return m, cmd
	case routeWeather:
		m.weather, cmd = m.weather.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) View() string {
	var content string
	switch m.route {
	case routeRecent:
		content = m.recent.View()
	case routeSearch:
		content = m.search.View()
	case routeWeather:
		content = m.weather.View()
	default:
		content = "Unknown state (core)"
	}
	return content
}

func InitialModel(sink io.Writer) Model {
	return Model{
		sink:    sink,
		recent:  recent.New(),
		search:  search.New(),
		weather: weather.New(openmeteo.GeocodingResult{}, sink),
	}
}
