package Supplier

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
	"sort"
	"strconv"
)

// region " declarations "

type editForm struct {
	form *fyne.Container

	entryName    *widget.Entry
	entryMail    *widget.Entry
	entryDetails *widget.Entry

	entryAddressLine1 *widget.Entry
	entryAddressLine2 *widget.Entry
	selectCity        *widget.Select

	formClear func()
}

var log = data.Logger
var identifier string

var cities []data.PartialCity

var bind []binding.DataMap
var filter func([]data.PartialSupplier) []data.PartialSupplier

var currentCityId int

var table *widget.Table
var tableOptions *rtable.TableOptions

var updateForm editForm
var addForm editForm

// endregion " declarations "

// region " constructor "

// MakePage function creates a new set of tabs
func MakePage(_ fyne.Window) fyne.CanvasObject {

	ListTab := container.NewTabItem("Liste des fournisseurs", initListTab(nil))
	addTab := container.NewTabItem("Ajouter un fournisseur", initAddTab(nil))
	tabs := container.NewAppTabs(
		ListTab,
		addTab,
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		if ti == ListTab {
			tableRefresh()
			updateForm.formClear()
		} else if ti == addTab {
			addForm.formClear()
		}
	}
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// endregion " constructor "

// region " design initializers "

// region " tabs "

// initListTab implements a dynamic table bound to an editing form
func initListTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.SUPPLIER.initListTab()"

	//region datas
	// retrieve structs from data package
	Supplier := data.IndSupplier

	resp, label := getSuppliers()
	if !resp {
		return label
	}

	cityNames, cityMap := getAndMapCityNames()

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

	updateForm = initForm(cityNames, cityMap)

	//region " design elements initialization "
	buttonsContainer := initButtonContainer(&Supplier)
	buttonsContainer.Hide()
	mainContainer := initMainContainer(updateForm.form, buttonsContainer)
	//endregion
	updateForm.form.Hide()

	updateForm.formClear = func() {
		updateForm.form.Hide()
		table.UnselectAll()

		var entryArray = []*widget.Entry{updateForm.entryName, updateForm.entryMail, updateForm.entryDetails, updateForm.entryAddressLine1, updateForm.entryAddressLine2}

		for _, entry := range entryArray {
			entry.Text = ""
			entry.Refresh()
		}

		updateForm.selectCity.Selected = " "
		updateForm.selectCity.Refresh()

		Supplier.ID = -1
		buttonsContainer.Hide()

		currentCityId = -1
	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		tableOnSelected(cell, Columns, &Supplier, cityNames, cityMap, buttonsContainer)
	}
	//endregion

	return mainContainer
}

// Form to add and send a new object to the API endpoint (POST)
func initAddTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.SUPPLIER.initAddTab"

	cityNames, cityMap := getAndMapCityNames()

	addForm = initForm(cityNames, cityMap)

	addForm.formClear = func() {

		var entryArray = []*widget.Entry{addForm.entryName, addForm.entryMail, addForm.entryDetails, addForm.entryAddressLine1, addForm.entryAddressLine2}

		for _, entry := range entryArray {
			entry.Text = ""
			entry.Refresh()
		}

		addForm.selectCity.Selected = " "
		addForm.selectCity.Refresh()

		currentCityId = -1
	}

	addBtn := widget.NewButtonWithIcon("Ajouter ce fournisseur", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		addSuppliers()
	}

	buttonsContainer := container.NewHBox(addBtn)

	layoutForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), addForm.form))
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
		container.NewTabItem("Modifier le fournisseur", layoutWithButtons),
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

func initForm(cityNames []string, cityMap map[string]int) editForm {
	form := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelName := widget.NewLabel("Nom")
	entryName := widget.NewEntry()

	labelMail := widget.NewLabel("Email")
	entryMail := widget.NewEntry()

	labelDetails := widget.NewLabel("Description")
	entryDetails := widget.NewEntry()

	labelAddress1 := widget.NewLabel("Adresse 1")
	entryAddress1 := widget.NewEntry()

	labelAddress2 := widget.NewLabel("Adresse 2")
	entryAddress2 := widget.NewEntry()

	labelZIPCode := widget.NewLabel("Code postal")
	entryZIPCode := widget.NewEntry()

	labelCity := widget.NewLabel("Ville")
	selectCity := widget.NewSelect(cityNames, func(s string) {
		currentCityId = cityMap[s]

		var zipCode int
		for _, n := range cities {
			if n.Id == currentCityId {
				zipCode = n.ZipCode
				break
			}
		}

		entryZIPCode.SetText(strconv.Itoa(zipCode))
	})
	selectCity.PlaceHolder = " "

	//Supplier's header
	layoutHeader := &fyne.Container{Layout: layout.NewFormLayout()}

	layoutHeader.Add(labelName)
	layoutHeader.Add(entryName)
	layoutHeader.Add(labelMail)
	layoutHeader.Add(entryMail)
	layoutHeader.Add(labelDetails)
	layoutHeader.Add(entryDetails)
	layoutHeader.Add(labelAddress1)
	layoutHeader.Add(entryAddress1)
	layoutHeader.Add(labelAddress2)
	layoutHeader.Add(entryAddress2)
	layoutHeader.Add(labelCity)
	layoutHeader.Add(selectCity)
	layoutHeader.Add(labelZIPCode)
	layoutHeader.Add(entryZIPCode)

	form.Add(layoutHeader)
	form.Add(widget.NewLabel(""))

	formStruct := editForm{
		form:              form,
		entryName:         entryName,
		entryMail:         entryMail,
		entryDetails:      entryDetails,
		entryAddressLine1: entryAddress1,
		entryAddressLine2: entryAddress2,
		selectCity:        selectCity,
	}

	return formStruct
}

func initButtonContainer(Supplier *data.Supplier) *fyne.Container {

	//var source = "WIDGETS.SUPPLIER.initButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier ce fournisseur", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer ce fournisseur", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {
		updateSupplier(Supplier)
	}

	deleteBtn.OnTapped = func() {
		deleteSupplier(Supplier.ID)
	}

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion
	return buttonsContainer
}

// endregion " containers "

// endregion " design initializers "

// region " data "

// region " suppliers "

func getSuppliers() (bool, *widget.Label) {
	var source = "WIDGETS.SUPPLIER.getSuppliers() "
	Suppliers := data.SupplierData
	response := data.AuthGetRequest("supplier")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&Suppliers); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	//filter data
	if filter != nil {
		Suppliers = filter(Suppliers)
	}

	//order datas by Id
	sort.SliceStable(Suppliers, func(i, j int) bool {
		return Suppliers[i].Id < Suppliers[j].Id
	})

	bind = nil

	for i := 0; i < len(Suppliers); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := Suppliers[i]
		id := strconv.Itoa(t.Id)
		Suppliers[i].ID = id

		// binding supplier data
		bind = append(bind, binding.BindStruct(&Suppliers[i]))

	}

	return true, widget.NewLabel("")
}

func addSuppliers() {
	var source = "WIDGETS.SUPPLIER.addSuppliers"

	name := addForm.entryName.Text
	mail := addForm.entryMail.Text
	details := addForm.entryDetails.Text

	address := &data.Address{
		AddressLine1: addForm.entryAddressLine1.Text,
		AddressLine2: addForm.entryAddressLine2.Text,
		CityId:       currentCityId,
	}

	supplier := &data.Supplier{
		Name:    name,
		Email:   mail,
		Address: address,
		Details: details,
	}
	jsonValue, _ := json.Marshal(supplier)
	postData := data.AuthPostRequest("Supplier/AddSupplier", bytes.NewBuffer(jsonValue))
	if postData != 201 {
		fmt.Println("Error on add")
		message := "Error on supplier " + identifier + " add"
		log(true, source, message)
		return
	} else {
		fmt.Println("Success on add")
	}
	tableRefresh()
}

func updateSupplier(Supplier *data.Supplier) {
	var source = "WIDGETS.SUPPLIER.updateSuppliers"

	name := updateForm.entryName.Text
	mail := updateForm.entryMail.Text
	details := updateForm.entryDetails.Text

	address := &data.Address{
		AddressLine1: updateForm.entryAddressLine1.Text,
		AddressLine2: updateForm.entryAddressLine2.Text,
		CityId:       currentCityId,
	}

	if Supplier.Address != nil {
		address.ID = Supplier.Address.ID
	}

	supplier := &data.Supplier{
		ID:      Supplier.ID,
		Name:    name,
		Email:   mail,
		Details: details,
		Address: address,
	}
	jsonValue, _ := json.Marshal(supplier)
	postData := data.AuthPostRequest("Supplier/UpdateSupplier", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on update")
		message := "Error on supplier " + identifier + " update"
		log(true, source, message)
		return
	}
	fmt.Println("Success on update")
	tableRefresh()
}

func deleteSupplier(id int) {
	var source = "WIDGETS.SUPPLIER.deleteSuppliers"
	jsonValue, _ := json.Marshal(strconv.Itoa(id))

	postData := data.AuthPostRequest("Supplier/DeleteSupplier", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on delete")
		message := "Error on supplier " + identifier + " delete"
		log(true, source, message)
		return
	}
	tableRefresh()
	updateForm.formClear()
}

// endregion " suppliers "

// region " cities "

func getAndMapCityNames() ([]string, map[string]int) {
	cities = getAllCityName()
	cityNames := make([]string, 0)
	for i := 0; i < len(cities); i++ {
		cities[i].ID = strconv.Itoa(cities[i].Id)
		name := cities[i].Name
		cityNames = append(cityNames, name)
	}

	cityMap := make(map[string]int)
	for i := 0; i < len(cities); i++ {
		id := cities[i].Id
		name := cities[i].Name
		cityMap[name] = id
	}

	return cityNames, cityMap
}

func getAllCityName() []data.PartialCity {
	var source = "WIDGETS.SUPPLIER.getAllCityName"
	cityData := data.CityData
	response := data.AuthGetRequest("City")
	if response == nil {
		fmt.Println("No result returned")
		return nil
	}
	if err := json.NewDecoder(response).Decode(&cityData); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return nil
	}
	return cityData
}

// endregion " cities "

// region " filters "

func beginByE(Suppliers []data.PartialSupplier) []data.PartialSupplier {

	n := 0
	for _, supplier := range Suppliers {
		if string([]rune(supplier.Name)[0]) == "e" {
			Suppliers[n] = supplier
			n++
		}
	}

	Suppliers = Suppliers[:n]

	return Suppliers
}

// endregion " filters "

// endregion " data "

// region " events "

// region " table "
func tableOnSelected(cell widget.TableCellID, Columns []rtable.ColAttr, Supplier *data.Supplier, cityNames []string, cityMap map[string]int, buttonsContainer *fyne.Container) {
	var source = "WIDGETS.SUPPLIER.tableOnSelected"
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
		// Fetch individual supplier to fill form
		response := data.AuthGetRequest("Supplier/" + identifier)
		if err := json.NewDecoder(response).Decode(&Supplier); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		updateForm.form.Show()
		buttonsContainer.Show()

		updateForm.entryName.SetText(Supplier.Name)
		updateForm.entryMail.SetText(Supplier.Email)
		updateForm.entryDetails.SetText(Supplier.Details)

		if Supplier != nil && Supplier.Address != nil {

			updateForm.entryAddressLine1.SetText(Supplier.Address.AddressLine1)
			updateForm.entryAddressLine2.SetText(Supplier.Address.AddressLine2)

			if Supplier.Address.City != nil {
				updateForm.selectCity.SetSelected(Supplier.Address.City.Name)
			} else {
				updateForm.selectCity.ClearSelected()
			}
		} else {
			updateForm.entryAddressLine1.SetText("")
			updateForm.entryAddressLine2.SetText("")
			updateForm.selectCity.ClearSelected()
		}

		updateForm.form.Refresh()
	} else {
		log(true, source, err.Error())
	}
}

func tableRefresh() {
	if table != nil && tableOptions != nil {
		getSuppliers()
		tableOptions.Bindings = bind
		table.Refresh()
	}
}

// endregion "table"

// endregion " events "

// region " Other functions "

// endregion " Other functions "
