package app

import (
	"example/ordi/internal/modules/organizer"
	"example/ordi/internal/ui/menu"

	tea "github.com/charmbracelet/bubbletea"
)

type ActiveModule int

const (
	Menu ActiveModule = iota
	Organizer
	Duplicates
	Images
)

type sessionState int

const (
	stateMenu sessionState = iota
	stateOrganize
	stateImageSort
)

type Model struct {
	state     sessionState
	menu      menu.Model
	organizer organizer.Model
}

func (m Model) Init() tea.Cmd {
	return m.organizer.Init()
}

func New() Model {
	return Model{
		state:     stateMenu,
		menu:      menu.New(),
		organizer: organizer.New(),
	}
}
