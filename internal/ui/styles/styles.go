package styles

import "github.com/charmbracelet/lipgloss"



var (
	Keyword = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	Subtle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	Dot     = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(" • ")
	Main    = lipgloss.NewStyle().MarginLeft(2)
	Cursor  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)
