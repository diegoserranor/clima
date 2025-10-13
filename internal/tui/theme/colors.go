package theme

import "github.com/charmbracelet/lipgloss"

var (
	Accent = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	Label  = lipgloss.NewStyle().Width(20).Foreground(lipgloss.Color("8"))
	Subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
