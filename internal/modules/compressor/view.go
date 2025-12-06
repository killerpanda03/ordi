package compressor

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return m.filePicker.Init()
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filePicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if len(m.selectedFiles) == 0 {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filePicker.Styles.Selected.Render(m.selectedFiles[len(m.selectedFiles)-1]))
	}
	s.WriteString("\n\n" + m.filePicker.View() + "\n")
	return s.String()
}
