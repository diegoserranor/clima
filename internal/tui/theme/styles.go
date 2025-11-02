package theme

import "github.com/charmbracelet/lipgloss"

var (
	AccentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("13"))

	LabelStyle = lipgloss.NewStyle().
			Width(20).
			Foreground(lipgloss.Color("8"))

	SubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	OuterFrameStyle = lipgloss.NewStyle().
			Padding(1, 2)

	KeyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("13")).
			PaddingLeft(1).
			PaddingRight(1)
)
