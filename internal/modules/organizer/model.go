package organizer

import (
	"example/ordi/internal/modules/organizer/styles"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type ScanCompleteMsg struct {
	Files      []FilePreview
	TotalFiles int
	Err        error
}

type OrganizeProgressMsg struct {
	Current int
	Total   int
}

type OrganizeCompleteMsg struct {
	Stats CategoryStats
	Err   error
}

type ProcessSuccessMsg struct{ Path string }

type ProcessFinishedMsg struct{}

type ProcessErrorMsg struct{ Err error }

type BackMsg struct{}

type state int

const (
	stateInput state = iota
	stateScanning
	statePreview
	stateOrganizing
	stateFinished
)

type FilePreview struct {
	Name     string
	Category string
	Icon     string
	Size     int64
}

type CategoryStats struct {
	Categories map[string]int
	TotalMoved int
}

type Model struct {
	TextInput textinput.Model
	Spinner   spinner.Model
	styles    styles.Styles
	State     state
	Err       error
	Path      string
	Result    string

	// Preview state
	files      []FilePreview
	totalFiles int

	// Progress state
	progress int
	total    int

	// Results
	stats CategoryStats
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "Pfad zum Ordner eingeben"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80

	return Model{
		TextInput: ti,
		Spinner:   s,
		styles:    *styles.DefaulStyles(),
		State:     stateInput,
	}
}
