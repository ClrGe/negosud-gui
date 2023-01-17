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
		"producers": {"Producteurs",
			"Liste des producteurs",
			retrieveProducers,
			true,
		},
		"bottles": {"Bouteilles",
			"Liste des bouteilles",
			retrieveBottles,
			true,
		},
		"prod_orders": {"Commandes fournisseurs",
			"Historique des commandes fournisseurs",
			retrieveOrders,
			true,
		},
		"cl_orders": {"Commandes clients",
			"Historique des commandes clients",
			retrieveOrders,
			true,
		},
	}

	ComponentIndex = map[string][]string{
		"": {"welcome", "producers", "bottles", "prod_orders", "cl_orders"},
	}
)
