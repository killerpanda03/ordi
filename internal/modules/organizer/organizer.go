package organizer

import (
	"os"
	"path/filepath"
)

func Organize(dirPath string) error {
	files, err := readFiles(dirPath)
	if err != nil {

		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		category := getCategory(file.Name())
		err := createDir(dirPath, category)
		if err != nil {
			return err
		}
		srcPath := filepath.Join(dirPath, file.Name())
		destDir := filepath.Join(dirPath, category)

		err = moveFile(srcPath, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCategory(fileName string) string {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff":
		return "Bilder"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv":
		return "Videos"
	case ".mp3", ".wav", ".flac", ".aac":
		return "Musik"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return "Dokumente"
	case ".zip", ".rar", ".tar", ".gz":
		return "Archive"
	default:
		return "Sonstiges"
	}
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
