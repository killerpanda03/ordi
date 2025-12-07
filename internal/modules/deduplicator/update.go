package deduplicator

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateInput:
			switch msg.String() {
			case "enter":
				if m.textInput.Value() != "" {
					m.dirPath = m.textInput.Value()
					m.state = stateScanning
					return m, tea.Batch(
						m.spinner.Tick,
						scanDirectory(m.dirPath),
					)
				}
			case "esc":
				return m, func() tea.Msg { return BackMsg{} }
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

		case stateResults:
			switch msg.String() {
			case "enter":
				if len(m.duplicates) > 0 {
					m.state = stateSelection
					m.cursor = 0
					// Auto-select all duplicates except first in each group
					for i := range m.duplicates {
						for j := 1; j < len(m.duplicates[i].Files); j++ {
							m.duplicates[i].Files[j].Selected = true
						}
					}
					return m, nil
				}
				return m, func() tea.Msg { return BackMsg{} }
			case "esc":
				return m, func() tea.Msg { return BackMsg{} }
			}

		case stateSelection:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				maxCursor := 0
				for _, group := range m.duplicates {
					maxCursor += len(group.Files) - 1 // -1 because first file is not selectable
				}
				if m.cursor < maxCursor-1 {
					m.cursor++
				}
			case " ": // Space to toggle selection
				currentItem := 0
				for groupIdx := range m.duplicates {
					for fileIdx := 1; fileIdx < len(m.duplicates[groupIdx].Files); fileIdx++ {
						if currentItem == m.cursor {
							m.duplicates[groupIdx].Files[fileIdx].Selected = !m.duplicates[groupIdx].Files[fileIdx].Selected
							return m, nil
						}
						currentItem++
					}
				}
			case "esc":
				m.state = stateResults
				return m, nil
			case "enter":
				m.state = stateDeleting
				return m, tea.Batch(
					m.spinner.Tick,
					deleteDuplicates(m.duplicates),
				)
			}

		case stateFinished:
			if msg.String() == "enter" || msg.String() == "esc" {
				return m, func() tea.Msg { return BackMsg{} }
			}
		}

	case ScanCompleteMsg:
		if msg.Err != nil {
			m.err = fmt.Errorf("Fehler beim Scannen: %w", msg.Err)
			m.state = stateFinished
			return m, nil
		}
		m.scannedFiles = msg.Files
		m.hashTotal = len(msg.Files)
		m.hashProgress = 0
		m.state = stateHashing

		return m, findDuplicates(msg.Files)

	case HashProgressMsg:
		m.hashProgress = msg.Current
		return m, nil

	case HashCompleteMsg:
		if msg.Err != nil {
			m.err = fmt.Errorf("Fehler beim Hashen: %w", msg.Err)
			m.state = stateFinished
			return m, nil
		}
		m.duplicates = msg.Duplicates
		m.similarImages = msg.SimilarImages
		m.totalSize = msg.TotalSize
		m.duplicateSize = msg.DuplicateSize
		m.state = stateResults
		return m, nil

	case DeleteCompleteMsg:
		if msg.Err != nil {
			m.err = fmt.Errorf("Fehler beim Löschen: %w", msg.Err)
		} else {
			m.savingsSize = msg.FreedSpace
		}
		m.state = stateFinished
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
