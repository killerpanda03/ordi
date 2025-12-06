package app

func (m Model) View() string {
	s := ""

	switch m.state {
	case stateMenu:
		s += m.menu.View()
	case stateOrganize:
		s += m.organizer.View()
	case stateImageSort:
		s += "Feature kommt bald..."
	}

	return s
}
