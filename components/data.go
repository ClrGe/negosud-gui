package components

import (
	"fyne.io/fyne/v2"
)

type Component struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	Components = map[string]Component{
		"welcome": {"Accueil",
			"Bienvenue",
			welcomeScreen,
			true},
		"box": {"Producteurs",
			"Liste des producteurs",
			retrieveProducers,
			true,
		},
		"table": {"Bouteilles",
			"Liste des bouteilles",
			retrieveBottles,
			true,
		},
	}

	ComponentIndex = map[string][]string{
		"": {"welcome", "box", "table"},
	}
)
