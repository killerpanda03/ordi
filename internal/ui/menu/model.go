package menu

type SelectMsg int

type Model struct {
	choices []string
	cursor  int      
}

func New() Model {
	return Model{
		choices: []string{"📂 Ein Verzeichnis organisieren", "Bilder sortieren", "Bild komprimieren", "Beenden"},
		cursor:  0,
	}
}
