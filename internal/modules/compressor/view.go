package compressor

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	successStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	warningStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	infoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	subtleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkMarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	crossMarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.filePicker.Init(), m.spinner.Tick)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	switch m.state {
	case stateToolCheck:
		s.WriteString(titleStyle.Render("📦 Datei-Komprimierung") + "\n\n")

		if m.err != nil {
			s.WriteString(errorStyle.Render("❌ " + m.err.Error()) + "\n\n")
		}

		s.WriteString(m.renderToolStatus())
		s.WriteString("\n" + subtleStyle.Render("Drücke Enter um fortzufahren oder q zum Beenden"))

	case stateFileSelection:
		s.WriteString(titleStyle.Render("📁 Dateien auswählen") + "\n\n")

		
		s.WriteString(subtleStyle.Render("Verzeichnis: "+m.filePicker.CurrentDirectory) + "\n")

		
		formats := m.tools.SupportedFormats()
		if len(formats) > 0 {
			s.WriteString(infoStyle.Render("Unterstützte Formate: ") + subtleStyle.Render(strings.Join(formats, ", ")) + "\n\n")
		} else {
			s.WriteString(warningStyle.Render("⚠ Keine Tools verfügbar - bitte installiere mindestens ein Tool!") + "\n\n")
		}

		if m.err != nil {
			s.WriteString(errorStyle.Render(m.err.Error()) + "\n\n")
		}
		if len(m.selectedFiles) > 0 {
			s.WriteString(successStyle.Render(fmt.Sprintf("✓ %d Datei(en) ausgewählt", len(m.selectedFiles))) + "\n\n")
		}
		s.WriteString(m.filePicker.View() + "\n")
		s.WriteString(subtleStyle.Render("\nEnter: Auswählen | Tab: Mehrfachauswahl | Esc: Weiter | q: Zurück"))

	case stateOutputSelection:
		s.WriteString(titleStyle.Render("📂 Ausgabeverzeichnis") + "\n\n")
		s.WriteString(infoStyle.Render(fmt.Sprintf("%d Datei(en) werden komprimiert", len(m.selectedFiles))) + "\n\n")
		s.WriteString(m.outputInput.View() + "\n\n")
		s.WriteString(subtleStyle.Render("Enter: Komprimierung starten | q: Zurück"))

	case stateProcessing:
		s.WriteString(titleStyle.Render("⚙️  Komprimiere Dateien...") + "\n\n")
		s.WriteString(m.spinner.View() + " " + m.status + "\n\n")
		s.WriteString(infoStyle.Render(fmt.Sprintf("Fortschritt: %d/%d", m.success+m.failed, len(m.selectedFiles))))

	case stateFinished:
		s.WriteString(titleStyle.Render("✅ Komprimierung abgeschlossen") + "\n\n")
		if m.err != nil {
			s.WriteString(errorStyle.Render("❌ Fehler: "+m.err.Error()) + "\n\n")
		} else {
			s.WriteString(successStyle.Render(fmt.Sprintf("✓ %d Datei(en) erfolgreich komprimiert", m.success)) + "\n")
			if m.failed > 0 {
				s.WriteString(warningStyle.Render(fmt.Sprintf("⚠ %d Datei(en) fehlgeschlagen", m.failed)) + "\n")
			}
			if m.originalSize > 0 && m.compressedSize > 0 {
				saved := m.originalSize - m.compressedSize
				percent := float64(saved) / float64(m.originalSize) * 100
				s.WriteString(infoStyle.Render(fmt.Sprintf("\n💾 Größe reduziert um %.1f%% (%s gespart)",
					percent, formatBytes(saved))) + "\n")
			}
		}
		s.WriteString("\n" + subtleStyle.Render("Drücke Enter um zum Menü zurückzukehren"))
	}

	return s.String()
}

func (m Model) renderToolStatus() string {
	var s strings.Builder

	s.WriteString(infoStyle.Render("Komprimierungs-Tools (alle optional):") + "\n\n")

	tools := []struct {
		name      string
		available bool
		desc      string
		formats   string
	}{
		{"ffmpeg", m.tools.FFmpeg, "Videos & Audio", ".mp4, .mov, .avi, .mkv, .mp3, .wav"},
		{"ImageMagick", m.tools.ImageMagick, "Bilder", ".jpg, .png, .gif, .bmp"},
		{"Ghostscript", m.tools.Ghostscript, "PDF-Dokumente", ".pdf"},
		{"7zip", m.tools.SevenZip, "Office & Archive", ".docx, .xlsx, .pptx, .zip"},
	}

	availableCount := 0
	for _, tool := range tools {
		var mark, status, formatInfo string
		if tool.available {
			mark = checkMarkStyle.Render("✓")
			status = successStyle.Render(tool.name)
			formatInfo = subtleStyle.Render(" → " + tool.formats)
			availableCount++
		} else {
			mark = crossMarkStyle.Render("✗")
			status = subtleStyle.Render(tool.name)
			formatInfo = subtleStyle.Render(" (nicht verfügbar)")
		}
		s.WriteString(fmt.Sprintf("  %s %s - %s%s\n", mark, status, tool.desc, formatInfo))
	}

	s.WriteString("\n")

	if !m.tools.HasAnyTool() {
		s.WriteString(errorStyle.Render("⚠ Keine Komprimierungs-Tools gefunden!") + "\n")
		s.WriteString(subtleStyle.Render("Installiere mindestens ein Tool um fortzufahren.") + "\n")
		s.WriteString(m.renderInstallHelp())
	} else {
		s.WriteString(successStyle.Render(fmt.Sprintf("✓ %d von 4 Tools verfügbar - bereit zum Komprimieren!", availableCount)) + "\n")

		missingCount := 4 - availableCount
		if missingCount > 0 {
			s.WriteString(subtleStyle.Render(fmt.Sprintf("\nOptional: Installiere %d weitere Tool(s) für mehr Dateitypen:", missingCount)) + "\n")
			s.WriteString(m.renderInstallHelp())
		}
	}

	return s.String()
}

func (m Model) renderInstallHelp() string {
	var s strings.Builder

	s.WriteString("\n" + infoStyle.Render("Installation:") + "\n")
	s.WriteString(subtleStyle.Render("Windows (choco):") + "\n")
	s.WriteString("  choco install ffmpeg imagemagick ghostscript 7zip\n")
	s.WriteString(subtleStyle.Render("Linux (apt):") + "\n")
	s.WriteString("  sudo apt-get install ffmpeg imagemagick ghostscript p7zip-full\n")

	return s.String()
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
