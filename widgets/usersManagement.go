package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
	"negosud-gui/data"
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

// makeUsersTabs function creates a new set of tabs
func makeHomeTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Accueil", welcomeScreen(nil)),
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
	// retrieve structs from data package
	Users := data.Users

	apiUrl := data.UserAPIConfig()

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

// addUserForm to add an authorized user
func addUserForm(w fyne.Window) fyne.CanvasObject {

	nameLabel := widget.NewLabelWithStyle("Nom", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	name := widget.NewEntry()
	name.SetPlaceHolder("Jean Bon")

	emailLabel := widget.NewLabelWithStyle("Email", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	emailInput := widget.NewEntry()
	emailInput.SetPlaceHolder("truc@example.com")
	emailInput.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	passwordLabel := widget.NewLabelWithStyle("Mot de passe", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	passwordInput := widget.NewPasswordEntry()
	passwordInput.SetPlaceHolder("******")

	roleLabel := widget.NewLabelWithStyle("Rôle", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	roleUser := widget.NewSelectEntry([]string{"Administrateur", "Employé", "Intérimaire"})
	roleUser.SetPlaceHolder("Veuillez sélectionner un rôle...")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: name},
			{Text: "", Widget: emailLabel},
			{Text: "", Widget: emailInput},
			{Text: "", Widget: passwordLabel},
			{Text: "", Widget: passwordInput},
			{Text: "", Widget: roleLabel},
			{Text: "", Widget: roleUser},
		},
		OnSubmit: func() {
			email := emailInput.Text
			password := passwordInput.Text

			apiUrl := data.LoginAPIConfig(email, password)

			resp, err := http.Post(apiUrl, "application/json", bytes.NewBufferString(email))

			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error on login form submitting: " + err.Error(),
				})
				return
			}
			if resp.StatusCode == 204 {
				fmt.Println("Could not log in")
				return
			}
			fmt.Println("User is connected")
		},
		SubmitText: "Envoyer",
		CancelText: "",
	}

	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))

	return mainContainer
}
