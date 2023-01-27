package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"negosud-gui/data"
)

// LoginForm to authenticate and receive a token
func LoginForm(w fyne.Window) fyne.CanvasObject {
	appLogo := canvas.NewImageFromFile("media/logo.png")
	appLogo.FillMode = canvas.ImageFillContain
	appLogo.SetMinSize(fyne.NewSize(100, 100))

	text := widget.NewLabelWithStyle("Merci de vous identifier pour accéder à l'application", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	emailLabel := canvas.NewText("Email", color.Black)
	emailInput := widget.NewEntry()
	emailInput.SetPlaceHolder("exemple@negosud.fr")
	emailInput.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "Adresse e-mail invalide !")

	passwordLabel := canvas.NewText("Mot de passe", color.Black)
	passwordInput := widget.NewPasswordEntry()
	passwordInput.SetPlaceHolder("******")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: emailLabel},
			{Text: "", Widget: emailInput},
			{Text: "", Widget: passwordLabel},
			{Text: "", Widget: passwordInput},
		},
		OnSubmit: func() {
			data.LoginAndSaveToken(emailInput.Text, passwordInput.Text)

			if data.LoginAndSaveToken(emailInput.Text, passwordInput.Text) != 200 {
				data.Logger(false, "LOGIN ", "Failed : Incorrect email or password")
				text.SetText("Identifiants incorrects !")
			} else {
				a := fyne.CurrentApp()
				content := container.NewMax(homePage(w))
				window := w
				changePage := func(c Component) {
					if fyne.CurrentDevice().IsMobile() {
						newPage := a.NewWindow(c.Title)
						window = newPage
						newPage.SetContent(c.View(window))
						newPage.Show()
						newPage.SetOnClosed(func() {
							window = w
							return
						})
					}
					content.Objects = []fyne.CanvasObject{c.View(w)}
					content.Refresh()
				}
				page := container.NewBorder(container.NewVBox(widget.NewSeparator()), nil, nil, nil, content)
				// responsive
				if fyne.CurrentDevice().IsMobile() {
					w.SetContent(Navigation(changePage, false))
				} else {
					split := container.NewHSplit(Navigation(changePage, true), page)
					split.Offset = 0.2
					w.SetContent(split)
				}
			}
		},
		OnCancel: func() {
			// close app when clicking on cancel button
			fyne.CurrentApp().Quit()
		},
		SubmitText: "Envoyer",
		CancelText: "Quitter",
	}
	form.Resize(fyne.NewSize(800, 400))
	form.Move(fyne.NewPos(555, 50))

	// LAYOUT
	spacer := widget.NewLabel("")
	formContainer := container.NewWithoutLayout(form)
	layoutPage := container.NewVBox(spacer, spacer, appLogo, widget.NewSeparator(), text, widget.NewSeparator(), formContainer)
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(1900, 800), layoutPage))
	return mainContainer
}
