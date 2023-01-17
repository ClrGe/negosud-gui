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
			displayProducers,
			true,
		},
		"bottles": {"Bouteilles",
			"Liste des bouteilles",
			displayBottles,
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
		"users": {"Utilisateurs",
			"Liste des utilisateurs",
			displayUsers,
			true,
		},
		"addProd": {"Ajouter producteur",
			"Ajouter un producteur",
			producerForm,
			true,
		},
		"addBottle": {"Ajouter produit",
			"Ajouter un produit",
			bottleForm,
			true,
		},
	}

	ComponentIndex = map[string][]string{
		"": {"welcome", "producers", "bottles", "prod_orders", "cl_orders", "users", "addProd", "addBottle"},
	}
)
