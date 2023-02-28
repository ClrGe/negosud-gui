package WineLabel

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

	entryLabel *widget.Entry

	formClear func()
}

var log = data.Logger
var identifier string

var bind []binding.DataMap
var filter func([]data.PartialWineLabel) []data.PartialWineLabel

var table *widget.Table
var tableOptions *rtable.TableOptions

var updateForm editForm
var addForm editForm

// endregion " declarations "

// region " constructor "

// MakePage function creates a new set of tabs
func MakePage(_ fyne.Window) fyne.CanvasObject {

	ListTab := container.NewTabItem("Liste des labels", initListTab(nil))
	addTab := container.NewTabItem("Ajouter un label", initAddTab(nil))
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
	//var source = "WIDGETS.WineLabel.initListTab()"

	//region datas
	// retrieve structs from data package
	WineLabel := data.IndWineLabel

	resp, label := getWineLabels()
	if !resp {
		return label
	}

	// Columns defines the header row for the table
	var Columns = []rtable.ColAttr{
		{ColName: "ID", Header: "ID", WidthPercent: 50},
		{ColName: "Label", Header: "Nom", WidthPercent: 90},
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
	buttonsContainer := initButtonContainer(&WineLabel)
	buttonsContainer.Hide()
	mainContainer := initMainContainer(updateForm.form, buttonsContainer)
	//endregion
	updateForm.form.Hide()

	updateForm.formClear = func() {
		updateForm.form.Hide()
		table.UnselectAll()
		updateForm.entryLabel.Text = ""
		updateForm.entryLabel.Refresh()
		WineLabel.ID = -1
		buttonsContainer.Hide()
	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		tableOnSelected(cell, Columns, &WineLabel, buttonsContainer)
	}
	//endregion

	return mainContainer
}

// Form to add and send a new object to the API endpoint (POST)
func initAddTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.WineLabel.initAddTab"

	addForm = initForm()

	addForm.formClear = func() {
		addForm.entryLabel.Text = ""
		addForm.entryLabel.Refresh()
	}

	addBtn := widget.NewButtonWithIcon("Ajouter ce label", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		addWineLabels()
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
		container.NewTabItem("Modifier le label", layoutWithButtons),
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
	labelLabel := widget.NewLabel("Nom")
	entryLabel := widget.NewEntry()

	//WineLabel's header
	layoutHeader := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutHeader.Add(labelLabel)
	layoutHeader.Add(entryLabel)

	form.Add(layoutHeader)
	form.Add(widget.NewLabel(""))
	form.Add(widget.NewSeparator())

	formStruct := editForm{
		form:       form,
		entryLabel: entryLabel,
	}

	return formStruct
}

func initButtonContainer(WineLabel *data.WineLabel) *fyne.Container {

	//var source = "WIDGETS.WineLabel.initButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier ce label ", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer ce label", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {
		updateWineLabel(WineLabel)
	}

	deleteBtn.OnTapped = func() {
		deleteWineLabel(WineLabel.ID)
	}

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion
	return buttonsContainer
}

// endregion " containers "

// endregion " design initializers "

// region " data "

// region " wineLabels "

func getWineLabels() (bool, *widget.Label) {
	var source = "WIDGETS.WineLabel.getWineLabels() "
	WineLabels := data.WineLabelData
	response := data.AuthGetRequest("wineLabel")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&WineLabels); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	//filter data
	if filter != nil {
		WineLabels = filter(WineLabels)
	}

	//order datas by Id
	sort.SliceStable(WineLabels, func(i, j int) bool {
		return WineLabels[i].Id < WineLabels[j].Id
	})

	bind = nil

	for i := 0; i < len(WineLabels); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := WineLabels[i]
		id := strconv.Itoa(t.Id)
		WineLabels[i].ID = id

		// binding wineLabel data
		bind = append(bind, binding.BindStruct(&WineLabels[i]))

	}

	return true, widget.NewLabel("")
}

func addWineLabels() {
	var source = "WIDGETS.WineLabel.addWineLabels"

	label := addForm.entryLabel.Text

	wineLabel := &data.WineLabel{
		Label: label,
	}
	jsonValue, _ := json.Marshal(wineLabel)
	postData := data.AuthPostRequest("WineLabel/AddWineLabel", bytes.NewBuffer(jsonValue))
	if postData != 201 {
		fmt.Println("Error on add")
		message := "Error on wineLabel " + identifier + " add"
		log(true, source, message)
		return
	} else {
		fmt.Println("Success on add")
	}
	tableRefresh()
}

func updateWineLabel(WineLabel *data.WineLabel) {
	var source = "WIDGETS.WineLabel.updateWineLabels"

	label := updateForm.entryLabel.Text

	wineLabel := &data.WineLabel{
		ID:    WineLabel.ID,
		Label: label,
	}
	jsonValue, _ := json.Marshal(wineLabel)
	postData := data.AuthPostRequest("WineLabel/UpdateWineLabel", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on update")
		message := "Error on wineLabel " + identifier + " update"
		log(true, source, message)
		return
	}
	fmt.Println("Success on update")
	tableRefresh()
}

func deleteWineLabel(id int) {
	var source = "WIDGETS.WineLabel.deleteWineLabels"
	jsonValue, _ := json.Marshal(strconv.Itoa(id))

	postData := data.AuthPostRequest("WineLabel/DeleteWineLabel", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on delete")
		message := "Error on wineLabel " + identifier + " delete"
		log(true, source, message)
		return
	}
	tableRefresh()
	updateForm.formClear()
}

// endregion " wineLabels "

// region " filters "

func beginByE(WineLabels []data.PartialWineLabel) []data.PartialWineLabel {

	n := 0
	for _, wineLabel := range WineLabels {
		if string([]rune(wineLabel.Label)[0]) == "e" {
			WineLabels[n] = wineLabel
			n++
		}
	}

	WineLabels = WineLabels[:n]

	return WineLabels
}

// endregion " filters "

// endregion " data "

// region " events "

// region " table "
func tableOnSelected(cell widget.TableCellID, Columns []rtable.ColAttr, WineLabel *data.WineLabel, buttonsContainer *fyne.Container) {
	var source = "WIDGETS.WineLabel.tableOnSelected"
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
		// Fetch individual wineLabel to fill form
		response := data.AuthGetRequest("WineLabel/" + identifier)
		if err := json.NewDecoder(response).Decode(&WineLabel); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		updateForm.form.Show()
		buttonsContainer.Show()

		updateForm.entryLabel.SetText(WineLabel.Label)

		updateForm.form.Refresh()
	} else {
		log(true, source, err.Error())
	}
}

func tableRefresh() {
	if table != nil && tableOptions != nil {
		getWineLabels()
		tableOptions.Bindings = bind
		table.Refresh()
	}
}

// endregion "table"

// endregion " events "

// region " Other functions "

// endregion " Other functions "
