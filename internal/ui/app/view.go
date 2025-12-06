package app

func (m Model) View() string {
	s := "\n Ordi - Der File Organizer \n"

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
