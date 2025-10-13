package weather

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	newSearch       key.Binding
	recentLocations key.Binding
	refresh         key.Binding
	quit            key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.newSearch, k.recentLocations, k.refresh, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.newSearch}, {k.recentLocations},
		{k.refresh}, {k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
		newSearch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new search"),
		),
		recentLocations: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "recent locations"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}
