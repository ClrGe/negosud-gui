package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Order struct {
	ID       int    `json:"id"`
	Product  string `json:"bottle_id"`
	Quantity string `json:"quantity"`
	Seller   int    `json:"producer_id"`
}

func retrieveOrders(fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Historique des commandes", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		widget.NewLabel(""),
	))
}
