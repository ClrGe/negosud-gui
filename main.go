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

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

// main sets up the window configuration and behaviour
func main() {
	a := app.NewWithID("negosud")
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("NEGOSUD")

	content := container.NewMax(loginForm(w))
	negosudLogo, _ := fyne.LoadResourceFromPath("media/logo.png")

	w.SetIcon(negosudLogo)
	w.Resize(fyne.NewSize(1920, 1080))
	w.SetContent(content)

	w.ShowAndRun()
}

// loginForm to perform an authentication to access the API
func loginForm(w fyne.Window) fyne.CanvasObject {
	negosudLogo := canvas.NewImageFromFile("media/logo.png")
	negosudLogo.FillMode = canvas.ImageFillContain
	negosudLogo.SetMinSize(fyne.NewSize(100, 100))

	text := canvas.NewText("Merci de vous identifier pour accéder à l'application", color.Black)
	text.TextSize = 15
	text.Alignment = fyne.TextAlignCenter

	emailLabel := canvas.NewText("Email", color.Black)
	emailInput := widget.NewEntry()
	emailInput.SetPlaceHolder("exemple@negosud.fr")
	emailInput.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "Saisir une adresse valide")

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
			a := fyne.CurrentApp()

			content := container.NewMax(homePage(w))
			title := widget.NewLabelWithStyle("Negosud", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

			currentPage = w
			changePage := func(t widgets.Component) {
				if fyne.CurrentDevice().IsMobile() {
					newPage := a.NewWindow(t.Title)

					currentPage = newPage
					newPage.SetContent(t.View(currentPage))
					newPage.Show()

					newPage.SetOnClosed(func() {
						currentPage = w
						return
					})
				}
				title.SetText(t.Title)
				content.Objects = []fyne.CanvasObject{t.View(w)}
				content.Refresh()
			}

			page := container.NewBorder(container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)

			// ! Without this : multiple windows open when clicking on a page !
			if fyne.CurrentDevice().IsMobile() {
				w.SetContent(makeNavigation(changePage, false))
			} else {
				split := container.NewHSplit(makeNavigation(changePage, true), page)
				split.Offset = 0.2
				w.SetContent(split)
			}
		},
		SubmitText: "Envoyer",
		CancelText: "",
	}

	formContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(800, 1080), form))
	box := container.NewVBox(negosudLogo, widget.NewSeparator(), text, widget.NewSeparator(), formContainer)
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(1920, 1080), box))
	return mainContainer
}

func isConnected(b bool) bool {
	if !b {
		b = false
		return false
	} else {
		b = true
		return true
	}
}

// homePage with logo and message
func homePage(w fyne.Window) fyne.CanvasObject {

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

// makeNavigation implements the navigation panel on the left of the screen
func makeNavigation(setTab func(component widgets.Component), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

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
			t, ok := widgets.Components[uid]
			if !ok {
				fyne.LogError("Missing something : "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)

			obj.(*widget.Label).TextStyle = fyne.TextStyle{}

		},
		OnSelected: func(uid string) {
			if t, ok := widgets.Components[uid]; ok {
				a.Preferences().SetString(currentTab, uid)
				setTab(t)
			}
		},
	}

	return container.NewBorder(nil, nil, nil, nil, arborescence)
}
