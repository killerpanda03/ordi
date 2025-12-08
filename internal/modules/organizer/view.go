package organizer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	categoryStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Bold(true)

	groupStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		MarginBottom(1)
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var b strings.Builder

	switch m.State {
	case stateInput:
		b.WriteString(titleStyle.Render("📂 Verzeichnis organisieren"))
		b.WriteString("\n\n")
		b.WriteString("Geben Sie den Pfad zum Ordner ein:\n\n")
		b.WriteString(m.TextInput.View())
		b.WriteString("\n\n")
		if m.Err != nil {
			b.WriteString(errorStyle.Render(fmt.Sprintf("⚠️  %v", m.Err)))
			b.WriteString("\n\n")
		}
		b.WriteString(helpStyle.Render("Enter = Scannen • Esc = Zurück zum Menü"))

	case stateScanning:
		b.WriteString(titleStyle.Render("🔍 Scanne Dateien..."))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s Durchsuche Verzeichnis: %s\n", m.Spinner.View(), m.Path))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Bitte warten..."))

	case statePreview:
		b.WriteString(titleStyle.Render("📋 Vorschau"))
		b.WriteString("\n\n")

		// Count files per category
		categoryCount := make(map[string]int)
		for _, file := range m.files {
			categoryCount[file.Category]++
		}

		b.WriteString(infoStyle.Render(fmt.Sprintf("Gefundene Dateien: %d\n", m.totalFiles)))
		b.WriteString("\n")

		// Show category breakdown in fixed order
		stats := []string{"Kategorien:"}
		categoryOrder := []string{"Bilder", "Videos", "Musik", "Dokumente", "Archive", "Sonstiges"}
		for _, category := range categoryOrder {
			if count, exists := categoryCount[category]; exists {
				icon := getCategoryIcon(category)
				stats = append(stats, fmt.Sprintf("  %s %-12s %d Dateien", icon, category+":", count))
			}
		}
		b.WriteString(categoryStyle.Render(lipgloss.JoinVertical(lipgloss.Left, stats...)))
		b.WriteString("\n\n")

		// Show first 10 files
		if len(m.files) > 0 {
			maxShow := 10
			if len(m.files) < maxShow {
				maxShow = len(m.files)
			}

			var samples []string
			for i := 0; i < maxShow; i++ {
				file := m.files[i]
				samples = append(samples, fmt.Sprintf("  %s %s → %s", file.Icon, truncate(file.Name, 40), file.Category))
			}

			if len(m.files) > maxShow {
				samples = append(samples, fmt.Sprintf("  ... und %d weitere", len(m.files)-maxShow))
			}

			b.WriteString(groupStyle.Render(lipgloss.JoinVertical(lipgloss.Left, samples...)))
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Enter = Organisieren starten • Esc = Abbrechen"))

	case stateOrganizing:
		b.WriteString(titleStyle.Render("📦 Organisiere Dateien..."))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s Verschiebe Dateien in Kategorien...\n", m.Spinner.View()))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Bitte warten..."))

	case stateFinished:
		if m.Err != nil {
			b.WriteString(titleStyle.Render("❌ Fehler"))
			b.WriteString("\n\n")
			b.WriteString(errorStyle.Render(fmt.Sprintf("Fehler: %v", m.Err)))
		} else {
			b.WriteString(titleStyle.Render("✅ Organisation abgeschlossen"))
			b.WriteString("\n\n")
			b.WriteString(successStyle.Render(fmt.Sprintf("✓ %d Dateien erfolgreich organisiert!", m.stats.TotalMoved)))
			b.WriteString("\n\n")

			if len(m.stats.Categories) > 0 {
				stats := []string{"Dateien pro Kategorie:"}
				categoryOrder := []string{"Bilder", "Videos", "Musik", "Dokumente", "Archive", "Sonstiges"}
				for _, category := range categoryOrder {
					if count, exists := m.stats.Categories[category]; exists {
						icon := getCategoryIcon(category)
						stats = append(stats, fmt.Sprintf("  %s %-12s %d Dateien", icon, category+":", count))
					}
				}
				b.WriteString(infoStyle.Render(lipgloss.JoinVertical(lipgloss.Left, stats...)))
			}
		}
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter = Zurück zum Menü"))
	}

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
