package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func loginForm(w fyne.Window) fyne.CanvasObject {

	email := widget.NewEntry()
	email.SetPlaceHolder("truc@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Mot de passe")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Email", Widget: email},
		},
		OnCancel: func() {
			fmt.Println("Annulation")
		},
		OnSubmit: func() {
			connectedScreen(w)

			loginSuccessDialog(w)
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title: "Form for: " + email.Text,
			})
		},
	}
	form.Append("Password", password)
	return form
}

// connectedScreen loads a tab panel for containers
func connectedScreen(w fyne.Window) fyne.CanvasObject {
	content := container.NewBorder(
		widget.NewLabelWithStyle("Top", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle("Bottom", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabel("Left"),
		widget.NewLabel("Right"),
		widget.NewLabel("Border Container"))
	return container.NewCenter(content)
}
