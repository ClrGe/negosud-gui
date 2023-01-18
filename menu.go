package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"net/url"
)

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
