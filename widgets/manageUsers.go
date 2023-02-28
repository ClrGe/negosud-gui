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
	"strconv"
)

var BindUser []binding.DataMap

var UserTableRefreshMethod func()

// UsersColumns defines the header row for the table
var UsersColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 40},
	{ColName: "Name", Header: "Nom", WidthPercent: 120},
	{ColName: "Email", Header: "Email", WidthPercent: 120},
	{ColName: "Role", Header: "Rôle", WidthPercent: 120},
}

// makeUsersPage function creates a new set of tabs
func makeUsersPage(_ fyne.Window) fyne.CanvasObject {
	userListTab := container.NewTabItem("Liste des utilisateurs", displayUsers(nil))
	tabs := container.NewAppTabs(
		userListTab,
		container.NewTabItem("Ajouter un utilisateur", addUserForm(nil)),
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		if ti == userListTab {
			UserTableRefreshMethod()
		}
	}
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func getUsers() {
	//retrieve structs from data package
	BindUser = nil
	Users := data.Users
	source := "WIDGETS.USERS "
	response := data.AuthGetRequest("User")

	if response != nil {

		if err := json.NewDecoder(response).Decode(&Users); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}

		for i := 0; i < len(Users); i++ {
			Users[i].ID = strconv.Itoa(Users[i].Id)
			BindUser = append(BindUser, binding.BindStruct(&Users[i]))
		}
	}

}

func displayUsers(_ fyne.Window) fyne.CanvasObject {
	getUsers()
	userTableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: UsersColumns,
		Bindings: BindUser,
	}
	userTable := rtable.CreateTable(userTableOptions)
	UserTableRefreshMethod = func() {
		getUsers()
		userTableOptions.Bindings = BindUser
		userTable.Refresh()
	}
	return userTable
}

// addUserForm to add an authorized user
func addUserForm(_ fyne.Window) fyne.CanvasObject {
	source := "WIDGETS.USERS "
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
	roleUser := widget.NewSelectEntry([]string{"Administrateur", "Employé", "Client"})
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
			user := &data.User{
				FirstName: name.Text,
				Email:     emailInput.Text,
				Password:  passwordInput.Text,
				//Role:     roleUser.Text,
			}
			// convert struct to json
			jsonValue, err := json.Marshal(user)
			if err != nil {
				log(true, source, err.Error())
				fmt.Println(err)
				return
			}
			fmt.Print(bytes.NewBuffer(jsonValue))
			// Send data to API
			postData := data.AuthPostRequest("User/AddUser", bytes.NewBuffer(jsonValue))
			if postData != 201|200 {
				log(true, source, "Error on user creation"+string(jsonValue))
				fmt.Println("Error on user creation")
				return
			}
			fmt.Println("User created")
			// Refresh after add
			//UserRefreshMethod()
		},
		SubmitText: "Envoyer",
		CancelText: "",
	}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))
	return mainContainer
}
