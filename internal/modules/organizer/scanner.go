package organizer

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// scanFiles scans directory and returns preview of files with categories
func scanFiles(dirPath string) tea.Cmd {
	return func() tea.Msg {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return ScanCompleteMsg{Err: err}
		}

		var previews []FilePreview
		totalFiles := 0

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			info, err := file.Info()
			if err != nil {
				continue
			}

			category := getCategory(file.Name())
			previews = append(previews, FilePreview{
				Name:     file.Name(),
				Category: category.Name,
				Icon:     category.Icon,
				Size:     info.Size(),
			})
			totalFiles++
		}

		return ScanCompleteMsg{
			Files:      previews,
			TotalFiles: totalFiles,
		}
	}
}

// organizeFiles organizes files and reports progress
func organizeFiles(dirPath string) tea.Cmd {
	return func() tea.Msg {
		stats, err := Organize(dirPath)
		if err != nil {
			return OrganizeCompleteMsg{Err: err}
		}

		return OrganizeCompleteMsg{
			Stats: stats,
		}
	}
}
