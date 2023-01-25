package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/widgets"
)

// --------------------------------------------------------------------

const currentTab = "currentTab"

var activePage fyne.Window

// define and start the window
func main() {
	a := app.NewWithID("negosud")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("NEGOSUD")

	activePage = w

	w.SetMainMenu(appBarMenu(a, w))
	w.SetMaster()

	content := container.NewMax(homePage(w))

	title := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	setTab := func(t widgets.Component) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			activePage = child
			child.SetContent(t.View(activePage))
			child.Show()
			child.SetOnClosed(func() {
				activePage = w
				return
			})
		}
		title.SetText(t.Title)
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tab := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)

	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNavigation(setTab, false))
	} else {
		split := container.NewHSplit(makeNavigation(setTab, true), tab)
		split.Offset = 0.2
		w.SetContent(split)
	}

	negosudLogo, _ := fyne.LoadResourceFromPath("media/logo.png")

	w.SetIcon(negosudLogo)
	w.Resize(fyne.NewSize(1920, 1080))
	w.ShowAndRun()
}

// static homepage with logo and welcome message
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

func makeNavigation(setTab func(component widgets.Component), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	arborescence := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return widgets.ComponentIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := widgets.ComponentIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Nouvel onglet")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := widgets.Components[uid]
			if !ok {
				fyne.LogError("Missing something : "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)

			obj.(*widget.Label).TextStyle = fyne.TextStyle{}

		},
		OnSelected: func(uid string) {
			if t, ok := widgets.Components[uid]; ok {
				a.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}

	return container.NewBorder(nil, nil, nil, nil, arborescence)
}

// TODO : implement functions for menu items
func appBarMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	param := func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	}
	parametersItem := fyne.NewMenuItem("Paramètres", param)

	performFind := func() { fmt.Println("Recherche") }
	findItem := fyne.NewMenuItem("Recherche", performFind)
	negosudMenu := fyne.NewMenu("Negosud", parametersItem, findItem)
	negosudMenu.Items = append(negosudMenu.Items, fyne.NewMenuItemSeparator())

	helpMenu := fyne.NewMenu("Support",
		fyne.NewMenuItem("FAQ", func() {}),
		fyne.NewMenuItem("Contacter les développeurs", func() {}),
		fyne.NewMenuItem("Documentation", func() {}),
	)
	main := fyne.NewMainMenu(
		negosudMenu,
		helpMenu,
	)

	return main
}
