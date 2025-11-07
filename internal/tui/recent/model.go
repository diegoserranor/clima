package recent

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/tui/theme"
)

func New() Model {
	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.AccentStyle

	keys := newKeyMap()

	header := theme.OuterFrameStyle.Render("Recent locations:")

	help := help.New().View(keys)
	footer := theme.OuterFrameStyle.Render(help)

	return Model{
		windowReady: false,
		dataReady:   false,
		ellipsis:    ellipsis,
		keys:        keys,
		header:      header,
		footer:      footer,
	}
}

type Model struct {
	windowReady bool
	dataReady   bool
	errStr      string
	ellipsis    spinner.Model
	list        list.Model
	keys        keyMap
	header      string
	footer      string
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
		otherWidth, _ := theme.OuterFrameStyle.GetFrameSize()
		otherHeight := lipgloss.Height(m.header) + lipgloss.Height(m.footer)

		if !m.windowReady {
			m.windowReady = true
			listDelegate := list.NewDefaultDelegate()
			listDelegate.ShowDescription = false
			selectedStyle := list.NewDefaultItemStyles().SelectedTitle
			listDelegate.Styles.SelectedTitle = selectedStyle.Foreground(theme.AccentColor).BorderForeground(theme.AccentColor)
			list := list.New([]list.Item{}, listDelegate, msg.Width-otherWidth, msg.Height-otherHeight)
			list.SetShowStatusBar(false)
			list.SetFilteringEnabled(false)
			list.SetShowHelp(false)
			list.SetShowTitle(false)
			m.list = list
		} else {
			m.list.SetWidth(msg.Width - otherWidth)
			m.list.SetHeight(msg.Height - otherHeight)
		}
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
		m.dataReady = true
		return m, nil
	case errorMsg:
		m.dataReady = true
		m.errStr = msg.err.Error()
		return m, nil
	}

	// Forward messages to sub-components
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if !m.windowReady {
		return theme.OuterFrameStyle.Render("Init...")
	}

	if !m.dataReady {
		return theme.OuterFrameStyle.Render("Loading...")
	}

	if m.errStr != "" {
		return theme.OuterFrameStyle.Render(m.errStr)
	}

	list := theme.OuterFrameStyle.Render(m.list.View())
	return fmt.Sprintf("%s%s%s", m.header, list, m.footer)
}

func (m Model) Reset() Model {
	m.list.ResetSelected()
	return m
}
