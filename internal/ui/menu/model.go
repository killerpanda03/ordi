package menu

type SelectMsg int

type Model struct {
	choices []string
	cursor  int
}

func New() Model {
	return Model{
		choices: []string{"Ein Verzeichnis organisieren", "Duplikate finden", "Dateien komprimieren",  "Beenden"},
		cursor:  0,
	}
}
