package main

import (
	"fmt"
	"os"

	"example/ordi/internal/ui/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	mainModel := app.New()

	p := tea.NewProgram(mainModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Oje, es gab einen Fehler: %v\n", err)
		os.Exit(1)
	}
}
