package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeAppTabsTab(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Éditer le producteur", container.New(nil)),
		container.NewTabItem("Ajouter un producteur", widget.NewLabel("Content of tab 2")),
		container.NewTabItem("Détails", widget.NewLabel("Content of tab 3")),
	)
	//for i := 4; i <= 12; i++ {
	//	tabs.Append(container.NewTabItem(fmt.Sprintf("Tab %d", i), widget.NewLabel(fmt.Sprintf("Content of tab %d", i))))
	//}
	return container.NewBorder(nil, nil, nil, nil, tabs)
}
