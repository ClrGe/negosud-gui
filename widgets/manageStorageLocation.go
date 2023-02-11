package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
	"negosud-gui/data"
	"negosud-gui/widgets/controls"
	"strconv"
)

var BindStorageLocation []binding.DataMap
var StorageLocationTableRefreshMethod func()

// StorageLocationColumns defines the header row for the table
var StorageLocationColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 50},
	{ColName: "Name", Header: "Nom", WidthPercent: 90},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
	{ColName: "CreatedAt", Header: "Crée le", WidthPercent: 50},
}

// makeStorageLocationPage function creates a new set of tabs
func makeStorageLocationPage(_ fyne.Window) fyne.CanvasObject {
	storageLocationListTab := container.NewTabItem("Liste des emplacements", displayAndUpdateStorageLocations(nil))
	tabs := container.NewAppTabs(
		storageLocationListTab,
		container.NewTabItem("Ajouter un emplacement", addNewStorageLocation(nil)),
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		if ti == storageLocationListTab {
			StorageLocationTableRefreshMethod()
		}
	}
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func getStorageLocations() (bool, *widget.Label) {
	var source = "WIDGETS.STORAGELOCATION.getStorageLocations() "
	StorageLocationData := data.StorageLocationsData
	response := data.AuthGetRequest("storageLocation")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&StorageLocationData); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	BindStorageLocation = nil

	for i := 0; i < len(StorageLocationData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := StorageLocationData[i]
		id := strconv.Itoa(t.Id)
		StorageLocationData[i].ID = id

		// binding storageLocation data
		BindStorageLocation = append(BindStorageLocation, binding.BindStruct(&StorageLocationData[i]))
	}
	return true, widget.NewLabel("")
}

// displayAndUpdateStorageLocations implements a dynamic table bound to an editing form
func displayAndUpdateStorageLocations(_ fyne.Window) fyne.CanvasObject {
	var source = "WIDGETS.STORAGELOCATION.displayAndUpdateStorageLocations()"

	//region datas
	// retrieve structs from data package
	StorageLocation := data.IndStorageLocation

	var identifier string

	resp, label := getStorageLocations()
	if !resp {
		return label
	}

	bottle := getAllBottleName()
	bottleNames := make([]string, 0)
	for i := 0; i < len(bottle); i++ {
		bottle[i].ID = strconv.Itoa(bottle[i].Id)
		name := bottle[i].FullName
		bottleNames = append(bottleNames, name)
	}

	bottleMap := make(map[string]int)
	for i := 0; i < len(bottle); i++ {
		id := bottle[i].Id
		name := bottle[i].FullName
		bottleMap[name] = id
	}
	//endregion

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: StorageLocationColumns,
		Bindings: BindStorageLocation,
	}
	table := rtable.CreateTable(tableOptions)
	StorageLocationTableRefreshMethod = func() {
		getStorageLocations()
		tableOptions.Bindings = BindStorageLocation
		table.Refresh()
	}

	// Region UPDATE FORM
	updateForm := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelName := widget.NewLabel("Nom")
	entryName := widget.NewEntry()
	//entryName.Disable()

	//Header
	layoutControlItemName := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutControlItemName.Add(labelName)
	layoutControlItemName.Add(entryName)

	//BottleStorageLocation List

	// List Title
	BSLListTitle := widget.NewLabel("Bouteilles")
	BSLListTitle.TextStyle.Bold = true

	// List headers
	labelBottle := widget.NewLabel("Nom")
	labelQuantity := widget.NewLabel("Quantité")

	// List items

	//item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)
	//
	//item.SelectEntryBottle.OnChanged = func(name string) {
	//	item.BindBottleId(name)
	//}

	//selectEntryBottle := item.SelectEntryBottle
	//entryQuantity := item.EntryQuantity

	gridContainerHeader := &fyne.Container{Layout: layout.NewGridLayout(3)}
	// List headers
	gridContainerHeader.Add(widget.NewLabel(""))
	gridContainerHeader.Add(labelBottle)
	gridContainerHeader.Add(labelQuantity)
	// List items
	gridContainerItems := &fyne.Container{Layout: layout.NewGridLayout(3)}
	//gridContainerHeader.Add(selectEntryBottle)
	//gridContainerHeader.Add(entryQuantity)

	AddItemBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(),
		func() {
			addItem(bottleNames, bottleMap, gridContainerItems)
		})

	updateForm.Add(layoutControlItemName)
	updateForm.Add(widget.NewLabel(""))
	updateForm.Add(widget.NewSeparator())
	updateForm.Add(BSLListTitle)
	updateForm.Add(gridContainerHeader)
	updateForm.Add(gridContainerItems)
	updateForm.Add(AddItemBtn)

	editBtn := widget.NewButtonWithIcon("Modifier cet emplacement", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer cet emplacement", theme.WarningIcon(),
		func() {})

	//region " table events "
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
			//Sort method
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
		}

		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]
		_, err = rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(err.Error())
			log(true, source, err.Error())

			//return
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
			}
			// Fill form fields with fetched data
			entryName.SetText(StorageLocation.Name)

			gridContainerItems.RemoveAll()

			//
			for _, bsl := range StorageLocation.BottleStorageLocations {
				item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)
				item.Bind(bsl.Bottle.ID, bsl.StorageLocation.ID)
				item.SelectEntryBottle.Text = bsl.Bottle.FullName
				item.EntryQuantity.Text = strconv.Itoa(bsl.Quantity)

				var deleteItemBtn *widget.Button
				deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
					func() {
						gridContainerItems.Remove(deleteItemBtn)
						gridContainerItems.Remove(item.SelectEntryBottle)
						gridContainerItems.Remove(item.EntryQuantity)
					})

				gridContainerItems.Add(deleteItemBtn)
				gridContainerItems.Add(item.SelectEntryBottle)
				gridContainerItems.Add(item.EntryQuantity)
			}
			updateForm.Refresh()
		} else {
			log(true, source, err.Error())
		}
	}
	//endregion

	//region " updateForm events "
	editBtn.OnTapped = func() {
		storageLocation := &data.StorageLocation{
			ID:   StorageLocation.ID,
			Name: entryName.Text,
		}
		jsonValue, _ := json.Marshal(storageLocation)
		postData := data.AuthPostRequest("StorageLocation/UpdateStorageLocation/"+identifier, bytes.NewBuffer(jsonValue))
		if postData != 200 {
			fmt.Println("Error on update")
			message := "Error on storageLocation " + identifier + " update"
			log(true, source, message)
			return
		}
		fmt.Println("Success on update")
		StorageLocationTableRefreshMethod()
	}
	//updateForm.OnCancel = func() {
	//	fmt.Println("Canceled")
	//}
	//endregion

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion

	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithDelete := container.NewBorder(layoutUpdateForm, buttonsContainer, nil, nil)

	// Define layout
	individualTabs := container.NewAppTabs(
		container.NewTabItem("Modifier l'emplacement", layoutWithDelete),
	)

	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)
	mainContainer := container.New(layout.NewGridLayout(2))
	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new storageLocation to the API endpoint (POST)
func addNewStorageLocation(_ fyne.Window) fyne.CanvasObject {
	var source = "WIDGETS.STORAGELOCATION.addNewStorageLocation"
	nameLabel := widget.NewLabel("Nom")
	nameStorageLocation := widget.NewEntry()
	detailsLabel := widget.NewLabel("Description")
	detailsStorageLocation := widget.NewMultiLineEntry()
	pictureLabel := widget.NewLabel("Image")
	pictureStorageLocation := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	title := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: title},
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameStorageLocation},
			{Text: "", Widget: detailsLabel},
			{Text: "", Widget: detailsStorageLocation},
			{Text: "", Widget: pictureLabel},
			{Text: "", Widget: pictureStorageLocation},
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
			StorageLocationTableRefreshMethod()
		},
		SubmitText: "Envoyer",
	}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))

	return mainContainer
}

func getAllBottleName() []data.PartialBottle {
	var source = "WIDGETS.STORAGELOCATION.GetAllBottleName"
	bottleData := data.BottleData
	response := data.AuthGetRequest("Bottle")
	if response == nil {
		fmt.Println("No result returned")
		return nil
	}
	if err := json.NewDecoder(response).Decode(&bottleData); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return nil
	}
	return bottleData
}

func addItem(bottleNames []string, bottleMap map[string]int, gridContainerItems *fyne.Container) {
	item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)

	var deleteItemBtn *widget.Button
	deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
		func() {
			gridContainerItems.Remove(deleteItemBtn)
			gridContainerItems.Remove(item.SelectEntryBottle)
			gridContainerItems.Remove(item.EntryQuantity)
		})

	gridContainerItems.Add(deleteItemBtn)
	gridContainerItems.Add(item.SelectEntryBottle)
	gridContainerItems.Add(item.EntryQuantity)
}
