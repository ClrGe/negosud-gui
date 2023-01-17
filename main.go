package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/components"
	"net/url"
)

const currentTab = "currentTab"

var topWindow fyne.Window

func main() {
	a := app.NewWithID("negosud")
	w := a.NewWindow("NEGOSUD")
	topWindow = w
	a.Settings().SetTheme(theme.LightTheme())

	w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("Onglet")

	setTab := func(t components.Component) {
		child := a.NewWindow(t.Title)
		topWindow = child
		child.SetContent(t.View(topWindow))
		child.Show()
		child.SetOnClosed(func() {
			topWindow = w
			return
		})

		title.SetText(t.Title)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tab := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	split := container.NewHSplit(makeNav(setTab, true), tab)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(1920, 1080))
	w.ShowAndRun()
}

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("Nouveau", nil)
	settingsItem := fyne.NewMenuItem("Paramètres", func() {
		// à faire : fonction ouverture settings
	})

	cutItem := fyne.NewMenuItem("Couper", func() {
	})

	copyItem := fyne.NewMenuItem("Copier", func() {
	})

	pasteItem := fyne.NewMenuItem("Coller", func() {
	})
	performFind := func() { fmt.Println("Chercher") }
	findItem := fyne.NewMenuItem("Chercher", performFind)

	helpMenu := fyne.NewMenu("Aide",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://negosud.com")
			_ = a.OpenURL(u)
		}))

	file := fyne.NewMenu("Fichier", newItem)
	file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)

	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Édition", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)

	return main
}

func makeNav(setTab func(component components.Component), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return components.ComponentIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := components.ComponentIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Nouvel onglet")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := components.Components[uid]
			if !ok {
				fyne.LogError("Missing something : "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)

			obj.(*widget.Label).TextStyle = fyne.TextStyle{}

		},
		OnSelected: func(uid string) {
			if t, ok := components.Components[uid]; ok {
				a.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}
