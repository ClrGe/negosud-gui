package StorageLocation

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
	"negosud-gui/widgets/StorageLocation/controls"
	"sort"
	"strconv"
)

// region " declarations "

var log = data.Logger
var identifier string

var bind []binding.DataMap
var filter func([]data.PartialStorageLocation) []data.PartialStorageLocation

var table *widget.Table
var tableOptions *rtable.TableOptions

var bottleStorageLocationControls map[*controls.BottleStorageLocationItem]int

var updateFormClearMethod func()
var addFormClearMethod func()

// endregion " declarations "

// region " constructor "

// MakePage function creates a new set of tabs
func MakePage(_ fyne.Window) fyne.CanvasObject {

	ListTab := container.NewTabItem("Liste des emplacements", initListTab(nil))
	addTab := container.NewTabItem("Ajouter un emplacement", initAddTab(nil))
	tabs := container.NewAppTabs(
		ListTab,
		addTab,
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		if ti == ListTab {
			tableRefresh()
			updateFormClearMethod()
		} else if ti == addTab {
			addFormClearMethod()
		}
	}
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// endregion " constructor "

// region " design initializers "

// region " tabs "

// initListTab implements a dynamic table bound to an editing form
func initListTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.STORAGELOCATION.initListTab()"

	//region datas
	// retrieve structs from data package
	StorageLocation := data.IndStorageLocation

	resp, label := getStorageLocations()
	if !resp {
		return label
	}

	bottleNames, bottleMap := getAndMapBottleNames()

	// Columns defines the header row for the table
	var Columns = []rtable.ColAttr{
		{ColName: "ID", Header: "ID", WidthPercent: 50},
		{ColName: "Name", Header: "Nom", WidthPercent: 90},
	}

	tableOptions = &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: Columns,
		Bindings: bind,
	}
	table = rtable.CreateTable(tableOptions)

	//region UPDATE FORM

	updateForm, entryName, gridContainerItems := initForm(bottleNames, bottleMap)

	//region " design elements initialization "
	buttonsContainer := initButtonContainer(&StorageLocation, entryName)
	buttonsContainer.Hide()
	mainContainer := initMainContainer(updateForm, buttonsContainer)
	//endregion
	updateForm.Hide()

	updateFormClearMethod = func() {
		updateForm.Hide()
		table.UnselectAll()
		entryName.Text = ""
		entryName.Refresh()
		gridContainerItems.RemoveAll()
		StorageLocation.ID = -1
		bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
		buttonsContainer.Hide()
	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		tableOnSelected(cell, Columns, &StorageLocation, entryName, gridContainerItems, bottleNames, bottleMap, updateForm, buttonsContainer)
	}
	//endregion

	return mainContainer
}

// Form to add and send a new object to the API endpoint (POST)
func initAddTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.STORAGELOCATION.initAddTab"

	bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)

	bottleNames, bottleMap := getAndMapBottleNames()

	addForm, entryName, gridContainerItems := initForm(bottleNames, bottleMap)

	addFormClearMethod = func() {
		entryName.Text = ""
		entryName.Refresh()
		gridContainerItems.RemoveAll()
		bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
	}

	addBtn := widget.NewButtonWithIcon("Ajouter cet emplacement", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		addStorageLocations(entryName.Text)
	}

	buttonsContainer := container.NewHBox(addBtn)

	layoutForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), addForm))
	layoutWithButtons := container.NewBorder(layoutForm, buttonsContainer, nil, nil)

	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), layoutWithButtons))

	return mainContainer
}

// endregion " tabs "

// region " containers "

func initMainContainer(updateForm *fyne.Container, buttonsContainer *fyne.Container) *fyne.Container {

	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithButtons := container.NewBorder(layoutUpdateForm, buttonsContainer, nil, nil)

	// Define layout
	individualTabs := container.NewAppTabs(
		container.NewTabItem("Modifier l'emplacement", layoutWithButtons),
	)

	filterContainer := initFilterContainer()

	tableContainer := container.NewBorder(filterContainer, nil, nil, nil, table)

	leftContainer := tableContainer
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)
	mainContainer := container.New(layout.NewGridLayout(2))
	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

func initFilterContainer() *fyne.Container {
	filterStrings := []string{"Tous", "E"}

	selectFilter := widget.NewSelect(filterStrings, func(s string) {
		if s == "Tous" {
			filter = nil
		} else {
			filter = beginByE
		}
		tableRefresh()
	})

	selectFilter.PlaceHolder = "Selectionner un filtre"

	selectFilter.Selected = "Tous"

	label := widget.NewLabel("Filtre : ")

	border2 := container.NewBorder(nil, nil, label, nil, selectFilter)

	filterContainer := container.NewCenter(border2)

	return filterContainer
}

func initForm(bottleNames []string, bottleMap map[string]int) (*fyne.Container, *widget.Entry, *fyne.Container) {
	// Region UPDATE FORM
	updateForm := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelName := widget.NewLabel("Nom")
	entryName := widget.NewEntry()

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
			addBottleStorageLocationControl(bottleNames, bottleMap, gridContainerItems)
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

func initButtonContainer(StorageLocation *data.StorageLocation, entryName *widget.Entry) *fyne.Container {

	var source = "WIDGETS.STORAGELOCATION.initButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier cet emplacement", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer cet emplacement", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {
		updateStorageLocation(StorageLocation, entryName.Text)
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
		tableRefresh()
		updateFormClearMethod()
	}

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion
	return buttonsContainer
}

// endregion " containers "

// endregion " design initializers "

// region " data "

// region " storageLocations "

func getStorageLocations() (bool, *widget.Label) {
	var source = "WIDGETS.STORAGELOCATION.getStorageLocations() "
	StorageLocations := data.StorageLocationsData
	response := data.AuthGetRequest("storageLocation")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&StorageLocations); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	//filter data
	if filter != nil {
		StorageLocations = filter(StorageLocations)
	}

	//order datas by Id
	sort.SliceStable(StorageLocations, func(i, j int) bool {
		return StorageLocations[i].Id < StorageLocations[j].Id
	})

	bind = nil

	for i := 0; i < len(StorageLocations); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := StorageLocations[i]
		id := strconv.Itoa(t.Id)
		StorageLocations[i].ID = id

		// binding storageLocation data
		bind = append(bind, binding.BindStruct(&StorageLocations[i]))

	}

	return true, widget.NewLabel("")
}

func addStorageLocations(name string) {
	var source = "WIDGETS.STORAGELOCATION.addStorageLocations"
	bottleStorageLocations := make([]data.BottleStorageLocation, 0)

	uniqueIds := make(map[int]struct{})
	// Modify duplicate values to exclude them later
	for item, _ := range bottleStorageLocationControls {
		if _, has := uniqueIds[item.BottleId]; has {
			//duplicate = true
			item.BottleId = -1
		}
		uniqueIds[item.BottleId] = struct{}{}
	}

	for control, _ := range bottleStorageLocationControls {
		// Exclude duplicate values
		if control.BottleId > 0 {

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
	}

	storageLocation := &data.StorageLocation{
		Name:                   name,
		BottleStorageLocations: bottleStorageLocations,
	}
	jsonValue, _ := json.Marshal(storageLocation)
	postData := data.AuthPostRequest("StorageLocation/AddStorageLocation", bytes.NewBuffer(jsonValue))
	if postData != 201 {
		fmt.Println("Error on add")
		message := "Error on storageLocation " + identifier + " add"
		log(true, source, message)
		return
	} else {
		fmt.Println("Success on add")
	}
	tableRefresh()
}

func updateStorageLocation(StorageLocation *data.StorageLocation, name string) {
	var source = "WIDGETS.STORAGELOCATION.updateStorageLocations"
	bottleStorageLocations := make([]data.BottleStorageLocation, 0)

	for control, _ := range bottleStorageLocationControls {
		if control.BottleId > 0 {
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
	}

	storageLocation := &data.StorageLocation{
		ID:                     StorageLocation.ID,
		Name:                   name,
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
	tableRefresh()
}

// endregion " storageLocations "

// region " bottles "

func getAndMapBottleNames() ([]string, map[string]int) {
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

func getAllBottleName() []data.PartialBottle {
	var source = "WIDGETS.STORAGELOCATION.getAllBottleName"
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

// endregion " bottles "

// region " filters "

func beginByE(StorageLocations []data.PartialStorageLocation) []data.PartialStorageLocation {

	n := 0
	for _, storageLocation := range StorageLocations {
		if string([]rune(storageLocation.Name)[0]) == "e" {
			StorageLocations[n] = storageLocation
			n++
		}
	}

	StorageLocations = StorageLocations[:n]

	return StorageLocations
}

// endregion " filters "

// endregion " data "

// region " events "

// region " table "
func tableOnSelected(cell widget.TableCellID, Columns []rtable.ColAttr, StorageLocation *data.StorageLocation, entryName *widget.Entry, gridContainerItems *fyne.Container, bottleNames []string, bottleMap map[string]int, updateForm *fyne.Container, buttonsContainer *fyne.Container) {
	var source = "WIDGETS.STORAGELOCATION.tableOnSelected"
	if cell.Row < 0 || cell.Row > len(bind) { // 1st col is header
		fmt.Println("*-> Row out of limits")
		log(true, source, "*-> Row out of limits")
		return
	}
	if cell.Col < 0 || cell.Col >= len(Columns) {
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
		updateForm.Show()
		buttonsContainer.Show()

		entryName.SetText(StorageLocation.Name)

		gridContainerItems.RemoveAll()

		bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)

		//
		for _, bsl := range StorageLocation.BottleStorageLocations {
			item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)
			item.Bind(bsl.Bottle.ID, bsl.StorageLocation.ID)
			item.SelectBottle.Selected = bsl.Bottle.FullName
			item.EntryQuantity.Text = strconv.Itoa(bsl.Quantity)

			bottleStorageLocationControls[item] = item.BottleId

			var deleteItemBtn *widget.Button
			deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
				func() {
					gridContainerItems.Remove(deleteItemBtn)
					gridContainerItems.Remove(item.SelectBottle)
					gridContainerItems.Remove(item.EntryQuantity)
					delete(bottleStorageLocationControls, item)
				})

			gridContainerItems.Add(deleteItemBtn)
			gridContainerItems.Add(item.SelectBottle)
			gridContainerItems.Add(item.EntryQuantity)
		}
		updateForm.Refresh()
	} else {
		log(true, source, err.Error())
	}
}

func tableRefresh() {
	getStorageLocations()
	tableOptions.Bindings = bind
	table.Refresh()
}

// endregion "table"

// endregion " events "

// region " Other functions "

func addBottleStorageLocationControl(bottleNames []string, bottleMap map[string]int, gridContainerItems *fyne.Container) {
	item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)

	bottleStorageLocationControls[item] = len(bottleStorageLocationControls) + 1

	var deleteItemBtn *widget.Button
	deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
		func() {
			gridContainerItems.Remove(deleteItemBtn)
			gridContainerItems.Remove(item.SelectBottle)
			gridContainerItems.Remove(item.EntryQuantity)
			delete(bottleStorageLocationControls, item)
		})

	gridContainerItems.Add(deleteItemBtn)
	gridContainerItems.Add(item.SelectBottle)
	gridContainerItems.Add(item.EntryQuantity)
}

// endregion " Other functions "
