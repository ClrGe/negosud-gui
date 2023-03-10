package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"negosud-gui/widgets"
)

var window fyne.Window

// main sets up the window configuration and behaviour
func main() {
	a := app.NewWithID("negosud")
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("NEGOSUD")

	content := container.NewMax(widgets.LoginForm(w))
	negosudLogo, _ := fyne.LoadResourceFromPath("media/logo.png")

	w.SetIcon(negosudLogo)
	w.SetContent(content)
	w.CenterOnScreen()
	w.SetMaster()

	w.CenterOnScreen()

	w.ShowAndRun()
}
