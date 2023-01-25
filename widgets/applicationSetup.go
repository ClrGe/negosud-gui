package widgets

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
		"homepage": {"Accueil",
			"Bienvenue",
			welcomeScreen,
			true},
		"login": {"Connexion",
			"Connexion",
			loginForm,
			true,
		},
		"users_management": {"Gestion des utilisateurs",
			"utilisateurs",
			makeUsersTabs,
			true,
		},
		"producers_management": {"Gestion des producteurs",
			"producteurs",
			makeProducerTabs,
			true,
		},
		"bottles_management": {"Gestion des produits",
			"produits",
			makeBottleTabs,
			true,
		},
		"orders_producers": {"Gestion des commandes fournisseurs",
			"commandes",
			makeOrdersTabs,
			true,
		},
		"orders_customers": {"Gestion des commandes clients",
			"commandes",
			makeCusOrdersTabs,
			true,
		},
		"support_tab": {"Support",
			"aide",
			logoScreen,
			true,
		},
		"faq_tab": {"FAQ",
			"aide",
			displayFAQ,
			true},
		"contact_tab": {"Demander de l'aide",
			"aide",
			contactForm,
			true},
	}

	ComponentIndex = map[string][]string{
		"":            {"homepage", "login", "users_management", "producers_management", "bottles_management", "orders_producers", "orders_customers", "support_tab"},
		"support_tab": {"faq_tab", "contact_tab"},
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
