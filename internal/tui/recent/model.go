package recent

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/theme"
)

type view int

const (
	viewList = iota
	viewError
)

type Model struct {
	size theme.Size
	view view
	list list.Model
	keys keyMap
	help help.Model
}

func (m Model) Init() tea.Cmd {
	return getRecentLocationsCmd()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.pick) {
			picked, ok := m.list.SelectedItem().(recentLocationItem)
			if ok {
				return m, pickCmd(picked.GeocodingResult, true)
			}
		}
		if key.Matches(msg, m.keys.newSearch) {
			return m, requestNewSearchCmd()
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
		if len(msg.locations) == 0 {
			return m, pickCmd(openmeteo.GeocodingResult{}, false)
		}
		if len(msg.locations) == 1 {
			return m, pickCmd(msg.locations[0], true)
		}

		items := make([]list.Item, len(msg.locations))
		for i, loc := range msg.locations {
			items[i] = recentLocationItem{loc}
		}
		m.list.SetItems(items)
		return m, nil
	case errorMsg:
		m.view = viewError
		return m, nil
	}

	// Forward messages to sub-components
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if !m.size.Ready {
		return theme.OuterFrame.Render("Init...")
	}

	frameX, frameY := theme.OuterFrame.GetFrameSize()
	innerWidth := m.size.Width - frameX
	innerHeight := m.size.Height - frameY

	var content string
	switch m.view {
	case viewList:
		content = lipgloss.JoinVertical(lipgloss.Left, "Recent locations:", "", m.list.View(), m.help.View(m.keys))
	case viewError:
		content = "Error occurred"
	default:
		content = "unknown state (recent)"
	}

	return theme.OuterFrame.Render(lipgloss.Place(
		innerWidth,
		innerHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	))
}

func (m Model) Reset() Model {
	m.list.ResetSelected()
	return m
}

func New() Model {
	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false
	list := list.New([]list.Item{}, listDelegate, 50, 14)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	return Model{
		view: viewList,
		list: list,
		keys: newKeyMap(),
		help: help.New(),
	}
}
