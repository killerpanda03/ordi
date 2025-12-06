package organizer

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func StartOrganizing(path string) tea.Cmd {
	return func() tea.Msg {
		// Die blockierende Funktion Organize wird hier aufgerufen
		time.Sleep(time.Millisecond * 500)
		err := Organize(path)

		if err != nil {
			return ProcessErrorMsg{Err: err}
		}
		// Bei Erfolg geben wir den Pfad zurück
		return ProcessSuccessMsg{Path: path}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.State {
	case stateInput:
		return m.updateInput(msg) // Initial State
	case stateProcessing:
		return m.updateProcessing(msg) // Loading State
	case stateFinished:
		return m.updateFinished(msg) // Success/Error State
	}

	return m, cmd
}

func (m Model) updateInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Path = m.TextInput.Value()

			if m.Path == "" {
				m.Err = fmt.Errorf("Bitte gib einen Pfad ein.") // Angepasste Fehlermeldung
				return m, nil
			}

			if _, err := os.Stat(m.Path); os.IsNotExist(err) {
				m.Err = fmt.Errorf("Pfad existiert nicht: %v", m.Path)
				return m, nil
			}

			m.State = stateProcessing
			m.Err = nil

			return m, tea.Batch(m.Spinner.Tick, StartOrganizing(m.Path))

		case tea.KeyEsc:
			return m, func() tea.Msg { return "" }
		}
	}
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) updateProcessing(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case ProcessSuccessMsg:
		m.State = stateFinished
		return m, nil

	case ProcessErrorMsg:
		m.State = stateFinished
		m.Err = msg.Err
		return m, nil
	}
	return m, nil
}

func (m Model) updateFinished(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			m.State = stateInput
			m.Result = ""
			m.Err = nil
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	return m, nil
}
