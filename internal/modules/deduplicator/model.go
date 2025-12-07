package deduplicator

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type ScanCompleteMsg struct {
	Files []string
	Err   error
}

type HashProgressMsg struct {
	Current int
	Total   int
}

type HashCompleteMsg struct {
	Duplicates      []DuplicateGroup
	SimilarImages   []SimilarGroup
	TotalSize       int64
	DuplicateSize   int64
	Err             error
}

type DeleteCompleteMsg struct {
	DeletedCount int
	FreedSpace   int64
	Err          error
}

type BackMsg struct{}

type state int

const (
	stateInput state = iota
	stateScanning
	stateHashing
	stateResults
	stateSelection
	stateDeleting
	stateFinished
)

type DuplicateGroup struct {
	Hash  string
	Files []FileInfo
	Size  int64
}

type SimilarGroup struct {
	Files      []FileInfo
	Similarity float64 // 0-100%
}

type FileInfo struct {
	Path     string
	Size     int64
	Selected bool // For deletion
}

type Model struct {
	textInput     textinput.Model
	spinner       spinner.Model
	progress      progress.Model
	state         state
	err           error

	// Scanning state
	dirPath       string
	scannedFiles  []string

	// Hashing state
	hashProgress  int
	hashTotal     int

	// Results state
	duplicates    []DuplicateGroup
	similarImages []SimilarGroup
	totalSize     int64
	duplicateSize int64
	savingsSize   int64

	// Selection state
	cursor        int
	selectedGroup int

	// UI state
	width         int
	height        int
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Pfad zum Ordner eingeben (z.B. C:\\Downloads)"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		textInput: ti,
		spinner:   s,
		progress:  p,
		state:     stateInput,
	}
}
