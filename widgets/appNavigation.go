package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// homePage with logo and message
func homePage(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("media/logo-large.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(900, 600))
	}
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Bienvenue dans l'utilitaire de gestion de NEGOSUD !", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		widget.NewLabel(""),
	))
}

type Component struct {
	Title      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
}

var (
	Components = map[string]Component{
		"home": {"Accueil",
			homePage,
			true,
		},
		"users_management": {"Gestion des utilisateurs",
			makeUsersTabs,
			true,
		},
		"storageLocations_management": {"Gestion des emplacements de stock",
			makeStorageLocationTabs,
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
		"":            {"home", "storageLocations_management", "producers_management", "bottles_management", "orders_tab", "admin_tab", "support_tab"},
		"support_tab": {"faq_tab", "contact_tab"},
		"orders_tab":  {"orders_producers", "orders_customers"},
		"admin_tab":   {"users_management", "documents_tab", "statistics_tab", "website_management"},
	}
)

const currentTab = "currentTab"

// Navigation implements the left-side navigation panel with layout defined in widgets/navigationLayout
func Navigation(setTab func(component Component), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	arborescence := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return ComponentIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := ComponentIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Nouvel onglet")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, _ := Components[uid]
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := Components[uid]; ok {
				a.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}

	// close a when hitting button
	disconnectUser := widget.NewButton("Déconnexion", func() {
		fmt.Println("user disconnected")
		fyne.CurrentApp().Quit()
	})

	return container.NewBorder(nil, disconnectUser, nil, nil, arborescence)
}
