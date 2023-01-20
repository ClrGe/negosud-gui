package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
			producers.displayProducers,
			true,
		},
		"bottles": {"Produits",
			"Produits",
			displayBottles,
			true,
		},
		"prod_orders": {"Commandes producteurs",
			"Historique des commandes producteurs",
			displayOrders,
			true,
		},
		"cl_orders": {"Commandes clients",
			"Historique des commandes clients",
			displayOrders,
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
			producers.producerForm,
			true,
		},
		"addBottle": {"Ajouter un produit",
			"Ajouter un produit",
			bottleForm,
			true,
		},
		"producerTab": {"Gestion des producteurs",
			"producteurs",
			makeProducerTabs,
			true,
		},
		"bottleTab": {"Gestion des produits",
			"produits",
			makeBottleTabs,
			true},
	}

	ComponentIndex = map[string][]string{
		"":                    {"welcome", "login", "connected", "producerTab", "bottleTab"},
		"connected":           {"producer_management", "bottle_management", "customer_management"},
		"producer_management": {"producers", "prod_orders", "addProd"},
		"bottle_management":   {"bottles", "addBottle"},
		"customer_management": {"cl_orders", "users"},
	}
)

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("media/logo-large.png")
	logo.FillMode = canvas.ImageFillContain

	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(1364, 920))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Bienvenue sur l'utilitaire de gestion de stock de NEGOSUD !", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		widget.NewLabel(""),
	))
}

func logoScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("media/logo-large.png")
	logo.FillMode = canvas.ImageFillContain

	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(1364, 920))
	}

	return container.NewCenter(container.NewVBox(
		logo,
		widget.NewLabel(""),
	))
}
