package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// makeWebsiteTabs creates a new set of tabs for bottle management
func makeWebsiteTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Statut du site", websiteStatus(nil)),
		container.NewTabItem("Statistiques", displayStats(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func websiteStatus(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Statut du site ecommerce (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}

func displayStats(_ fyne.Window) fyne.CanvasObject {
	image := canvas.NewImageFromFile("media/stat2.png")
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(1280, 720))

	return container.NewCenter(container.NewVBox(
		image,
		widget.NewLabelWithStyle("Statistiques du site ecommerce (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
