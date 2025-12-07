package deduplicator

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	groupStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)
)

func (m Model) View() string {
	var b strings.Builder

	switch m.state {
	case stateInput:
		b.WriteString(titleStyle.Render("Duplikate finden"))
		b.WriteString("\n\n")
		b.WriteString("Geben Sie den Pfad zum Ordner ein, der durchsucht werden soll:\n\n")
		b.WriteString(m.textInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter = Scannen starten • Esc = Zurück zum Menü"))

	case stateScanning:
		b.WriteString(titleStyle.Render("Scanne Dateien..."))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s Durchsuche Verzeichnis: %s\n", m.spinner.View(), m.dirPath))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Bitte warten..."))

	case stateHashing:
		b.WriteString(titleStyle.Render("Berechne Hashes..."))
		b.WriteString("\n\n")
		percent := float64(m.hashProgress) / float64(m.hashTotal)
		b.WriteString(m.progress.ViewAs(percent))
		b.WriteString(fmt.Sprintf("\n\n%d / %d Dateien verarbeitet", m.hashProgress, m.hashTotal))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Dies kann einige Minuten dauern..."))

	case stateResults:
		b.WriteString(titleStyle.Render("📊 Scan-Ergebnisse"))
		b.WriteString("\n\n")

		if len(m.duplicates) == 0 && len(m.similarImages) == 0 {
			b.WriteString(successStyle.Render("✓ Keine Duplikate gefunden!"))
			b.WriteString("\n\n")
			b.WriteString(fmt.Sprintf("Gescannte Dateien: %d\n", len(m.scannedFiles)))
		} else {
			stats := []string{
				fmt.Sprintf("Gescannte Dateien:       %d", len(m.scannedFiles)),
				fmt.Sprintf("Exakte Duplikate:        %d Gruppen", len(m.duplicates)),
				fmt.Sprintf("Ähnliche Bilder:         %d Gruppen", len(m.similarImages)),
				fmt.Sprintf("Verschwendeter Speicher: %s", formatBytes(m.duplicateSize)),
			}
			b.WriteString(infoStyle.Render(lipgloss.JoinVertical(lipgloss.Left, stats...)))
			b.WriteString("\n\n")

			// Show first few duplicate groups
			shown := 0
			maxShow := 3
			for i, group := range m.duplicates {
				if shown >= maxShow {
					remaining := len(m.duplicates) - shown
					b.WriteString(fmt.Sprintf("\n... und %d weitere exakte Duplikat-Gruppen\n", remaining))
					break
				}

				groupContent := fmt.Sprintf("🔄 Exakte Duplikate - Gruppe %d (%s pro Datei)\n", i+1, formatBytes(group.Size))
				for j, file := range group.Files {
					if j >= 3 {
						groupContent += fmt.Sprintf("  ... und %d weitere\n", len(group.Files)-3)
						break
					}
					groupContent += fmt.Sprintf("  • %s\n", truncatePath(file.Path, 70))
				}
				b.WriteString(groupStyle.Render(groupContent))
				shown++
			}

			// Show similar images
			if len(m.similarImages) > 0 {
				b.WriteString("\n")
				shownSimilar := 0
				maxShowSimilar := 2
				for i, group := range m.similarImages {
					if shownSimilar >= maxShowSimilar {
						remaining := len(m.similarImages) - shownSimilar
						b.WriteString(fmt.Sprintf("\n... und %d weitere ähnliche Bild-Gruppen\n", remaining))
						break
					}

					groupContent := fmt.Sprintf("🖼️  Ähnliche Bilder - Gruppe %d (%.1f%% ähnlich)\n", i+1, group.Similarity)
					for j, file := range group.Files {
						if j >= 3 {
							groupContent += fmt.Sprintf("  ... und %d weitere\n", len(group.Files)-3)
							break
						}
						groupContent += fmt.Sprintf("  • %s (%s)\n", truncatePath(file.Path, 60), formatBytes(file.Size))
					}
					b.WriteString(groupStyle.Render(groupContent))
					shownSimilar++
				}
			}
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Enter = Bereinigung starten • Esc = Zurück zum Menü"))

	case stateSelection:
		b.WriteString(titleStyle.Render("🗑️  Duplikate zur Löschung auswählen"))
		b.WriteString("\n\n")

		totalToDelete := 0
		sizeToFree := int64(0)

		// Show all duplicate groups with selection checkboxes
		currentItem := 0
		for groupIdx, group := range m.duplicates {
			b.WriteString(fmt.Sprintf("\n📁 Gruppe %d - %s pro Datei:\n", groupIdx+1, formatBytes(group.Size)))

			// First file is always kept (not selectable)
			b.WriteString(fmt.Sprintf("  [KEEP] %s\n", truncatePath(group.Files[0].Path, 70)))

			// Other files can be selected for deletion
			for fileIdx := 1; fileIdx < len(group.Files); fileIdx++ {
				file := group.Files[fileIdx]
				checkbox := "[ ]"
				style := lipgloss.NewStyle()

				if file.Selected {
					checkbox = "[✓]"
					totalToDelete++
					sizeToFree += file.Size
					style = selectedStyle
				}

				cursor := "  "
				if currentItem == m.cursor {
					cursor = "> "
					style = style.Foreground(lipgloss.Color("205"))
				}

				b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, style.Render(truncatePath(file.Path, 65))))
				currentItem++
			}
		}

		b.WriteString("\n")
		b.WriteString(infoStyle.Render(fmt.Sprintf("📊 %d Dateien ausgewählt • %s werden freigegeben", totalToDelete, formatBytes(sizeToFree))))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("↑/↓ = Navigieren • Space = Auswählen/Abwählen • Enter = Löschen • Esc = Abbrechen"))

	case stateDeleting:
		b.WriteString(titleStyle.Render("🗑️  Lösche Duplikate..."))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s Räume auf...\n", m.spinner.View()))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Bitte warten..."))

	case stateFinished:
		b.WriteString(titleStyle.Render("✓ Bereinigung abgeschlossen"))
		b.WriteString("\n\n")
		if m.err != nil {
			b.WriteString(errorStyle.Render(fmt.Sprintf("Fehler: %v", m.err)))
		} else {
			b.WriteString(successStyle.Render(fmt.Sprintf("✓ %s Speicher freigegeben!", formatBytes(m.savingsSize))))
		}
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter = Zurück zum Menü"))
	}

	if m.err != nil && m.state != stateFinished {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf("Fehler: %v", m.err)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter = Zurück zum Menü"))
	}

	return b.String()
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

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}
