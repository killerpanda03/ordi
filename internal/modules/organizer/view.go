package organizer

import (
	"example/ordi/internal/ui/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// func (m Model) View() string {
// 	switch m.State {
// 	case stateInput:
// 		s := styles.Subtle.Render("  Bitte gib den Pfad zum Ordner ein:\n\n")
// 		s += fmt.Sprintf("%s\n\n", m.TextInput.View())
// 		if m.Err != nil {
// 			s += fmt.Sprintf("  ❌ %v\n", m.Err)
// 		}
// 		s += "  (Enter: Bestätigen, Esc: Zurück)\n"
// 		return s

//		case stateProcessing:
//			s := fmt.Sprintf("\n  %s Organisiere Dateien in %s...\n\n", m.Spinner.View(), styles.Keyword.Render(m.Path))
//			s += "  (Dies dauert einen Moment)\n"
//			return s
//		}
//		return ""
//	}
func (m Model) Init() tea.Cmd {

	return textinput.Blink
}

func (m Model) View() string {
	var viewContent string
	viewContent += "\n"

	switch m.State {
	case stateInput: // Initial State
		viewContent += lipgloss.JoinVertical(lipgloss.Left,
			"Bitte gib den Pfad zum Ordner ein:\n",
			m.styles.InputField.Render(m.TextInput.View()),
			"(Enter: Bestätigen, Esc: Zurück)")
		if m.Err != nil {
			viewContent += fmt.Sprintf("❌ %v\n\n", m.Err)
		}

	case stateProcessing: // Loading State
		viewContent += fmt.Sprintf("%s Organisiere Dateien in %s...\n\n", m.Spinner.View(), styles.Keyword.Render(m.Path))
		viewContent += "(Bitte warten...)\n"

	case stateFinished: // Success / Error State
		if m.Err != nil {
			// Fehler anzeigen
			viewContent += fmt.Sprintf("❌ Organisation in %s fehlgeschlagen!\n", styles.Keyword.Render(m.Path))
			viewContent += fmt.Sprintf("Fehler: %v\n\n", m.Err)
		} else {
			// Erfolg anzeigen
			viewContent += styles.Keyword.Render(m.Result) + "\n\n"
		}
		viewContent += "(Drücke Enter oder Esc, um zum Menü zurückzukehren)"
	}

	return styles.Main.Render(viewContent) + "\n"
}
