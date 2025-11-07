package search

import (
	"fmt"
	"strings"

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

var inputStyle = theme.OuterFrameStyle.PaddingTop(0)

func New() Model {
	inputKeys := newInputKeyMap()

	inputHeader := theme.OuterFrameStyle.Render("Location search:")

	inputHelp := help.New().View(inputKeys)
	inputFooter := theme.OuterFrameStyle.Render(inputHelp)

	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = theme.AccentStyle

	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false
	selectedStyle := list.NewDefaultItemStyles().SelectedTitle
	listDelegate.Styles.SelectedTitle = selectedStyle.Foreground(theme.AccentColor).BorderForeground(theme.AccentColor)
	list := list.New([]list.Item{}, listDelegate, 0, 0)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	listKeys := newListKeyMap()

	listHeader := theme.OuterFrameStyle.Render("Pick a location:")

	listHelp := help.New().View(listKeys)
	listFooter := theme.OuterFrameStyle.Render(listHelp)

	return Model{
		view:        viewInput,
		inputKeys:   inputKeys,
		inputHeader: inputHeader,
		inputFooter: inputFooter,
		ellipsis:    ellipsis,
		list:        list,
		listKeys:    newListKeyMap(),
		listHeader:  listHeader,
		listFooter:  listFooter,
	}
}

type view int

const (
	viewInput = iota
	viewLoading
	viewPick
	viewError
)

type Model struct {
	windowReady bool
	view        view
	input       textinput.Model
	inputKeys   inputKeyMap
	inputHeader string
	inputFiller string
	inputFooter string
	ellipsis    spinner.Model
	list        list.Model
	listKeys    listKeyMap
	listHeader  string
	listFooter  string
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
		otherWidth, _ := theme.OuterFrameStyle.GetFrameSize()
		otherHeight := lipgloss.Height(m.inputHeader) + lipgloss.Height(m.inputFooter)

		if !m.windowReady {
			m.windowReady = true

			input := textinput.New()
			input.Placeholder = "Salinas"
			input.Focus()
			input.CharLimit = 256
			input.Width = msg.Width - otherWidth
			input.Cursor.Style = theme.AccentStyle
			m.input = input

			inputRendered := inputStyle.Render(m.input.View())
			usedHeight := lipgloss.Height(m.inputHeader) +
				lipgloss.Height(inputRendered) +
				lipgloss.Height(m.inputFooter)
			inputFillerHeight := max(msg.Height-usedHeight, 0)
			m.inputFiller = strings.Repeat("\n", inputFillerHeight)
		} else {
			m.input.Width = msg.Width - otherWidth
		}

		m.list.SetWidth(msg.Width - otherWidth)
		m.list.SetHeight(msg.Height - otherHeight)

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
	if !m.windowReady {
		return theme.OuterFrameStyle.Render("Init...")
	}

	var content string
	switch m.view {
	case viewInput:
		input := inputStyle.Render(m.input.View())
		content = fmt.Sprintf("%s\n%s%s\n%s", m.inputHeader, input, m.inputFiller, m.inputFooter)
	case viewLoading:
		content = fmt.Sprintf("Finding location%s", m.ellipsis.View())
		content = theme.OuterFrameStyle.Render(content)
	case viewPick:
		list := theme.OuterFrameStyle.Render(m.list.View())
		content = fmt.Sprintf("%s%s%s", m.listHeader, list, m.listFooter)
	case viewError:
		content = theme.OuterFrameStyle.Render("Error ... try again or quit")
	default:
		content = theme.OuterFrameStyle.Render("Unknown state (search)")
	}

	return content
}

func (m Model) Reset() Model {
	m.input.Reset()
	m.list.ResetSelected()
	m.view = viewInput
	return m
}
