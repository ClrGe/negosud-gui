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
	"sort"
	"strconv"
)

var BindStorageLocation []binding.DataMap
var StorageLocationTableRefreshMethod func()
var BottleStorageLocationControls map[*controls.BottleStorageLocationItem]int
var StorageLocationUpdateFormClearMethod func()
var StorageLocationAddFormClearMethod func()

// makeStorageLocationPage function creates a new set of tabs
func makeStorageLocationPage(_ fyne.Window) fyne.CanvasObject {
	storageLocationListTab := container.NewTabItem("Liste des emplacements", displayAndUpdateStorageLocations(nil))
	addStorageLocationTab := container.NewTabItem("Ajouter un emplacement", addNewStorageLocation(nil))
	tabs := container.NewAppTabs(
		storageLocationListTab,
		addStorageLocationTab,
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		if ti == storageLocationListTab {
			StorageLocationTableRefreshMethod()
			StorageLocationUpdateFormClearMethod()
		}
		if ti == addStorageLocationTab {
			StorageLocationAddFormClearMethod()
			//BottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
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

	//order datas by Id
	sort.SliceStable(StorageLocationData, func(i, j int) bool {
		return StorageLocationData[i].Id < StorageLocationData[j].Id
	})

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
	//var source = "WIDGETS.STORAGELOCATION.displayAndUpdateStorageLocations()"

	//region datas
	// retrieve structs from data package
	StorageLocation := data.IndStorageLocation

	resp, label := getStorageLocations()
	if !resp {
		return label
	}

	bottleNames, bottleMap := StorageLocation_GetAndMapBottleNames()

	// StorageLocationColumns defines the header row for the table
	var StorageLocationColumns = []rtable.ColAttr{
		{ColName: "ID", Header: "ID", WidthPercent: 50},
		{ColName: "Name", Header: "Nom", WidthPercent: 90},
		{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
		{ColName: "CreatedAt", Header: "Crée le", WidthPercent: 50},
	}

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

	//region UPDATE FORM

	updateForm, entryName, gridContainerItems := StorageLocation_InitForm(bottleNames, bottleMap)

	StorageLocationUpdateFormClearMethod = func() {
		table.UnselectAll()
		entryName.Text = ""
		entryName.Refresh()
		gridContainerItems.RemoveAll()
		BottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		StorageLocation_TableOnSelected(cell, StorageLocationColumns, tableOptions, &StorageLocation, entryName, gridContainerItems, bottleNames, bottleMap, updateForm)
	}
	//endregion

	//region " design elements initialization "
	buttonsContainer := StorageLocation_InitButtonContainer(&StorageLocation, entryName)
	mainContainer := StorageLocation_InitMainContainer(updateForm, table, buttonsContainer)
	//endregion

	return mainContainer
}

func StorageLocation_TableOnSelected(cell widget.TableCellID, StorageLocationColumns []rtable.ColAttr, tableOptions *rtable.TableOptions, StorageLocation *data.StorageLocation, entryName *widget.Entry, gridContainerItems *fyne.Container, bottleNames []string, bottleMap map[string]int, updateForm *fyne.Container) {
	if cell.Row < 0 || cell.Row > len(BindStorageLocation) { // 1st col is header
		fmt.Println("*-> Row out of limits")
		log(true, source, "*-> Row out of limits")
		return
	}
	if cell.Col < 0 || cell.Col >= len(StorageLocationColumns) {
		fmt.Println("*-> Column out of limits")
		log(true, source, "*-> Column out of limits")
		return
	}
	// Handle header row clicked
	if cell.Row == 0 { // fmt.Println("-->", tblOpts.ColAttrs[cell.Col].Header)
		//Sort method
		return
	}
	//Handle non-header row clicked

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
		response := data.AuthGetRequest("StorageLocation/" + identifier)
		if err := json.NewDecoder(response).Decode(&StorageLocation); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		entryName.SetText(StorageLocation.Name)

		gridContainerItems.RemoveAll()

		BottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)

		//
		for _, bsl := range StorageLocation.BottleStorageLocations {
			item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)
			item.Bind(bsl.Bottle.ID, bsl.StorageLocation.ID)
			item.SelectEntryBottle.Text = bsl.Bottle.FullName
			item.EntryQuantity.Text = strconv.Itoa(bsl.Quantity)

			BottleStorageLocationControls[item] = len(BottleStorageLocationControls) + 1

			var deleteItemBtn *widget.Button
			deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
				func() {
					gridContainerItems.Remove(deleteItemBtn)
					gridContainerItems.Remove(item.SelectEntryBottle)
					gridContainerItems.Remove(item.EntryQuantity)
					delete(BottleStorageLocationControls, item)
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

// region " design initializers"
func StorageLocation_InitForm(bottleNames []string, bottleMap map[string]int) (*fyne.Container, *widget.Entry, *fyne.Container) {
	// Region UPDATE FORM
	updateForm := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelName := widget.NewLabel("Nom")
	entryName := widget.NewEntry()
	//entryName.Disable()

	//StorageLocation's header
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
	gridContainerHeader := &fyne.Container{Layout: layout.NewGridLayout(3)}
	// List headers
	gridContainerHeader.Add(widget.NewLabel(""))
	gridContainerHeader.Add(labelBottle)
	gridContainerHeader.Add(labelQuantity)
	// List items
	gridContainerItems := &fyne.Container{Layout: layout.NewGridLayout(3)}

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

	return updateForm, entryName, gridContainerItems
}

func StorageLocation_InitButtonContainer(StorageLocation *data.StorageLocation, entryName *widget.Entry) *fyne.Container {

	var source = "WIDGETS.STORAGELOCATION.StorageLocation_InitButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier cet emplacement", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer cet emplacement", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {

		bottleStorageLocations := make([]data.BottleStorageLocation, 0)

		for control, _ := range BottleStorageLocationControls {
			bottle := data.Bottle{
				ID: control.BottleId,
			}

			quantity, _ := strconv.ParseInt(control.EntryQuantity.Text, 10, 0)

			bottleStorageLocation := data.BottleStorageLocation{
				Bottle:   bottle,
				Quantity: int(quantity),
			}

			bottleStorageLocations = append(bottleStorageLocations, bottleStorageLocation)
		}

		storageLocation := &data.StorageLocation{
			ID:                     StorageLocation.ID,
			Name:                   entryName.Text,
			BottleStorageLocations: bottleStorageLocations,
		}
		jsonValue, _ := json.Marshal(storageLocation)
		postData := data.AuthPostRequest("StorageLocation/UpdateStorageLocation", bytes.NewBuffer(jsonValue))
		if postData != 200 {
			fmt.Println("Error on update")
			message := "Error on storageLocation " + identifier + " update"
			log(true, source, message)
			return
		}
		fmt.Println("Success on update")
		StorageLocationTableRefreshMethod()
	}

	deleteBtn.OnTapped = func() {
		jsonValue, _ := json.Marshal(strconv.Itoa(StorageLocation.ID))

		postData := data.AuthPostRequest("StorageLocation/DeleteStorageLocation", bytes.NewBuffer(jsonValue))
		if postData != 200 {
			fmt.Println("Error on delete")
			message := "Error on storageLocation " + identifier + " delete"
			log(true, source, message)
			return
		}
		StorageLocationTableRefreshMethod()
		StorageLocationUpdateFormClearMethod()
	}

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion
	return buttonsContainer
}

func StorageLocation_InitMainContainer(updateForm *fyne.Container, table *widget.Table, buttonsContainer *fyne.Container) *fyne.Container {

	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithButtons := container.NewBorder(layoutUpdateForm, buttonsContainer, nil, nil)

	// Define layout
	individualTabs := container.NewAppTabs(
		container.NewTabItem("Modifier l'emplacement", layoutWithButtons),
	)

	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)
	mainContainer := container.New(layout.NewGridLayout(2))
	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

func StorageLocation_GetAndMapBottleNames() ([]string, map[string]int) {
	bottles := getAllBottleName()
	bottleNames := make([]string, 0)
	for i := 0; i < len(bottles); i++ {
		bottles[i].ID = strconv.Itoa(bottles[i].Id)
		name := bottles[i].FullName
		bottleNames = append(bottleNames, name)
	}

	bottleMap := make(map[string]int)
	for i := 0; i < len(bottles); i++ {
		id := bottles[i].Id
		name := bottles[i].FullName
		bottleMap[name] = id
	}

	return bottleNames, bottleMap
}

// endregion

// Form to add and send a new storageLocation to the API endpoint (POST)
func addNewStorageLocation(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.STORAGELOCATION.addNewStorageLocation"
	//nameLabel := widget.NewLabel("Nom")
	//nameStorageLocation := widget.NewEntry()
	//detailsLabel := widget.NewLabel("Description")
	//detailsStorageLocation := widget.NewMultiLineEntry()
	//pictureLabel := widget.NewLabel("Image")
	//pictureStorageLocation := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })
	//
	//title := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	//
	//form := &widget.Form{
	//	Items: []*widget.FormItem{
	//		{Text: "", Widget: title},
	//		{Text: "", Widget: nameLabel},
	//		{Text: "", Widget: nameStorageLocation},
	//		{Text: "", Widget: detailsLabel},
	//		{Text: "", Widget: detailsStorageLocation},
	//		{Text: "", Widget: pictureLabel},
	//		{Text: "", Widget: pictureStorageLocation},
	//	},
	//	OnSubmit: func() {
	//		storageLocation := &data.StorageLocation{
	//			Name: nameStorageLocation.Text,
	//		}
	//		// convert storageLocation struct to json
	//		jsonValue, err := json.Marshal(&storageLocation)
	//		if err != nil {
	//			fmt.Println(err)
	//			log(true, source, err.Error())
	//			return
	//		}
	//		postData := data.AuthPostRequest("StorageLocation/AddStorageLocation", bytes.NewBuffer(jsonValue))
	//		if postData != 200|201 {
	//			message := "StatusCode " + strconv.Itoa(postData)
	//			log(true, source, message)
	//			fmt.Println(message)
	//			return
	//		}
	//		fmt.Println("New storageLocation added with success")
	//		StorageLocationTableRefreshMethod()
	//	},
	//	SubmitText: "Envoyer",
	//}

	BottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)

	bottleNames, bottleMap := StorageLocation_GetAndMapBottleNames()

	addForm, entryName, gridContainerItems := StorageLocation_InitForm(bottleNames, bottleMap)

	StorageLocationAddFormClearMethod = func() {
		entryName.Text = ""
		entryName.Refresh()
		gridContainerItems.RemoveAll()
		BottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
	}

	addBtn := widget.NewButtonWithIcon("Ajouter cet emplacement", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		bottleStorageLocations := make([]data.BottleStorageLocation, 0)

		for control, _ := range BottleStorageLocationControls {
			bottle := data.Bottle{
				ID: control.BottleId,
			}

			quantity, _ := strconv.ParseInt(control.EntryQuantity.Text, 10, 0)

			bottleStorageLocation := data.BottleStorageLocation{
				Bottle:   bottle,
				Quantity: int(quantity),
			}

			bottleStorageLocations = append(bottleStorageLocations, bottleStorageLocation)
		}

		storageLocation := &data.StorageLocation{
			Name:                   entryName.Text,
			BottleStorageLocations: bottleStorageLocations,
		}
		jsonValue, _ := json.Marshal(storageLocation)
		postData := data.AuthPostRequest("StorageLocation/AddStorageLocation", bytes.NewBuffer(jsonValue))
		if postData != 200 {
			fmt.Println("Error on add")
			message := "Error on storageLocation " + identifier + " add"
			log(true, source, message)
			return
		}
		fmt.Println("Success on update")
		StorageLocationTableRefreshMethod()
	}

	buttonsContainer := container.NewHBox(addBtn)

	layoutForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), addForm))
	layoutWithButtons := container.NewBorder(layoutForm, buttonsContainer, nil, nil)

	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), layoutWithButtons))

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

	BottleStorageLocationControls[item] = len(BottleStorageLocationControls) + 1

	var deleteItemBtn *widget.Button
	deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
		func() {
			gridContainerItems.Remove(deleteItemBtn)
			gridContainerItems.Remove(item.SelectEntryBottle)
			gridContainerItems.Remove(item.EntryQuantity)
			delete(BottleStorageLocationControls, item)
		})

	gridContainerItems.Add(deleteItemBtn)
	gridContainerItems.Add(item.SelectEntryBottle)
	gridContainerItems.Add(item.EntryQuantity)
}
