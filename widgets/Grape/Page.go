package Grape

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

	entryGrapeType *widget.Entry

	formClear func()
}

var log = data.Logger
var identifier string

var bind []binding.DataMap
var filter func([]data.PartialGrape) []data.PartialGrape

var table *widget.Table
var tableOptions *rtable.TableOptions

var updateForm editForm
var addForm editForm

// endregion " declarations "

// region " constructor "

// MakePage function creates a new set of tabs
func MakePage(_ fyne.Window) fyne.CanvasObject {

	ListTab := container.NewTabItem("Liste des cepages", initListTab(nil))
	addTab := container.NewTabItem("Ajouter un cepage", initAddTab(nil))
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
	//var source = "WIDGETS.Grape.initListTab()"

	//region datas
	// retrieve structs from data package
	Grape := data.IndGrape

	resp, grape := getGrapes()
	if !resp {
		return grape
	}

	// Columns defines the header row for the table
	var Columns = []rtable.ColAttr{
		{ColName: "ID", Header: "ID", WidthPercent: 50},
		{ColName: "GrapeType", Header: "Nom", WidthPercent: 90},
	}

	tableOptions = &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: Columns,
		Bindings: bind,
	}
	table = rtable.CreateTable(tableOptions)

	//region UPDATE FORM

	updateForm = initForm()

	//region " design elements initialization "
	buttonsContainer := initButtonContainer(&Grape)
	buttonsContainer.Hide()
	mainContainer := initMainContainer(updateForm.form, buttonsContainer)
	//endregion
	updateForm.form.Hide()

	updateForm.formClear = func() {
		updateForm.form.Hide()
		table.UnselectAll()
		updateForm.entryGrapeType.Text = ""
		updateForm.entryGrapeType.Refresh()
		Grape.ID = -1
		buttonsContainer.Hide()

	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		tableOnSelected(cell, Columns, &Grape, buttonsContainer)
	}
	//endregion

	return mainContainer
}

// Form to add and send a new object to the API endpoint (POST)
func initAddTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.Grape.initAddTab"

	addForm = initForm()

	addForm.formClear = func() {
		addForm.entryGrapeType.Text = ""
		addForm.entryGrapeType.Refresh()
	}

	addBtn := widget.NewButtonWithIcon("Ajouter ce cépage", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		addGrapes()
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
		container.NewTabItem("Modifier le cépage", layoutWithButtons),
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

func initForm() editForm {
	form := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelGrapeType := widget.NewLabel("Nom")
	entryGrapeType := widget.NewEntry()

	//Grape's header
	layoutHeader := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutHeader.Add(labelGrapeType)
	layoutHeader.Add(entryGrapeType)

	form.Add(layoutHeader)
	form.Add(widget.NewLabel(""))
	form.Add(widget.NewSeparator())

	formStruct := editForm{
		form:           form,
		entryGrapeType: entryGrapeType,
	}

	return formStruct
}

func initButtonContainer(Grape *data.Grape) *fyne.Container {

	//var source = "WIDGETS.Grape.initButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier ce cepageGrapeType ", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer ce cepage", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {
		updateGrape(Grape)
	}

	deleteBtn.OnTapped = func() {
		deleteGrape(Grape.ID)
	}

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion
	return buttonsContainer
}

// endregion " containers "

// endregion " design initializers "

// region " data "

// region " grapes "

func getGrapes() (bool, *widget.Label) {
	var source = "WIDGETS.Grape.getGrapes() "
	Grapes := data.GrapeData
	response := data.AuthGetRequest("grape")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&Grapes); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	//filter data
	if filter != nil {
		Grapes = filter(Grapes)
	}

	//order datas by Id
	sort.SliceStable(Grapes, func(i, j int) bool {
		return Grapes[i].Id < Grapes[j].Id
	})

	bind = nil

	for i := 0; i < len(Grapes); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := Grapes[i]
		id := strconv.Itoa(t.Id)
		Grapes[i].ID = id

		// binding grape data
		bind = append(bind, binding.BindStruct(&Grapes[i]))

	}

	return true, widget.NewLabel("")
}

func addGrapes() {
	var source = "WIDGETS.Grape.addGrapes"

	grapeType := addForm.entryGrapeType.Text

	grape := &data.Grape{
		GrapeType: grapeType,
	}
	jsonValue, _ := json.Marshal(grape)
	postData := data.AuthPostRequest("Grape/AddGrape", bytes.NewBuffer(jsonValue))
	if postData != 201 {
		fmt.Println("Error on add")
		message := "Error on grape " + identifier + " add"
		log(true, source, message)
		return
	} else {
		fmt.Println("Success on add")
	}
	tableRefresh()
}

func updateGrape(Grape *data.Grape) {
	var source = "WIDGETS.Grape.updateGrapes"

	grapeType := updateForm.entryGrapeType.Text

	grape := &data.Grape{
		ID:        Grape.ID,
		GrapeType: grapeType,
	}
	jsonValue, _ := json.Marshal(grape)
	postData := data.AuthPostRequest("Grape/UpdateGrape", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on update")
		message := "Error on grape " + identifier + " update"
		log(true, source, message)
		return
	}
	fmt.Println("Success on update")
	tableRefresh()
}

func deleteGrape(id int) {
	var source = "WIDGETS.Grape.deleteGrapes"
	jsonValue, _ := json.Marshal(strconv.Itoa(id))

	postData := data.AuthPostRequest("Grape/DeleteGrape", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on delete")
		message := "Error on grape " + identifier + " delete"
		log(true, source, message)
		return
	}
	tableRefresh()
	updateForm.formClear()
}

// endregion " grapes "

// region " filters "

func beginByE(Grapes []data.PartialGrape) []data.PartialGrape {

	n := 0
	for _, grape := range Grapes {
		if string([]rune(grape.GrapeType)[0]) == "e" {
			Grapes[n] = grape
			n++
		}
	}

	Grapes = Grapes[:n]

	return Grapes
}

// endregion " filters "

// endregion " data "

// region " events "

// region " table "
func tableOnSelected(cell widget.TableCellID, Columns []rtable.ColAttr, Grape *data.Grape, buttonsContainer *fyne.Container) {
	var source = "WIDGETS.Grape.tableOnSelected"
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
		// Fetch individual grape to fill form
		response := data.AuthGetRequest("Grape/" + identifier)
		if err := json.NewDecoder(response).Decode(&Grape); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		updateForm.form.Show()
		buttonsContainer.Show()

		updateForm.entryGrapeType.SetText(Grape.GrapeType)

		updateForm.form.Refresh()
	} else {
		log(true, source, err.Error())
	}
}

func tableRefresh() {
	if table != nil && tableOptions != nil {
		getGrapes()
		tableOptions.Bindings = bind
		table.Refresh()
	}
}

// endregion "table"

// endregion " events "

// region " Other functions "

// endregion " Other functions "
