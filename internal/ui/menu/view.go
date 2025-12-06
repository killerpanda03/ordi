package menu

import (
	"example/ordi/internal/ui/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) View() string {
	s := "\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = styles.Cursor.Render(choice)
		}
		s += fmt.Sprintf("  %s %s\n", cursor, choice)
	}
	return s
}

func (m Model) Init() tea.Cmd {
	return nil
}
