package app

func (m Model) View() string {
	s := ""

	switch m.state {
	case stateMenu:
		s += m.menu.View()
	case stateOrganize:
		s += m.organizer.View()
	case stateCompress:
		s += m.compressor.View()
	case stateDeduplicate:
		s += m.deduplicator.View()
	}

	return s
}
