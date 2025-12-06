package compressor

import (
	"fmt"
	"os"
)

func createMediaArchive(files []string, targetPath string) error {
	archiveFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Archivdatei: %w", err)
	}
	defer archiveFile.Close()

	//TODO: Implementiere die Logik zum Hinzufügen von Dateien zum Archiv

	return nil
}
