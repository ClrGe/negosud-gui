package widgets

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
		"login": {"Accueil",
			"Connexion",
			makeHomeTabs,
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
			makeBottlesTabs,
			true,
		},
		"orders_tab": {"Commandes",
			"commandes",
			logoScreen,
			true},
		"orders_producers": {"Commandes fournisseurs",
			"commandes",
			makeOrdersTabs,
			true,
		},
		"orders_customers": {"Commandes clients",
			"commandes",
			makeCusOrdersTabs,
			true,
		},
		"website_management": {"Site e-commerce",
			"Gestion site web",
			makeWebsiteTabs,
			true},
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
		"":            {"login", "users_management", "producers_management", "bottles_management", "orders_tab", "support_tab", "website_management"},
		"support_tab": {"faq_tab", "contact_tab"},
		"orders_tab":  {"orders_producers", "orders_customers"},
	}
)
