package app

import (
	"example/ordi/internal/modules/organizer"
	"example/ordi/internal/ui/menu"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	switch m.state {
	case stateMenu:
		newMenu, newCmd := m.menu.Update(msg)
		m.menu = newMenu
		cmd = newCmd

		if selectMsg, ok := msg.(menu.SelectMsg); ok {
			switch selectMsg {
			case 0:
				m.state = stateOrganize

				m.organizer = organizer.New()
				return m, tea.Batch(m.organizer.Init(), m.organizer.TextInput.Focus())
			case 1:
			case 3:
				return m, tea.Quit
			}
		}

	case stateOrganize:
		newOrg, newCmd := m.organizer.Update(msg)
		m.organizer = newOrg
		cmd = newCmd

		if _, ok := msg.(organizer.BackMsg); ok {
			m.state = stateMenu
			m.organizer = organizer.New()
			return m, nil
		}
	}

	return m, cmd
}
