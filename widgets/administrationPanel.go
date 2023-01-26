package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// makeBottlesTabs creates a new set of tabs for bottles management
func makeDocumentsTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Devis", displayDocuments(nil)),
		container.NewTabItem("Factures", displayDocuments(nil)),
		container.NewTabItem("Autres", displayDocuments(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// makeBottlesTabs creates a new set of tabs for bottles management
func makeStatsTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Évolution des prix", displayStatistics(nil)),
		container.NewTabItem("Évolution des entrées/sorties", displayStatistics(nil)),
		container.NewTabItem("Popularité des produits", displayStatistics(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func displayDocuments(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewButtonWithIcon("Ajouter un document", theme.DocumentIcon(), func() { fmt.Print("Document envoyé") }),

		widget.NewLabelWithStyle("Liste des documents (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}

func displayStatistics(_ fyne.Window) fyne.CanvasObject {
	image := canvas.NewImageFromFile("media/stat2.png")
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(1280, 720))

	return container.NewCenter(container.NewVBox(
		image,
		widget.NewLabelWithStyle("Statistiques vente/entrées-sorties/prix etc (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
