package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
	"negosud-gui/data"
	"strconv"
)

var BindStorageLocation []binding.DataMap

// makeStorageLocationTabs function creates a new set of tabs
func makeStorageLocationTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des emplacements", displayAndUpdateStorageLocations(nil)),
		container.NewTabItem("Ajouter un emplacement", addNewStorageLocation(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// StorageLocationColumns defines the header row for the table
var StorageLocationColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 50},
	{ColName: "Name", Header: "Nom", WidthPercent: 90},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
	{ColName: "CreatedAt", Header: "Crée le", WidthPercent: 50},
}

// displayAndUpdateStorageLocations implements a dynamic table bound to an editing form
func displayAndUpdateStorageLocations(_ fyne.Window) fyne.CanvasObject {
	// retrieve structs from data package
	StorageLocation := data.IndStorageLocation
	StorageLocationData := data.StorageLocationData

	var identifier string
	var yPos, heightFields, widthForm float32
	yPos = 200
	heightFields = 50

	// DETAILS STORAGELOCATION
	// declare elements (empty or hidden until an identifier in the table gets clicked on)
	instructions := widget.NewLabelWithStyle("Cliquez sur un identifiant dans le tableau pour afficher les détails", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	instructions.Resize(fyne.NewSize(widthForm, heightFields))
	instructions.Move(fyne.NewPos(0, yPos-500))
	productImg := canvas.NewImageFromFile("media/wineyard.jpeg")
	productImg.FillMode = canvas.ImageFillContain
	productImg.SetMinSize(fyne.NewSize(600, 340))
	productImg.Hidden = true
	productTitle := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	productTitle.Resize(fyne.NewSize(widthForm, heightFields))
	productTitle.Move(fyne.NewPos(0, yPos-300))
	productDesc := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	productDesc.Resize(fyne.NewSize(widthForm, heightFields))
	productDesc.Move(fyne.NewPos(0, yPos-250))
	// UPDATE FORM
	// declare form elements
	nameLabel := widget.NewLabel("Nom")
	nameStorageLocation := widget.NewEntry()

	deleteBtn := widget.NewButtonWithIcon("Supprimer cet emplacement", theme.WarningIcon(),
		func() {})

	response := data.AuthGetRequest("storageLocation")
	if response == nil {
		fmt.Println("No result returned")
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&StorageLocationData); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return widget.NewLabel("Erreur de décodage du json")
	}

	for i := 0; i < len(StorageLocationData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := StorageLocationData[i]
		id := strconv.Itoa(t.Id)
		StorageLocationData[i].ID = id

		// binding storageLocation data
		BindStorageLocation = append(BindStorageLocation, binding.BindStruct(&StorageLocationData[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: StorageLocationColumns,
		Bindings: BindStorageLocation,
	}
	table := rtable.CreateTable(tableOptions)
	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 0 || cell.Row > len(BindStorageLocation) { // 1st col is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(StorageLocationColumns) {
			fmt.Println("*-> Column out of limits")
			return
		}
		// Handle header row clicked
		if cell.Row == 0 { // fmt.Println("-->", tblOpts.ColAttrs[cell.Col].Header)
			//// Add a row
			//BindStorageLocation = append(BindStorageLocation,
			//	binding.BindStruct(&data.StorageLocation{Name: "Belle Ambiance",
			//		Details: "brown", CreatedBy: "170"}))
			//tableOptions.Bindings = BindStorageLocation
			//table.Refresh()
			return
		}
		//Handle non-header row clicked
		var identifier string
		var err error
		if cell.Col == 0 {
			identifier, err = rtable.GetStrCellValue(cell, tableOptions)
			if err != nil {
				fmt.Println(err.Error())
				log(true, source, err.Error())
				return
			}
		} else {
			identifier, err = rtable.GetStrCellValue(widget.TableCellID{cell.Row, 0}, tableOptions)
			if err != nil {
				fmt.Println(err.Error())
				log(true, source, err.Error())
				return
			}
			//table.Select(widget.TableCellID{cell.Row, 0})
		}

		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]
		_, err = rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(err.Error())
			log(true, source, err.Error())

			return
		} else {
			instructions.Hidden = true
		}
		fmt.Println("-->", identifier)
		// Prevent app crash if other row than ID is clicked
		i, err := strconv.Atoi(identifier)
		if err == nil {
			fmt.Println(i)
			// Fetch individual storageLocation to fill form
			response := data.AuthGetRequest("storageLocation/" + identifier)
			if err := json.NewDecoder(response).Decode(&StorageLocation); err != nil {
				log(true, source, err.Error())
				fmt.Println(err)
			} else {
				productImg.Hidden = false
			}
			// Fill form fields with fetched data
			nameStorageLocation.SetText(StorageLocation.Name)
			productTitle.SetText(StorageLocation.Name)
		} else {
			log(true, source, err.Error())
		}
	}
	updateForm := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameStorageLocation},
		},
		OnSubmit: func() {
			storageLocation := &data.StorageLocation{
				ID:   StorageLocation.ID,
				Name: nameStorageLocation.Text,
			}
			jsonValue, _ := json.Marshal(storageLocation)
			postData := data.AuthPostRequest("StorageLocation/UpdateStorageLocation", bytes.NewBuffer(jsonValue))
			if postData != 201|200 {
				fmt.Println("Error on update")
				message := "Error on storageLocation " + identifier + " update"
				log(true, source, message)
			} else {
				fmt.Println("Success on update")
			}
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Envoyer",
		CancelText: "Annuler",
	}
	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewWithoutLayout(productTitle, productDesc))
	layoutDetailsTab := container.NewBorder(image, nil, nil, nil, textProduct, instructions)
	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithDelete := container.NewBorder(layoutUpdateForm, deleteBtn, nil, nil)

	// Define layout
	individualTabs := container.NewAppTabs(
		container.NewTabItem("Détails de l'emplacement", layoutDetailsTab),
		container.NewTabItem("Modifier l'emplacement", layoutWithDelete),
	)
	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)
	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new storageLocation to the API endpoint (POST)
func addNewStorageLocation(_ fyne.Window) fyne.CanvasObject {
	nameLabel := widget.NewLabel("Nom")
	nameStorageLocation := widget.NewEntry()

	title := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: title},
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameStorageLocation},
		},
		OnSubmit: func() {
			storageLocation := &data.StorageLocation{
				Name: nameStorageLocation.Text,
			}
			// convert storageLocation struct to json
			jsonValue, err := json.Marshal(&storageLocation)
			if err != nil {
				fmt.Println(err)
				log(true, source, err.Error())
				return
			}
			postData := data.AuthPostRequest("StorageLocation/AddStorageLocation", bytes.NewBuffer(jsonValue))
			if postData != 200|201 {
				message := "StatusCode " + strconv.Itoa(postData)
				log(true, source, message)
				fmt.Println(message)
				return
			}
			fmt.Println("New storageLocation added with success")
		},
		SubmitText: "Envoyer",
	}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))

	return mainContainer
}
