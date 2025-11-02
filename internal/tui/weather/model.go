package weather

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
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

type windowState int

const (
	windowInit windowState = iota
	windowReady
)

type Model struct {
	sink        io.Writer
	windowState windowState
	dataState   dataState
	viewport    viewport.Model
	keys        keyMap
	errStr      string
	ellipsis    spinner.Model
	location    openmeteo.GeocodingResult
	forecast    openmeteo.ForecastResponse
	help        string
	// help        help.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		saveRecentLocationCmd(m.location),
		getForecastCmd(m.location.Latitude, m.location.Longitude),
		m.ellipsis.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.sink != nil {
		now := time.Now()
		nowStr := now.Format(time.DateTime)
		fmt.Fprintf(m.sink, "%s [msg @ weather model] %T: %+v\n", nowStr, msg, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.newSearch) {
			cmds = append(cmds, requestNewSearchCmd())
		}
		if key.Matches(msg, m.keys.recentLocations) {
			cmds = append(cmds, requestRecentCmd())
		}
		if key.Matches(msg, m.keys.refresh) && m.dataState == dataReady {
			m.dataState = dataLoading
			batched := tea.Batch(getForecastCmd(m.location.Latitude, m.location.Longitude), m.ellipsis.Tick)
			cmds = append(cmds, batched)
		}
		if key.Matches(msg, m.keys.quit) {
			cmds = append(cmds, tea.Quit)
		}
	case tea.WindowSizeMsg:
		if m.windowState == windowInit {
			m.windowState = windowReady
			m.viewport = viewport.New(msg.Width, msg.Height-lipgloss.Height(m.help))
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - lipgloss.Height(m.help)
		}
	case dataMsg:
		m.forecast = msg.forecast
		m.dataState = dataReady

		// build the body once, when data is ready
		frameX, _ := theme.OuterFrameStyle.GetFrameSize()
		innerWidth := m.viewport.Width - frameX

		header := renderHeader(m.location, m.forecast)
		current := renderCurrent(m.forecast)
		hourly := renderHourly(innerWidth, m.forecast)
		daily := renderDaily(innerWidth, m.forecast)
		body := renderBody(innerWidth, header, current, hourly, daily)

		m.viewport.SetContent(theme.OuterFrameStyle.Render(body))
	case errorMsg:
		m.dataState = dataError
		m.errStr = msg.err.Error()
	}

	if m.dataState == dataLoading {
		m.ellipsis, cmd = m.ellipsis.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.windowState == windowInit {
		return theme.OuterFrameStyle.Render("Init...")
	}

	var content string
	switch m.dataState {
	case dataError:
		content = renderError()
	case dataLoading:
		content = renderLoading(m.ellipsis)
	case dataReady:
		content = fmt.Sprintf("%s\n%s", m.viewport.View(), m.help)
	default:
		content = "unknown state (weather)"
	}

	return content
}

func (m Model) Reset(location openmeteo.GeocodingResult) Model {
	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.AccentStyle
	m.ellipsis = ellipsis
	m.dataState = dataLoading
	m.location = location
	return m
}

func New(location openmeteo.GeocodingResult, sink io.Writer) Model {
	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.AccentStyle

	keys := newKeyMap()

	help := help.New().View(keys)
	help = theme.OuterFrameStyle.Render(help)

	return Model{
		sink:      sink,
		dataState: dataLoading,
		location:  location,
		ellipsis:  ellipsis,
		keys:      newKeyMap(),
		help:      help,
	}
}
