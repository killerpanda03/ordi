package compressor

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type clearErrorMsg struct{}

type BackMsg struct{}

type CompressCompleteMsg struct {
	Success        int
	Failed         int
	OriginalSize   int64
	CompressedSize int64
	Err            error
}

type state int

const (
	stateToolCheck state = iota
	stateFileSelection
	stateOutputSelection
	stateProcessing
	stateFinished
)

type Model struct {
	filePicker    filepicker.Model
	outputInput   textinput.Model
	spinner       spinner.Model
	state         state
	selectedFiles []string
	outputDir     string
	status        string
	quitting      bool
	err           error

	
	tools ToolAvailability

	
	success        int
	failed         int
	originalSize   int64
	compressedSize int64
}

func New() Model {
	
	tools := checkExternalTools()

	
	fp := filepicker.New()

	
	
	

	
	
	if cwd, err := os.Getwd(); err == nil {
		fp.CurrentDirectory = cwd
	} else {
		
		fp.CurrentDirectory, _ = os.UserHomeDir()
	}
	fp.ShowHidden = false
	fp.DirAllowed = true  
	fp.FileAllowed = true 

	
	ti := textinput.New()
	ti.Placeholder = "Ausgabeverzeichnis (leer = gleiches Verzeichnis)"
	ti.CharLimit = 256
	ti.Width = 80

	
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		filePicker:    fp,
		outputInput:   ti,
		spinner:       s,
		state:         stateToolCheck,
		selectedFiles: []string{},
		tools:         tools,
	}
}
