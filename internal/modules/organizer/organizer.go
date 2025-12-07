package organizer

import (
	"os"
	"path/filepath"
)

func Organize(dirPath string) (CategoryStats, error) {
	files, err := readFiles(dirPath)
	if err != nil {
		return CategoryStats{}, err
	}

	stats := CategoryStats{
		Categories: make(map[string]int),
		TotalMoved: 0,
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		category := getCategory(file.Name())
		err := createDir(dirPath, category.Name)
		if err != nil {
			return stats, err
		}
		srcPath := filepath.Join(dirPath, file.Name())
		destDir := filepath.Join(dirPath, category.Name)

		err = moveFile(srcPath, destDir)
		if err != nil {
			return stats, err
		}

		stats.Categories[category.Name]++
		stats.TotalMoved++
	}

	return stats, nil
}

type Category struct {
	Name string
	Icon string
}

var categories = map[string]Category{
	"images": {Name: "Bilder", Icon: "📷"},
	"videos": {Name: "Videos", Icon: "🎬"},
	"music":  {Name: "Musik", Icon: "🎵"},
	"docs":   {Name: "Dokumente", Icon: "📄"},
	"archives": {Name: "Archive", Icon: "📦"},
	"other":  {Name: "Sonstiges", Icon: "📁"},
}

func getCategory(fileName string) Category {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return categories["images"]
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv":
		return categories["videos"]
	case ".mp3", ".wav", ".flac", ".aac", ".ogg":
		return categories["music"]
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt":
		return categories["docs"]
	case ".zip", ".rar", ".tar", ".gz", ".7z":
		return categories["archives"]
	default:
		return categories["other"]
	}
}

func getCategoryIcon(categoryName string) string {
	for _, cat := range categories {
		if cat.Name == categoryName {
			return cat.Icon
		}
	}
	return "📁"
}

func readFiles(dirPath string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func createDir(dirPath, category string) error {
	categoryPath := filepath.Join(dirPath, category)
	return os.MkdirAll(categoryPath, os.ModePerm)
}

func moveFile(filePath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(filePath))
	return os.Rename(filePath, destPath)
}
