package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/data"
	"negosud-gui/widgets/Bottle"
	"negosud-gui/widgets/StorageLocation"
)

var log = data.Logger
var identifier string

type Component struct {
	Title      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
	Icon       fyne.Resource
}

var (
	Components = map[string]Component{
		"home": {"Accueil",
			homePage,
			true,
			theme.HomeIcon(),
		},
		"users_management": {"Gestion des utilisateurs",
			makeUsersPage,
			true,
			theme.AccountIcon(),
		},
		"storageLocations_management": {"Gestion des emplacements de stock",
			StorageLocation.MakePage,
			true,
			theme.StorageIcon(),
		},
		"producers_management": {"Gestion des producteurs",
			makeProducerPage,
			true,
			theme.FolderOpenIcon(),
		},
		"bottles_management": {"Gestion des produits",
			Bottle.MakePage,
			true,
			theme.FolderOpenIcon(),
		},
		"orders_tab": {"Commandes",
			makeOrdersTabs,
			true,
			theme.FolderIcon(),
		},
		"orders_producers": {"Commandes fournisseurs",
			makeOrdersTabs,
			true,
			theme.DownloadIcon(),
		},
		"orders_customers": {"Commandes clients",
			makeCusOrdersTabs,
			true,
			theme.DownloadIcon(),
		},
		//"website_management": {"Site e-commerce",
		//	makeWebsiteTabs,
		//	true},
		"support_tab": {"Support",
			displayFAQ,
			true,
			theme.InfoIcon(),
		},

		"faq_tab": {"FAQ",
			displayFAQ,
			true,
			theme.QuestionIcon(),
		},
		"contact_tab": {"Demander de l'aide",
			contactForm,
			true,
			theme.MailComposeIcon(),
		},
		"admin_tab": {"Administration",
			makeDocumentsTabs,
			true,
			theme.SettingsIcon(),
		},
		"documents_tab": {"Documents",
			makeDocumentsTabs,
			true,
			theme.DocumentIcon(),
			//"statistics_tab": {"Statistiques",
			//	makeStatsTabs,
			//	true,
			//},
			//"new_bottle_tab": {"NewBottle",
			//	makeBottleTabs,
			//	true,
			//},
		},
	}

	ComponentIndex = map[string][]string{
		"":            {"home", "storageLocations_management", "producers_management", "bottles_management", "orders_tab", "admin_tab", "support_tab", "new_bottle_tab"},
		"support_tab": {"faq_tab", "contact_tab"},
		"orders_tab":  {"orders_producers", "orders_customers"},
		//"admin_tab":   {"users_management", "documents_tab", "statistics_tab", "website_management"},

		"admin_tab": {"users_management", "documents_tab"},
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
			// for each branch, add icon to the left
			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			var icon fyne.CanvasObject
			if branch {
				icon = widget.NewIcon(theme.FolderIcon())
			} else {
				icon = widget.NewIcon(theme.DocumentIcon())
			}
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), icon, widget.NewLabel("Template Object"))
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			c := node.(*fyne.Container)
			// retrieve icon from map
			c.Objects[0].(*widget.Icon).SetResource(Components[uid].Icon)
			// set node height to 100

			l := c.Objects[1].(*widget.Label)
			l.SetText(Components[uid].Title)

		},
		OnSelected: func(uid string) {
			if t, ok := Components[uid]; ok {
				a.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}
	arborescence.ExtendBaseWidget(arborescence)
	// close a when hitting button
	disconnectUser := widget.NewButton("Déconnexion", func() {
		fmt.Println("user disconnected")
		fyne.CurrentApp().Quit()
	})

	return container.NewBorder(nil, disconnectUser, nil, nil, arborescence)
}

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
