package Bottle

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
	"negosud-gui/widgets/Bottle/controls"
	"sort"
	"strconv"
	"strings"
)

// region " declarations "

var log = data.Logger
var identifier string

var bind []binding.DataMap
var filter func([]data.PartialBottle) []data.PartialBottle

var table *widget.Table
var tableOptions *rtable.TableOptions

var bottleStorageLocationControls map[*controls.BottleStorageLocationItem]int

var updateFormClearMethod func()
var addFormClearMethod func()

// endregion " declarations "

// region " constructor "

// MakePage function creates a new set of tabs
func MakePage(_ fyne.Window) fyne.CanvasObject {

	ListTab := container.NewTabItem("Liste des bouteilles", initListTab(nil))
	addTab := container.NewTabItem("Ajouter un bouteille", initAddTab(nil))
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
	//var source = "WIDGETS.Bottle.initListTab()"

	//region datas
	// retrieve structs from data package
	Bottle := data.IndBottle

	resp, label := getBottles()
	if !resp {
		return label
	}

	storageLocationNames, storageLocationMap := getAndMapStorageLocationNames()

	// Columns defines the header row for the table
	var Columns = []rtable.ColAttr{
		{ColName: "ID", Header: "ID", WidthPercent: 50},
		{ColName: "FullName", Header: "Nom", WidthPercent: 90},
	}

	tableOptions = &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: Columns,
		Bindings: bind,
	}
	table = rtable.CreateTable(tableOptions)

	//region UPDATE FORM

	updateForm, entryName, entryCustomerPrice, entrySupplierPrice, gridContainerItems := initForm(storageLocationNames, storageLocationMap)

	//region " design elements initialization "
	buttonsContainer := initButtonContainer(&Bottle, entryName, entryCustomerPrice, entrySupplierPrice)
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
		Bottle.ID = -1
		bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
		buttonsContainer.Hide()
	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		tableOnSelected(cell, Columns, &Bottle, entryName, entryCustomerPrice, entrySupplierPrice, gridContainerItems, storageLocationNames, storageLocationMap, updateForm, buttonsContainer)
	}
	//endregion

	return mainContainer
}

// Form to add and send a new object to the API endpoint (POST)
func initAddTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.Bottle.initAddTab"

	bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)

	bottleNames, bottleMap := getAndMapStorageLocationNames()

	addForm, entryName, entryCustomerPrice, entrySupplierPrice, gridContainerItems := initForm(bottleNames, bottleMap)

	addFormClearMethod = func() {
		entryName.Text = ""
		entryName.Refresh()
		entryCustomerPrice.Text = ""
		entryCustomerPrice.Refresh()
		entrySupplierPrice.Text = ""
		entrySupplierPrice.Refresh()
		gridContainerItems.RemoveAll()
		bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)
	}

	addBtn := widget.NewButtonWithIcon("Ajouter cette bouteille", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		addBottles(entryName.Text, entryCustomerPrice.Text, entrySupplierPrice.Text)
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
		container.NewTabItem("Modifier la bouteille", layoutWithButtons),
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

func initForm(bottleNames []string, bottleMap map[string]int) (*fyne.Container, *widget.Entry, *widget.Entry, *widget.Entry, *fyne.Container) {
	// Region UPDATE FORM
	updateForm := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelName := widget.NewLabel("Nom")
	entryName := widget.NewEntry()

	labelCustomerPrice := widget.NewLabel("Prix du client")
	entryCustomerPrice := widget.NewEntry()

	labelSupplierPrice := widget.NewLabel("Prix du fournisseur")
	entrySupplierPrice := widget.NewEntry()

	//Bottle's header
	layoutControlItemName := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutControlItemName.Add(labelName)
	layoutControlItemName.Add(entryName)

	layoutControlItemCustomerPrice := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutControlItemCustomerPrice.Add(labelCustomerPrice)
	layoutControlItemCustomerPrice.Add(entryCustomerPrice)

	layoutControlItemSupplierPrice := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutControlItemSupplierPrice.Add(labelSupplierPrice)
	layoutControlItemSupplierPrice.Add(entrySupplierPrice)

	//BottleStorageLocation List

	// List Title
	BSLListTitle := widget.NewLabel("Emplacements")
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
	updateForm.Add(layoutControlItemCustomerPrice)
	updateForm.Add(layoutControlItemSupplierPrice)
	updateForm.Add(widget.NewLabel(""))
	updateForm.Add(widget.NewSeparator())
	updateForm.Add(BSLListTitle)
	updateForm.Add(gridContainerHeader)
	updateForm.Add(gridContainerItems)
	updateForm.Add(AddItemBtn)

	return updateForm, entryName, entryCustomerPrice, entrySupplierPrice, gridContainerItems
}

func initButtonContainer(Bottle *data.Bottle, entryName *widget.Entry, entryCustomerPrice *widget.Entry, entrySupplierPrice *widget.Entry) *fyne.Container {

	var source = "WIDGETS.Bottle.initButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier cette bouteille", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer cette bouteille", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {
		updateBottle(Bottle, entryName.Text, entryCustomerPrice.Text, entrySupplierPrice.Text)
	}

	deleteBtn.OnTapped = func() {
		jsonValue, _ := json.Marshal(strconv.Itoa(Bottle.ID))

		postData := data.AuthPostRequest("Bottle/DeleteBottle", bytes.NewBuffer(jsonValue))
		if postData != 200 {
			fmt.Println("Error on delete")
			message := "Error on Bottle " + identifier + " delete"
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

// region " Bottles "

func getBottles() (bool, *widget.Label) {
	var source = "WIDGETS.Bottle.getBottles() "
	Bottles := data.BottleData
	response := data.AuthGetRequest("Bottle")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&Bottles); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	//filter data
	if filter != nil {
		Bottles = filter(Bottles)
	}

	//order datas by Id
	sort.SliceStable(Bottles, func(i, j int) bool {
		return Bottles[i].Id < Bottles[j].Id
	})

	bind = nil

	for i := 0; i < len(Bottles); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := Bottles[i]
		id := strconv.Itoa(t.Id)
		Bottles[i].ID = id

		// binding Bottle data
		bind = append(bind, binding.BindStruct(&Bottles[i]))

	}

	return true, widget.NewLabel("")
}

func addBottles(name string, customerPriceString string, supplierPriceString string) {
	var source = "WIDGETS.Bottle.addBottles"
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

	customerPriceString = strings.Replace(customerPriceString, ",", ".", 1)
	supplierPriceString = strings.Replace(supplierPriceString, ",", ".", 1)

	customerPrice, err := strconv.ParseFloat(customerPriceString, 32)
	if err != nil {
		customerPrice = 0
	}

	supplierPrice, err := strconv.ParseFloat(supplierPriceString, 32)
	if err != nil {
		supplierPrice = 0
	}

	Bottle := &data.Bottle{
		FullName:      name,
		CustomerPrice: float32(customerPrice),
		SupplierPrice: float32(supplierPrice),

		BottleStorageLocations: bottleStorageLocations,
	}
	jsonValue, _ := json.Marshal(Bottle)
	postData := data.AuthPostRequest("Bottle/AddBottle", bytes.NewBuffer(jsonValue))
	if postData != 201 {
		fmt.Println("Error on add")
		message := "Error on Bottle " + identifier + " add"
		log(true, source, message)
		return
	} else {
		fmt.Println("Success on add")
	}
	tableRefresh()
}

func updateBottle(Bottle *data.Bottle, name string, customerPriceString string, supplierPriceString string) {
	var source = "WIDGETS.Bottle.updateBottles"
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

	customerPriceString = strings.Replace(customerPriceString, ",", ".", 1)
	supplierPriceString = strings.Replace(supplierPriceString, ",", ".", 1)

	customerPrice, err := strconv.ParseFloat(customerPriceString, 32)
	if err != nil {
		customerPrice = 0
	}

	supplierPrice, err := strconv.ParseFloat(supplierPriceString, 32)
	if err != nil {
		supplierPrice = 0
	}

	bottle := &data.Bottle{
		ID:                     Bottle.ID,
		FullName:               name,
		CustomerPrice:          float32(customerPrice),
		SupplierPrice:          float32(supplierPrice),
		BottleStorageLocations: bottleStorageLocations,
	}
	jsonValue, _ := json.Marshal(bottle)
	postData := data.AuthPostRequest("Bottle/UpdateBottle", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on update")
		message := "Error on Bottle " + identifier + " update"
		log(true, source, message)
		return
	}
	fmt.Println("Success on update")
	tableRefresh()
}

// endregion " Bottles "

// region " bottles "

func getAndMapStorageLocationNames() ([]string, map[string]int) {
	storageLocations := getAllStorageLocationName()
	storageLocationNames := make([]string, 0)
	for i := 0; i < len(storageLocations); i++ {
		storageLocations[i].ID = strconv.Itoa(storageLocations[i].Id)
		name := storageLocations[i].Name
		storageLocationNames = append(storageLocationNames, name)
	}

	storageLocationMap := make(map[string]int)
	for i := 0; i < len(storageLocations); i++ {
		id := storageLocations[i].Id
		name := storageLocations[i].Name
		storageLocationMap[name] = id
	}

	return storageLocationNames, storageLocationMap
}

func getAllStorageLocationName() []data.PartialStorageLocation {
	var source = "WIDGETS.Bottle.getAllStorageLocationName"
	storageLocationData := data.StorageLocationData
	response := data.AuthGetRequest("StorageLocation")
	if response == nil {
		fmt.Println("No result returned")
		return nil
	}
	if err := json.NewDecoder(response).Decode(&storageLocationData); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return nil
	}
	return storageLocationData
}

// endregion " bottles "

// region " filters "

func beginByE(Bottles []data.PartialBottle) []data.PartialBottle {

	n := 0
	for _, Bottle := range Bottles {
		if string([]rune(Bottle.FullName)[0]) == "e" {
			Bottles[n] = Bottle
			n++
		}
	}

	Bottles = Bottles[:n]

	return Bottles
}

// endregion " filters "

// endregion " data "

// region " events "

// region " table "
func tableOnSelected(cell widget.TableCellID, Columns []rtable.ColAttr, Bottle *data.Bottle, entryName *widget.Entry, entryCustomerPrice *widget.Entry, entrySupplierPrice *widget.Entry, gridContainerItems *fyne.Container, bottleNames []string, bottleMap map[string]int, updateForm *fyne.Container, buttonsContainer *fyne.Container) {
	var source = "WIDGETS.Bottle.tableOnSelected"
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
		// Fetch individual Bottle to fill form
		response := data.AuthGetRequest("Bottle/" + identifier)
		if err := json.NewDecoder(response).Decode(&Bottle); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		updateForm.Show()
		buttonsContainer.Show()

		entryName.SetText(Bottle.FullName)
		entryCustomerPrice.SetText(fmt.Sprintf("%f", Bottle.CustomerPrice))
		entrySupplierPrice.SetText(fmt.Sprintf("%f", Bottle.SupplierPrice))

		gridContainerItems.RemoveAll()

		bottleStorageLocationControls = make(map[*controls.BottleStorageLocationItem]int)

		//
		for _, bsl := range Bottle.BottleStorageLocations {
			item := controls.NewBottleStorageLocationControl(bottleNames, bottleMap)
			item.Bind(bsl.Bottle.ID, bsl.Bottle.ID)
			item.SelectStorageLocation.Selected = bsl.StorageLocation.Name
			item.EntryQuantity.Text = strconv.Itoa(bsl.Quantity)

			bottleStorageLocationControls[item] = item.BottleId

			var deleteItemBtn *widget.Button
			deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
				func() {
					gridContainerItems.Remove(deleteItemBtn)
					gridContainerItems.Remove(item.SelectStorageLocation)
					gridContainerItems.Remove(item.EntryQuantity)
					delete(bottleStorageLocationControls, item)
				})

			gridContainerItems.Add(deleteItemBtn)
			gridContainerItems.Add(item.SelectStorageLocation)
			gridContainerItems.Add(item.EntryQuantity)
		}
		updateForm.Refresh()
	} else {
		log(true, source, err.Error())
	}
}

func tableRefresh() {
	getBottles()
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
			gridContainerItems.Remove(item.SelectStorageLocation)
			gridContainerItems.Remove(item.EntryQuantity)
			delete(bottleStorageLocationControls, item)
		})

	gridContainerItems.Add(deleteItemBtn)
	gridContainerItems.Add(item.SelectStorageLocation)
	gridContainerItems.Add(item.EntryQuantity)
}

// endregion " Other functions "
