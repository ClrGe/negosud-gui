package widgets

import (
	"fyne.io/fyne/v2"
)

type Component struct {
	Title      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
}

var (
	Components = map[string]Component{
		"home": {"Accueil",
			welcomeScreen,
			true,
		},
		"users_management": {"Gestion des utilisateurs",
			makeUsersTabs,
			true,
		},
		"producers_management": {"Gestion des producteurs",
			makeProducerTabs,
			true,
		},
		"bottles_management": {"Gestion des produits",
			makeBottlesTabs,
			true,
		},
		"orders_tab": {"Commandes",
			makeOrdersTabs,
			true},
		"orders_producers": {"Commandes fournisseurs",
			makeOrdersTabs,
			true,
		},
		"orders_customers": {"Commandes clients",
			makeCusOrdersTabs,
			true,
		},
		"website_management": {"Site e-commerce",
			makeWebsiteTabs,
			true},
		"support_tab": {"Support",
			displayFAQ,
			true,
		},
		"faq_tab": {"FAQ",
			displayFAQ,
			true},
		"contact_tab": {"Demander de l'aide",
			contactForm,
			true,
		},
		"admin_tab": {"Administration",
			makeDocumentsTabs,
			true,
		},
		"documents_tab": {"Documents",
			makeDocumentsTabs,
			true},
		"statistics_tab": {"Statistiques",
			makeStatsTabs,
			true,
		},
	}

	ComponentIndex = map[string][]string{
		"":            {"home", "users_management", "producers_management", "bottles_management", "orders_tab", "admin_tab", "support_tab"},
		"support_tab": {"faq_tab", "contact_tab"},
		"orders_tab":  {"orders_producers", "orders_customers"},
		"admin_tab":   {"documents_tab", "statistics_tab", "website_management"},
	}
)
