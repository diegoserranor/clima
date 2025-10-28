package search

import "github.com/charmbracelet/bubbles/key"

type inputKeyMap struct {
	submit     key.Binding
	exitSearch key.Binding
}

func (k inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.submit, k.exitSearch}
}

func (k inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.submit, k.exitSearch},
	}
}

func newInputKeyMap() inputKeyMap {
	return inputKeyMap{
		submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "search"),
		),
		exitSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit search"),
		),
	}
}

type listKeyMap struct {
	up        key.Binding
	down      key.Binding
	pick      key.Binding
	newSearch key.Binding
	quit      key.Binding
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.up, k.down, k.pick, k.newSearch, k.quit}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up},
		{k.down},
		{k.pick},
		{k.newSearch},
		{k.quit},
	}
}

func newListKeyMap() listKeyMap {
	return listKeyMap{
		up: key.NewBinding(
			key.WithKeys("↑"),
			key.WithHelp("↑", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("↓"),
			key.WithHelp("↓", "down"),
		),
		pick: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "pick"),
		),
		newSearch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new search"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}
