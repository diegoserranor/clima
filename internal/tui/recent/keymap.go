package recent

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	up        key.Binding
	down      key.Binding
	pick      key.Binding
	newSearch key.Binding
	quit      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.up, k.down, k.pick, k.newSearch, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up},
		{k.down},
		{k.pick},
		{k.newSearch},
		{k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
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
