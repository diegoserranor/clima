package weather

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/theme"
)

type dataState int

const (
	dataLoading dataState = iota
	dataReady
	dataError
)

type sizeState int

const (
	sizeInit sizeState = iota
	sizeReady
)

type Model struct {
	sink      io.Writer
	dataState dataState
	size      theme.Size
	errStr    string
	ellipsis  spinner.Model
	location  openmeteo.GeocodingResult
	forecast  openmeteo.ForecastResponse
	keys      keyMap
	help      help.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		saveRecentLocationCmd(m.location),
		getForecastCmd(m.location.Latitude, m.location.Longitude),
		m.ellipsis.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.sink != nil {
		now := time.Now()
		nowStr := now.Format(time.DateTime)
		fmt.Fprintf(m.sink, "%s [msg @ weather model] %T: %+v\n", nowStr, msg, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.newSearch) {
			return m, requestNewSearchCmd()
		}
		if key.Matches(msg, m.keys.recentLocations) {
			return m, requestRecentCmd()
		}
		if key.Matches(msg, m.keys.refresh) && m.dataState == dataReady {
			m.dataState = dataLoading
			return m, tea.Batch(getForecastCmd(m.location.Latitude, m.location.Longitude), m.ellipsis.Tick)
		}
		if key.Matches(msg, m.keys.quit) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height
		m.size.Ready = true
		return m, nil
	case dataMsg:
		m.forecast = msg.forecast
		m.dataState = dataReady
		return m, nil
	case errorMsg:
		m.dataState = dataError
		m.errStr = msg.err.Error()
		return m, nil
	}

	if m.dataState == dataLoading {
		var cmd tea.Cmd
		m.ellipsis, cmd = m.ellipsis.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if !m.size.Ready {
		return theme.OuterFrame.Render("Init...")
	}

	frameX, frameY := theme.OuterFrame.GetFrameSize()
	innerWidth := m.size.Width - frameX
	innerHeight := m.size.Height - frameY

	var content string
	switch m.dataState {
	case dataError:
		content = renderError()
	case dataLoading:
		content = renderLoading(m.ellipsis)
	case dataReady:
		header := renderHeader(m.location, m.forecast)
		current := renderCurrent(m.forecast)
		hourly := renderHourly(innerWidth, m.forecast)
		daily := renderDaily(innerWidth, m.forecast)
		body := renderBody(innerWidth, header, current, hourly, daily)
		help := m.help.View(m.keys)
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			body,
			"",
			help,
		)
	default:
		content = "unknown state (weather)"
	}

	return theme.OuterFrame.Render(lipgloss.Place(
		innerWidth,
		innerHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	))
}

func (m Model) Reset(location openmeteo.GeocodingResult) Model {
	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.Accent
	m.ellipsis = ellipsis
	m.dataState = dataLoading
	m.location = location
	return m
}

func New(location openmeteo.GeocodingResult, sink io.Writer) Model {
	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.Accent
	return Model{
		sink:      sink,
		dataState: dataLoading,
		location:  location,
		ellipsis:  ellipsis,
		keys:      newKeyMap(),
		help:      help.New(),
	}
}
