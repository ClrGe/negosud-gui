package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TODO : display the list of all custormers
func displayCustomers(fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Liste des clients", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
