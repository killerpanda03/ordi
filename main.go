package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count int
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return nil
}

// Initialer Zustand
func initialModel() model {
	return model{count: 0}
}

// Update: Logik bei Tastendruck
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			m.count++

		case "down":
			m.count--
		}
	}

	return m, nil
}

// View: Was im Terminal angezeigt wird
func (m model) View() string {
	return fmt.Sprintf(
		"Zähler: %d\n\n↑ erhöhen | ↓ verringern | q beenden\n",
		m.count,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Fehler:", err)
		os.Exit(1)
	}
}
