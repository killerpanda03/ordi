package compressor

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
)

type clearErrorMsg struct{}

type Model struct {
	filePicker    filepicker.Model
	selectedFiles []string
	targetArchive string
	status        string
	quitting      bool
	err           error
}

func New() Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".png", ".jpg", ".mp4", ".mp3", ".mov", ".avi", ".mkv"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	return Model{
		filePicker:    fp,
		selectedFiles: []string{},
	}
}
