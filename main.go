package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"negosud-gui/widgets"
)

const currentTab = "currentTab"

var currentPage fyne.Window

// main sets up the window configuration and behaviour
func main() {
	app := app.NewWithID("negosud")
	app.Settings().SetTheme(theme.LightTheme())

	w := app.NewWindow("NEGOSUD")

	content := container.NewMax(loginForm(w))
	negosudLogo, _ := fyne.LoadResourceFromPath("media/logo.png")

	w.SetIcon(negosudLogo)
	w.Resize(fyne.NewSize(1920, 1080))
	w.SetContent(content)

	w.ShowAndRun()
}

// loginForm to perform an authentication to access the API
func loginForm(w fyne.Window) fyne.CanvasObject {
	appLogo := canvas.NewImageFromFile("media/logo.png")
	appLogo.FillMode = canvas.ImageFillContain
	appLogo.SetMinSize(fyne.NewSize(100, 100))

	text := canvas.NewText("Merci de vous identifier pour accéder à l'application", color.Black)
	text.TextSize = 15
	text.Alignment = fyne.TextAlignCenter

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
			app := fyne.CurrentApp()
			content := container.NewMax(homePage(w))
			currentPage = w
			changePage := func(c widgets.Component) {
				if fyne.CurrentDevice().IsMobile() {
					newPage := app.NewWindow(c.Title)
					currentPage = newPage
					newPage.SetContent(c.View(currentPage))
					newPage.Show()
					newPage.SetOnClosed(func() {
						currentPage = w
						return
					})
				}
				content.Objects = []fyne.CanvasObject{c.View(w)}
				content.Refresh()
			}
			page := container.NewBorder(container.NewVBox(widget.NewSeparator()), nil, nil, nil, content)
			// responsive
			if fyne.CurrentDevice().IsMobile() {
				w.SetContent(makeNavigation(changePage, false))
			} else {
				split := container.NewHSplit(makeNavigation(changePage, true), page)
				split.Offset = 0.2
				w.SetContent(split)
			}
		},
		OnCancel: func() {
			// close app when clicking on cancel button
			fyne.CurrentApp().Quit()
		},
		SubmitText: "Envoyer",
		CancelText: "Quitter",
	}
	form.Resize(fyne.NewSize(800, 200))
	form.Move(fyne.NewPos(555, 100))

	// LAYOUT
	spacer := widget.NewLabel("")
	formContainer := container.NewWithoutLayout(form)
	layoutPage := container.NewVBox(spacer, spacer, appLogo, widget.NewSeparator(), text, widget.NewSeparator(), formContainer)
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(1900, 800), layoutPage))
	return mainContainer
}

// homePage with logo and message
func homePage(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("media/logo-large.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(900, 600))
	}
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Bienvenue dans l'utilitaire de gestion de NEGOSUD !", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		widget.NewLabel(""),
	))
}

// makeNavigation implements the left-side navigation panel with layout defined in widgets/navigationLayout
func makeNavigation(setTab func(component widgets.Component), loadPrevious bool) fyne.CanvasObject {
	app := fyne.CurrentApp()
	arborescence := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return widgets.ComponentIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := widgets.ComponentIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Nouvel onglet")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, _ := widgets.Components[uid]
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := widgets.Components[uid]; ok {
				app.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}

	// close app when hitting button
	disconnectUser := widget.NewButton("Déconnexion", func() {
		fmt.Println("user disconnected")
		fyne.CurrentApp().Quit()
	})

	return container.NewBorder(nil, disconnectUser, nil, nil, arborescence)
}
