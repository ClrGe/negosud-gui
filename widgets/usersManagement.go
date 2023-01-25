package widgets

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
	"image/color"
	"negosud-gui/config"
	"net/http"
	"strconv"
)

var BindUser []binding.DataMap

// makeUsersTabs function creates a new set of tabs
func makeUsersTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des utilisateurs", displayUsers(nil)),
		container.NewTabItem("Ajouter un utilisateur", addUserForm(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// UsersColumns defines the header row for the table
var UsersColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 10},
	{ColName: "Name", Header: "Nom", WidthPercent: 120},
	{ColName: "Email", Header: "Email", WidthPercent: 120},
	{ColName: "Role", Header: "Rôle", WidthPercent: 120},
}

func displayUsers(w fyne.Window) fyne.CanvasObject {
	// retrieve structs from config package
	Users := config.Users

	apiUrl := config.UserAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&Users); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(Users); i++ {
		t := Users[i]
		id := strconv.Itoa(t.Id)
		Users[i].ID = id
		BindUser = append(BindUser, binding.BindStruct(&Users[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: UsersColumns,
		Bindings: BindUser,
	}

	table := rtable.CreateTable(tableOptions)

	return table
}

// loginForm to perform an authentication to access the API
func loginForm(w fyne.Window) fyne.CanvasObject {

	var xPos, yPos, heightFields, heightLabels, widthForm float32

	xPos = 500
	yPos = 150
	heightFields = 50
	heightLabels = 20
	widthForm = 550

	text := canvas.NewText("Pour accéder à toutes les fonctionnalités, veuillez vous authentifier.", color.Black)
	text.TextSize = 20
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Resize(fyne.NewSize(widthForm, heightFields))
	text.Move(fyne.NewPos(430, 10))

	emailLabel := canvas.NewText("Email", color.Black)
	emailLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	emailLabel.Move(fyne.NewPos(xPos, yPos))
	email := widget.NewEntry()
	email.SetPlaceHolder("truc@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")
	email.Resize(fyne.NewSize(widthForm, heightFields))
	email.Move(fyne.NewPos(xPos, yPos+20))

	pwdLabel := canvas.NewText("Mot de passe", color.Black)
	pwdLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	pwdLabel.Move(fyne.NewPos(xPos, yPos+100))
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("****")
	password.Resize(fyne.NewSize(widthForm, heightFields))
	password.Move(fyne.NewPos(xPos, yPos+120))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(widthForm, heightFields))
	submitBtn.Move(fyne.NewPos(xPos, yPos+220))

	mainContainer := container.NewWithoutLayout(text, emailLabel, email, pwdLabel, password, submitBtn)
	return mainContainer
}

// addUserForm to add an authorized user
func addUserForm(w fyne.Window) fyne.CanvasObject {

	var xPos, yPos, heightFields, heightLabels, widthForm float32

	xPos = 500
	yPos = 150
	heightFields = 50
	heightLabels = 20
	widthForm = 550

	text := canvas.NewText("Pour ajouter un nouvel utilisateur, veuillez remplir ce formulaire.", color.Black)
	text.TextSize = 20
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Resize(fyne.NewSize(widthForm, heightFields))
	text.Move(fyne.NewPos(430, 10))

	nameLabel := canvas.NewText("Nom", color.Black)
	nameLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	nameLabel.Move(fyne.NewPos(xPos, yPos))
	name := widget.NewEntry()
	name.SetPlaceHolder("Jean Bon")
	name.Resize(fyne.NewSize(widthForm, heightFields))
	name.Move(fyne.NewPos(xPos, yPos+20))

	emailLabel := canvas.NewText("Email", color.Black)
	emailLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	emailLabel.Move(fyne.NewPos(xPos, yPos+100))
	email := widget.NewEntry()
	email.SetPlaceHolder("truc@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")
	email.Resize(fyne.NewSize(widthForm, heightFields))
	email.Move(fyne.NewPos(xPos, yPos+120))

	pwdLabel := canvas.NewText("Mot de passe", color.Black)
	pwdLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	pwdLabel.Move(fyne.NewPos(xPos, yPos+220))
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("****")
	password.Resize(fyne.NewSize(widthForm, heightFields))
	password.Move(fyne.NewPos(xPos, yPos+240))

	roleLabel := canvas.NewText("Rôle de l'utilisateur", color.Black)
	roleLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	roleLabel.Move(fyne.NewPos(xPos, yPos+340))
	roleUser := widget.NewSelectEntry([]string{"Administrateur", "Employé", "Intérimaire"})
	roleUser.SetPlaceHolder("Veuillez sélectionner un rôle...")
	roleUser.Resize(fyne.NewSize(widthForm, heightFields))
	roleUser.Move(fyne.NewPos(xPos, yPos+360))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(widthForm, heightFields))
	submitBtn.Move(fyne.NewPos(xPos, yPos+460))

	mainContainer := container.NewWithoutLayout(text, nameLabel, name, emailLabel, email, pwdLabel, password, roleLabel, roleUser, submitBtn)
	return mainContainer
}
