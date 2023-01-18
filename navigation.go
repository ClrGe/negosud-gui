package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/components"
)

func makeNav(setTab func(component components.Component), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	arborescence := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return components.ComponentIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := components.ComponentIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Nouvel onglet")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := components.Components[uid]
			if !ok {
				fyne.LogError("Missing something : "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)

			obj.(*widget.Label).TextStyle = fyne.TextStyle{}

		},
		OnSelected: func(uid string) {
			if t, ok := components.Components[uid]; ok {
				a.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}

	return container.NewBorder(nil, nil, nil, nil, arborescence)
}
