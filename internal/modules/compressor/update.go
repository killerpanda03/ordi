package compressor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateToolCheck:
			switch msg.String() {
			case "q", "esc":
				return m, func() tea.Msg { return BackMsg{} }
			case "enter":
				if !m.tools.HasAnyTool() {
					
					m.err = fmt.Errorf("Mindestens ein Tool wird benötigt. Installiere ffmpeg, ImageMagick, Ghostscript oder 7zip.")
					return m, clearErrorAfter(5 * time.Second)
				}
				
				m.state = stateFileSelection
				m.err = nil
				return m, nil
			}

		case stateFileSelection:
			switch msg.String() {
			case "q":
				return m, func() tea.Msg { return BackMsg{} }
			case "esc":
				if len(m.selectedFiles) > 0 {
					m.state = stateOutputSelection
					m.outputInput.Focus()
					return m, nil
				}
			case "tab":
				
				if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
					m.selectedFiles = append(m.selectedFiles, path)
				}
			}

		case stateOutputSelection:
			switch msg.String() {
			case "q":
				m.state = stateFileSelection
				m.selectedFiles = []string{}
				return m, nil
			case "enter":
				m.outputDir = m.outputInput.Value()
				m.state = stateProcessing
				return m, tea.Batch(m.spinner.Tick, compressFilesCmd(m.selectedFiles, m.outputDir))
			}

		case stateFinished:
			switch msg.String() {
			case "enter", "q":
				return m, func() tea.Msg { return BackMsg{} }
			}
		}

	case CompressCompleteMsg:
		m.state = stateFinished
		m.success = msg.Success
		m.failed = msg.Failed
		m.originalSize = msg.OriginalSize
		m.compressedSize = msg.CompressedSize
		m.err = msg.Err
		return m, nil

	case clearErrorMsg:
		m.err = nil

	case spinner.TickMsg:
		if m.state == stateProcessing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	
	var cmd tea.Cmd
	switch m.state {
	case stateFileSelection:
		m.filePicker, cmd = m.filePicker.Update(msg)

		
		if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
			
			if m.tools.IsFileSupported(path) {
				m.selectedFiles = append(m.selectedFiles, path)
				m.err = nil
			} else {
				ext := strings.ToLower(filepath.Ext(path))
				m.err = fmt.Errorf("%s wird nicht unterstützt. Benötigtes Tool für %s fehlt.",
					filepath.Base(path), ext)
				return m, clearErrorAfter(3 * time.Second)
			}
		}

		
		if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
			m.err = fmt.Errorf("%s ist nicht verfügbar.", filepath.Base(path))
			return m, clearErrorAfter(2 * time.Second)
		}

	case stateOutputSelection:
		m.outputInput, cmd = m.outputInput.Update(msg)
	}

	return m, cmd
}

func compressFilesCmd(files []string, outputDir string) tea.Cmd {
	return func() tea.Msg {
		success := 0
		failed := 0
		var originalSize, compressedSize int64

		for _, file := range files {
			
			var outputPath string
			if outputDir == "" {
				
				dir := filepath.Dir(file)
				base := filepath.Base(file)
				ext := filepath.Ext(base)
				name := base[:len(base)-len(ext)]
				outputPath = filepath.Join(dir, name+"_compressed"+ext)
			} else {
				base := filepath.Base(file)
				ext := filepath.Ext(base)
				name := base[:len(base)-len(ext)]
				outputPath = filepath.Join(outputDir, name+"_compressed"+ext)
			}

			
			if stat, err := os.Stat(file); err == nil {
				originalSize += stat.Size()
			}

			
			if err := compressFile(file, outputPath); err != nil {
				failed++
			} else {
				success++
				
				if stat, err := os.Stat(outputPath); err == nil {
					compressedSize += stat.Size()
				}
			}
		}

		return CompressCompleteMsg{
			Success:        success,
			Failed:         failed,
			OriginalSize:   originalSize,
			CompressedSize: compressedSize,
		}
	}
}
