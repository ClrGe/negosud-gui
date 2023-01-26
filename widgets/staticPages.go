package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("media/logo-large.png")
	logo.FillMode = canvas.ImageFillContain

	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(900, 600))
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
		logo.SetMinSize(fyne.NewSize(900, 600))
	}
	return container.NewCenter(container.NewVBox(
		logo,
		widget.NewLabel(""),
	))
}