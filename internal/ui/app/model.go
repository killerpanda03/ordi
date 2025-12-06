package app

import (
	"example/ordi/internal/modules/compressor"
	"example/ordi/internal/modules/organizer"
	"example/ordi/internal/ui/menu"

	tea "github.com/charmbracelet/bubbletea"
)

type ActiveModule int

const (
	Menu ActiveModule = iota
	Organizer
	Compresser
)

type sessionState int

const (
	stateMenu sessionState = iota
	stateOrganize
	stateCompress
)

type Model struct {
	state      sessionState
	menu       menu.Model
	organizer  organizer.Model
	compressor compressor.Model
}

func (m Model) Init() tea.Cmd {
	return m.organizer.Init()
}

func New() Model {
	return Model{
		state:      stateMenu,
		menu:       menu.New(),
		organizer:  organizer.New(),
		compressor: compressor.New(),
	}
}
