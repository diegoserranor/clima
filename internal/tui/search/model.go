package search

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/esferadigital/clima/internal/tui/theme"
)

const DEFAULT_SEARCH_COUNT = 10

type view int

const (
	viewInput = iota
	viewLoading
	viewPick
	viewError
)

type Model struct {
	size      theme.Size
	view      view
	input     textinput.Model
	inputKeys inputKeyMap
	ellipsis  spinner.Model
	list      list.Model
	listKeys  listKeyMap
	help      help.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.view == viewInput {
			if key.Matches(msg, m.inputKeys.submit) {
				m.view = viewLoading
				return m, tea.Batch(searchLocationsCmd(m.input.Value()), m.ellipsis.Tick)
			}
			if key.Matches(msg, m.inputKeys.exitSearch) {
				return m, requestRecentCmd()
			}
		}
		if m.view == viewPick {
			if key.Matches(msg, m.inputKeys.submit) {
				picked, ok := m.list.SelectedItem().(searchListItem)
				if ok {
					return m, pickCmd(picked.GeocodingResult)
				}
			}
			if key.Matches(msg, m.listKeys.newSearch) {
				m.input.Reset()
				m.list.ResetSelected()
				m.view = viewInput
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height
		m.size.Ready = true
		return m, nil
	case dataMsg:
		if len(msg.locations) == 1 {
			return m, pickCmd(msg.locations[0])
		}
		items := make([]list.Item, len(msg.locations))
		for i, loc := range msg.locations {
			items[i] = searchListItem{loc}
		}
		m.list.SetItems(items)
		m.view = viewPick
		return m, nil
	case errorMsg:
		m.view = viewError
		return m, nil
	}

	// Forward messages to sub-components
	var cmd tea.Cmd
	switch m.view {
	case viewInput:
		m.input, cmd = m.input.Update(msg)
	case viewLoading:
		m.ellipsis, cmd = m.ellipsis.Update(msg)
	case viewPick:
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if !m.size.Ready {
		return theme.OuterFrameStyle.Render("Init...")
	}

	frameX, frameY := theme.OuterFrameStyle.GetFrameSize()
	innerWidth := m.size.Width - frameX
	innerHeight := m.size.Height - frameY

	var content string
	switch m.view {
	case viewInput:
		content = lipgloss.JoinVertical(lipgloss.Left, "Location search:", m.input.View(), "", m.help.View(m.inputKeys))
	case viewLoading:
		content = fmt.Sprintf("Finding location%s", m.ellipsis.View())
	case viewPick:
		content = lipgloss.JoinVertical(lipgloss.Left, "Pick a location:", "", m.list.View(), "", m.help.View(m.listKeys))
	case viewError:
		content = "Error ... try again or quit"
	default:
		content = "unknown state (search)"
	}

	return theme.OuterFrameStyle.Render(lipgloss.Place(
		innerWidth,
		innerHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	))
}

func (m Model) Reset() Model {
	m.input.Reset()
	m.list.ResetSelected()
	return m
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "Salinas"
	input.Focus()
	input.CharLimit = 256
	input.Width = 50
	input.Cursor.Style = theme.AccentStyle

	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.AccentStyle

	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false
	list := list.New([]list.Item{}, listDelegate, 50, 14)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	return Model{
		view:      viewInput,
		input:     input,
		inputKeys: newInputKeyMap(),
		ellipsis:  ellipsis,
		list:      list,
		listKeys:  newListKeyMap(),
		help:      help.New(),
	}
}
