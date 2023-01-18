package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/components"
)

const currentTab = "currentTab"

var activePage fyne.Window

func main() {
	a := app.NewWithID("negosud")
	w := a.NewWindow("NEGOSUD - Utilitaire de gestion")

	activePage = w
	a.Settings().SetTheme(theme.LightTheme())
	w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	content := container.NewMax(welcomeScreen(w))

	title := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	setTab := func(t components.Component) {
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
		w.SetContent(makeNav(setTab, false))
	} else {
		split := container.NewHSplit(makeNav(setTab, true), tab)
		split.Offset = 0.2
		w.SetContent(split)
	}
	negosudLogo, _ := fyne.LoadResourceFromPath("media/logo.png")
	w.SetIcon(negosudLogo)

	w.Resize(fyne.NewSize(1920, 1080))
	w.ShowAndRun()
}
