package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// grouping all information dialogs needed in the package

func bottleSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Nouvelle bouteille ajoutée", w)
}

func bottleFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Impossible d'ajouter la nouvelle bouteille", w)
}

func producerSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Nouveau producteur ajouté", w)
}

func producerFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Impossible d'ajouter le nouveau producteur", w)
}
