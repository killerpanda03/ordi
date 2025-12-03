package main

import (
	"example/ordi/internal/organizer"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateInput sessionState = iota // 0: Eingabe des Pfades
	stateMenu                      // 1: Das Auswahlmenü
	test
)

type model struct {
	state     sessionState    // In welcher Phase sind wir?
	textInput textinput.Model // Das Textfeld-Objekt
	spinner   spinner.Model
	path      string   // Der gespeicherte Pfad
	choices   []string // Menüpunkte
	cursor    int      // Menü-Cursor
	status    string   // Status-Nachricht
	err       error    // Falls der Pfad ungültig ist
}

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

func (m model) Init() tea.Cmd {
	// Das Textfeld braucht einen "Blink"-Befehl für den Cursor beim Start
	if m.state == stateInput {
		return textinput.Blink
	}

	return m.spinner.Tick
}

func initialModel() model {

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// Wir konfigurieren das Textfeld
	ti := textinput.New()
	ti.Placeholder = keywordStyle.Render("/pfad/zum/ordner")
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		state:     stateMenu, // Wir starten im Eingabe-Modus!
		textInput: ti,
		spinner:   s,
		choices:   []string{"📂 Ein Verzeichnis organisieren", "Bilder sortieren", "Bild komprimieren", "❌ Beenden"},
		cursor:    0,
		status:    "",
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	// Logik-Weiche: Je nach Zustand reagieren wir anders
	switch m.state {

	// === ZUSTAND 1: PFAD EINGEBEN ===
	case stateInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// Pfad übernehmen
				inputPath := m.textInput.Value()

				// Validierung: Existiert der Pfad überhaupt?
				if _, err := os.Stat(inputPath); os.IsNotExist(err) {
					m.err = fmt.Errorf("Pfad existiert nicht: %s", inputPath)
					return m, nil
				}

				// Alles gut -> Zustand wechseln und Pfad speichern
				m.path = inputPath
				m.state = stateMenu // Wechsel zum Menü!
				m.status = "Bereit für Ordner: " + m.path
				m.err = nil
				return m, nil
			}
		}
		// Textfeld aktualisieren (Buchstaben annehmen, Blinken etc.)
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd

	// === ZUSTAND 2: MENÜ BEDIENEN ===
	case stateMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q": // Im Menü darf man mit q beenden
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				if m.cursor == 0 {
					// Organisieren ausführen
					err := organizer.Organize(m.path)
					if err != nil {
						m.status = "⚠️  Fehler: " + err.Error()
					} else {
						m.status = "✅  Erfolg! Verzeichnis organisiert."
					}
				} else if m.cursor == 1 {
					m.state = test
					return m, m.spinner.Tick
				} else {
					return m, tea.Quit
				}
			}
		}
	case test:
		var cmd tea.Cmd

		// Prüfen, ob wir abbrechen wollen
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" || msg.String() == "esc" {
				m.state = stateMenu
				return m, nil
			}
		}

		// Spinner updaten
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	s := "\n  ✨ Ordi - Der File Organizer ✨\n\n"

	// Ansicht-Weiche: Was zeigen wir an?
	if m.state == stateInput {
		s += subtleStyle.Render("  Bitte gib den Pfad zum Ordner ein:\n\n")
		s += fmt.Sprintf("  %s\n\n", m.textInput.View()) // Das Textfeld rendern

		if m.err != nil {
			s += fmt.Sprintf("  ❌ %v\n", m.err)
		}
		s += "  (Drücke Enter zum Bestätigen, Ctrl+C zum Beenden)\n"

	} else if m.state == test {
		s += fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
		s += "  Bilder sortieren - Funktion noch in Arbeit!\n\n"
	} else {
		// Menü-Ansicht (wie vorher)
		// s += fmt.Sprintf("  📂 Aktueller Ordner: %s\n\n", m.path)

		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				choice = fmt.Sprintf("[%s]", choice)
			}
			s += fmt.Sprintf("  %s %s\n", cursor, choice)
		}

		s += "\n  --------------------------------\n"
		s += fmt.Sprintf("  Status: %s\n", m.status)
	}

	return s
}

func main() {
	// Wir brauchen keine CLI Args mehr!
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fehler: %v", err)
		os.Exit(1)
	}
}
