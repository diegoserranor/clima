package theme

import "github.com/charmbracelet/lipgloss"

var (
	AccentStyle = lipgloss.NewStyle().
			Foreground(AccentColor)

	LabelStyle = lipgloss.NewStyle().
			Width(12).
			Foreground(SubtleColor)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)

	OuterFrameStyle = lipgloss.NewStyle().
			Padding(1, 2)

	KeyStyle = lipgloss.NewStyle().
			Background(AccentColor).
			PaddingLeft(1).
			PaddingRight(1)
)
