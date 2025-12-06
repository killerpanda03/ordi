package organizer

import (
	"example/ordi/internal/modules/organizer/styles"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type ProcessSuccessMsg struct{ Path string }

type ProcessFinishedMsg struct{}

type ProcessErrorMsg struct{ Err error }

type BackMsg struct{}

type state int

const (
	stateInput state = iota
	stateProcessing
	stateFinished
)

type Model struct {
	TextInput textinput.Model
	Spinner   spinner.Model
	styles    styles.Styles
	State     state
	Err       error
	Path      string
	Result    string
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "/pfad/zum/ordner"
	ti.CharLimit = 156

	return Model{
		TextInput: ti,
		Spinner:   s,
		styles:    *styles.DefaulStyles(),
		State:     stateInput,
	}
}
