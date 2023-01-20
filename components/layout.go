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
		"login": {"Connexion",
			"Connexion",
			loginForm,
			true,
		},
		"connected": {"Partie authentifiée",
			"Partie authentifiée",
			loginForm,
			true,
		},
		"producers": {"Producteurs",
			"Producteurs",
			displayProducers,
			true,
		},
		"bottles": {"Produits",
			"Produits",
			displayBottles,
			true,
		},
		"prod_orders": {"Commandes producteurs",
			"Historique des commandes producteurs",
			retrieveOrders,
			true,
		},
		"cl_orders": {"Commandes clients",
			"Historique des commandes clients",
			retrieveOrders,
			true,
		},
		"users": {"Liste des utilisateurs",
			"Liste des utilisateurs",
			displayUsers,
			true,
		},
		"producer_management": {"Gestion producteurs",
			"Gestion producteurs",
			logoScreen,
			true,
		},
		"bottle_management": {"Gestion produits",
			"Gestion produits",
			logoScreen,
			true,
		},
		"customer_management": {"Gestion clients",
			"Gestion clients",
			logoScreen,
			true,
		},
		"addProd": {"Ajouter producteur",
			"Ajouter un producteur",
			producerForm,
			true,
		},
		"addBottle": {"Ajouter un produit",
			"Ajouter un produit",
			bottleForm,
			true,
		},
		"tabs": {"Test onglets",
			"test",
			makeAppTabsTab,
			true},
	}

	ComponentIndex = map[string][]string{
		"":                    {"welcome", "login", "connected", "tabs"},
		"connected":           {"producer_management", "bottle_management", "customer_management"},
		"producer_management": {"producers", "prod_orders", "addProd"},
		"bottle_management":   {"bottles", "addBottle"},
		"customer_management": {"cl_orders", "users"},
	}
)
