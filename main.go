package main

import (
	"os"
	"example/ordi/internal/organizer"
)

func main() {
	if len(os.Args) < 2 {
		println("Bitte geben Sie den zu organisierenden Verzeichnis-Pfad an.")
	}

	dirPath := os.Args[1]
	err := organizer.Organize(dirPath)
	if err != nil {
		println("Fehler beim Organisieren des Verzeichnisses:", err.Error())

	} else {
		println("Verzeichnis erfolgreich organisiert.")
	}

}
