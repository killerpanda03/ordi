package organizer

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.State {
		case stateInput:
			switch msg.Type {
			case tea.KeyEnter:
				m.Path = m.TextInput.Value()

				if m.Path == "" {
					m.Err = fmt.Errorf("Bitte gib einen Pfad ein.")
					return m, nil
				}

				if _, err := os.Stat(m.Path); os.IsNotExist(err) {
					m.Err = fmt.Errorf("Pfad existiert nicht: %v", m.Path)
					return m, nil
				}

				m.State = stateScanning
				m.Err = nil
				return m, tea.Batch(m.Spinner.Tick, scanFiles(m.Path))

			case tea.KeyEsc:
				return m, func() tea.Msg { return BackMsg{} }

			default:
				m.TextInput, cmd = m.TextInput.Update(msg)
				return m, cmd
			}

		case statePreview:
			switch msg.String() {
			case "enter":
				m.State = stateOrganizing
				return m, tea.Batch(m.Spinner.Tick, organizeFiles(m.Path))
			case "esc":
				return m, func() tea.Msg { return BackMsg{} }
			}

		case stateFinished:
			if msg.String() == "enter" || msg.String() == "esc" {
				return m, func() tea.Msg { return BackMsg{} }
			}
		}

	case ScanCompleteMsg:
		if msg.Err != nil {
			m.Err = fmt.Errorf("Fehler beim Scannen: %w", msg.Err)
			m.State = stateFinished
			return m, nil
		}
		m.files = msg.Files
		m.totalFiles = msg.TotalFiles
		m.State = statePreview
		return m, nil

	case OrganizeCompleteMsg:
		if msg.Err != nil {
			m.Err = fmt.Errorf("Fehler beim Organisieren: %w", msg.Err)
		} else {
			m.stats = msg.Stats
		}
		m.State = stateFinished
		return m, nil

	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}
