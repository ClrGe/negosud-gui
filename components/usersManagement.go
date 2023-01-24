package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func displayUsers(fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Liste des utilisateurs", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}

func loginForm(w fyne.Window) fyne.CanvasObject {

	email := widget.NewEntry()
	email.SetPlaceHolder("truc@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")
	email.Resize(fyne.NewSize(300, 35))
	email.Move(fyne.NewPos(100, 150))

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Mot de passe")
	password.Resize(fyne.NewSize(300, 35))
	password.Move(fyne.NewPos(100, 200))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(300, 50))
	submitBtn.Move(fyne.NewPos(100, 300))

	mainContainer := container.NewWithoutLayout(email, password, submitBtn)
	return mainContainer
}
