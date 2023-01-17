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
		"welcome": {"Accueil", "", welcomeScreen, true},
		"box": {"Clients",
			"Liste des clients",
			makeClientsTab,
			true,
		},
		"table": {"Stock",
			"Liste des bouteilles",
			makeTableTab,
			true,
		},
	}

	ComponentIndex = map[string][]string{
		"": {"welcome", "box", "table"},
	}
)
