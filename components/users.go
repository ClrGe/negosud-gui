package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func displayUsers(fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Liste des utilisateurs", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		widget.NewLabel(""),
	))
}
