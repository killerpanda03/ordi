package menu

import (
	"example/ordi/internal/ui/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {

	s := "\n"
	s += lipgloss.JoinVertical(lipgloss.Center, "Wähle eine Option: \n\n")
	for i, choice := range m.choices {
		cursor := ""
		if m.cursor == i {
			cursor = styles.Cursor.Render(">")
			choice = styles.Cursor.Render(choice)
		}
		s += fmt.Sprintf("  %s %s\n", cursor, choice)
	}
	return s
}

func (m Model) Init() tea.Cmd {
	return nil
}
