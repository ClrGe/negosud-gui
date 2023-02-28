package SupplierOrder

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
	"negosud-gui/widgets/SupplierOrder/controls"
	"sort"
	"strconv"
)

// region " declarations "

type editForm struct {
	form *fyne.Container

	entryReference *widget.Entry
	selectSupplier *widget.Select

	gridContainerItems *fyne.Container

	formClear func()
}

var log = data.Logger
var identifier string

var suppliers []data.PartialSupplier

var bind []binding.DataMap
var filter func([]data.PartialSupplierOrder) []data.PartialSupplierOrder

var currentSupplierId int
var supplierNames []string
var supplierMap map[string]int

var table *widget.Table
var tableOptions *rtable.TableOptions

var updateForm editForm
var addForm editForm

var supplierOrderLineControls map[*controls.SupplierOrderLineItem]int

// endregion " declarations "

// region " constructor "

// MakePage function creates a new set of tabs
func MakePage(_ fyne.Window) fyne.CanvasObject {

	ListTab := container.NewTabItem("Liste des commandes fournisseurs", initListTab(nil))
	addTab := container.NewTabItem("Ajouter une commande fournisseur", initAddTab(nil))
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
	//var source = "WIDGETS.SupplierOrder.initListTab()"

	//region datas
	// retrieve structs from data package
	SupplierOrder := data.IndSupplierOrder

	resp, label := getSupplierOrders()
	if !resp {
		return label
	}

	bottleNames, bottleMap := getAndMapBottleNames()
	supplierNames, supplierMap = getAndMapSupplierNames()

	// Columns defines the header row for the table
	var Columns = []rtable.ColAttr{
		{ColName: "ID", Header: "ID", WidthPercent: 50},
		{ColName: "Reference", Header: "Reference", WidthPercent: 90},
	}

	tableOptions = &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: Columns,
		Bindings: bind,
	}
	table = rtable.CreateTable(tableOptions)

	//region UPDATE FORM

	updateForm = initForm(bottleNames, bottleMap)

	//region " design elements initialization "
	buttonsContainer := initButtonContainer(&SupplierOrder)
	buttonsContainer.Hide()
	mainContainer := initMainContainer(updateForm.form, buttonsContainer)
	//endregion
	updateForm.form.Hide()

	updateForm.formClear = func() {
		updateForm.form.Hide()
		table.UnselectAll()
		updateForm.entryReference.Text = ""
		updateForm.entryReference.Refresh()
		updateForm.selectSupplier.Selected = " "
		updateForm.selectSupplier.Refresh()
		updateForm.gridContainerItems.RemoveAll()
		SupplierOrder.ID = -1
		supplierOrderLineControls = make(map[*controls.SupplierOrderLineItem]int)
		buttonsContainer.Hide()

		currentSupplierId = -1
	}

	//endregion

	//region " table events "
	table.OnSelected = func(cell widget.TableCellID) {
		tableOnSelected(cell, Columns, &SupplierOrder, bottleNames, bottleMap, buttonsContainer)
	}
	//endregion

	return mainContainer
}

// Form to add and send a new object to the API endpoint (POST)
func initAddTab(_ fyne.Window) fyne.CanvasObject {
	//var source = "WIDGETS.SupplierOrder.initAddTab"

	supplierOrderLineControls = make(map[*controls.SupplierOrderLineItem]int)

	bottleNames, bottleMap := getAndMapBottleNames()
	supplierNames, supplierMap = getAndMapSupplierNames()

	addForm = initForm(bottleNames, bottleMap)

	addForm.formClear = func() {
		addForm.entryReference.Text = ""
		addForm.entryReference.Refresh()
		addForm.selectSupplier.Selected = " "
		addForm.selectSupplier.Refresh()
		addForm.gridContainerItems.RemoveAll()
		supplierOrderLineControls = make(map[*controls.SupplierOrderLineItem]int)
		currentSupplierId = -1
	}

	addBtn := widget.NewButtonWithIcon("Ajouter cette commande fournisseur", theme.ConfirmIcon(),
		func() {})

	addBtn.OnTapped = func() {
		addSupplierOrders()
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
		container.NewTabItem("Modifier la commande fournisseur", layoutWithButtons),
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

func initForm(bottleNames []string, bottleMap map[string]int) editForm {
	form := &fyne.Container{Layout: layout.NewVBoxLayout()}

	// declare form elements
	labelReference := widget.NewLabel("Reference")
	entryReference := widget.NewEntry()

	labelSupplier := widget.NewLabel("Fournisseur")
	selectSupplier := widget.NewSelect(supplierNames, func(s string) {
		currentSupplierId = supplierMap[s]
	})
	selectSupplier.PlaceHolder = " "

	//SupplierOrder's header
	layoutHeader := &fyne.Container{Layout: layout.NewFormLayout()}
	layoutHeader.Add(labelReference)
	layoutHeader.Add(entryReference)
	layoutHeader.Add(labelSupplier)
	layoutHeader.Add(selectSupplier)

	//SupplierOrderLine List

	// List Title
	BSLListTitle := widget.NewLabel("Produit")
	BSLListTitle.TextStyle.Bold = true

	// List headers
	labelBottle := widget.NewLabel("Reference")
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
			addSupplierOrderLineControl(bottleNames, bottleMap, gridContainerItems)
		})

	form.Add(layoutHeader)
	form.Add(widget.NewLabel(""))
	form.Add(widget.NewSeparator())
	form.Add(BSLListTitle)
	form.Add(gridContainerHeader)
	form.Add(gridContainerItems)
	form.Add(AddItemBtn)

	formStruct := editForm{
		form:               form,
		entryReference:     entryReference,
		selectSupplier:     selectSupplier,
		gridContainerItems: gridContainerItems,
	}

	return formStruct
}

func initButtonContainer(SupplierOrder *data.SupplierOrder) *fyne.Container {

	//var source = "WIDGETS.SupplierOrder.initButtonContainer"

	editBtn := widget.NewButtonWithIcon("Modifier cette commande fournisseur", theme.ConfirmIcon(),
		func() {})
	deleteBtn := widget.NewButtonWithIcon("Supprimer cette commande fournisseur", theme.WarningIcon(),
		func() {})

	//region " events "
	editBtn.OnTapped = func() {
		updateSupplierOrder(SupplierOrder)
	}

	deleteBtn.OnTapped = func() {
		deleteSupplierOrder(SupplierOrder.ID)
	}

	buttonsContainer := container.NewHBox(editBtn, deleteBtn)
	//endregion
	return buttonsContainer
}

// endregion " containers "

// endregion " design initializers "

// region " data "

// region " supplierOrders "

func getSupplierOrders() (bool, *widget.Label) {
	var source = "WIDGETS.SupplierOrder.getSupplierOrders() "
	SupplierOrders := data.SupplierOrderData
	response := data.AuthGetRequest("supplierOrder")
	if response == nil {
		fmt.Println("No result returned")
		return false, widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&SupplierOrders); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return false, widget.NewLabel("Erreur de décodage du json")
	}

	//filter data
	if filter != nil {
		SupplierOrders = filter(SupplierOrders)
	}

	//order datas by Id
	sort.SliceStable(SupplierOrders, func(i, j int) bool {
		return SupplierOrders[i].Id < SupplierOrders[j].Id
	})

	bind = nil

	for i := 0; i < len(SupplierOrders); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := SupplierOrders[i]
		id := strconv.Itoa(t.Id)
		SupplierOrders[i].ID = id

		// binding supplierOrder data
		bind = append(bind, binding.BindStruct(&SupplierOrders[i]))

	}

	return true, widget.NewLabel("")
}

func addSupplierOrders() {
	var source = "WIDGETS.SupplierOrder.addSupplierOrders"
	supplierOrderLines := make([]data.SupplierOrderLine, 0)

	uniqueIds := make(map[int]struct{})
	// Modify duplicate values to exclude them later
	for item, _ := range supplierOrderLineControls {
		if _, has := uniqueIds[item.BottleId]; has {
			//duplicate = true
			item.BottleId = -1
		}
		uniqueIds[item.BottleId] = struct{}{}
	}

	for control, _ := range supplierOrderLineControls {
		// Exclude duplicate values
		if control.BottleId > 0 {

			bottle := data.Bottle{
				ID: control.BottleId,
			}

			quantity, _ := strconv.ParseInt(control.EntryQuantity.Text, 10, 0)

			supplierOrderLine := data.SupplierOrderLine{
				BottleId: control.BottleId,
				Bottle:   bottle,
				Quantity: int(quantity),
			}

			supplierOrderLines = append(supplierOrderLines, supplierOrderLine)
		}
	}

	reference := addForm.entryReference.Text

	supplier := &data.Supplier{
		ID: currentSupplierId,
	}

	supplierOrder := &data.SupplierOrder{
		Reference: reference,
		Lines:     supplierOrderLines,
		Supplier:  supplier,
	}
	jsonValue, _ := json.Marshal(supplierOrder)
	postData := data.AuthPostRequest("SupplierOrder/AddSupplierOrder", bytes.NewBuffer(jsonValue))
	if postData != 201 {
		fmt.Println("Error on add")
		message := "Error on supplierOrder " + identifier + " add"
		log(true, source, message)
		return
	} else {
		fmt.Println("Success on add")
	}
	tableRefresh()
}

func updateSupplierOrder(SupplierOrder *data.SupplierOrder) {
	var source = "WIDGETS.SupplierOrder.updateSupplierOrders"
	supplierOrderLines := make([]data.SupplierOrderLine, 0)

	for control, _ := range supplierOrderLineControls {
		if control.BottleId > 0 {
			bottle := data.Bottle{
				ID: control.BottleId,
			}

			quantity, _ := strconv.ParseInt(control.EntryQuantity.Text, 10, 0)

			supplierOrderLine := data.SupplierOrderLine{
				BottleId: control.BottleId,
				Bottle:   bottle,
				Quantity: int(quantity),
			}

			supplierOrderLines = append(supplierOrderLines, supplierOrderLine)
		}
	}

	reference := updateForm.entryReference.Text

	supplier := &data.Supplier{
		ID: currentSupplierId,
	}

	supplierOrder := &data.SupplierOrder{
		ID:        SupplierOrder.ID,
		Reference: reference,
		Supplier:  supplier,
		Lines:     supplierOrderLines,
	}
	jsonValue, _ := json.Marshal(supplierOrder)
	postData := data.AuthPostRequest("SupplierOrder/UpdateSupplierOrder", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on update")
		message := "Error on supplierOrder " + identifier + " update"
		log(true, source, message)
		return
	}
	fmt.Println("Success on update")
	tableRefresh()
}

func deleteSupplierOrder(id int) {
	var source = "WIDGETS.SupplierOrder.deleteSupplierOrders"
	jsonValue, _ := json.Marshal(strconv.Itoa(id))

	postData := data.AuthPostRequest("SupplierOrder/DeleteSupplierOrder", bytes.NewBuffer(jsonValue))
	if postData != 200 {
		fmt.Println("Error on delete")
		message := "Error on supplierOrder " + identifier + " delete"
		log(true, source, message)
		return
	}
	tableRefresh()
	updateForm.formClear()
}

// endregion " supplierOrders "

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
	var source = "WIDGETS.SupplierOrder.getAllBottleName"
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

// region " suppliers "

func getAndMapSupplierNames() ([]string, map[string]int) {
	suppliers = getAllSupplierName()
	supplierNames = make([]string, 0)
	for i := 0; i < len(suppliers); i++ {
		suppliers[i].ID = strconv.Itoa(suppliers[i].Id)
		name := suppliers[i].Name
		supplierNames = append(supplierNames, name)
	}

	supplierMap = make(map[string]int)
	for i := 0; i < len(suppliers); i++ {
		id := suppliers[i].Id
		name := suppliers[i].Name
		supplierMap[name] = id
	}

	return supplierNames, supplierMap
}

func getAllSupplierName() []data.PartialSupplier {
	var source = "WIDGETS.SUPPLIER.getAllSupplierName"
	supplierData := data.SupplierData
	response := data.AuthGetRequest("Supplier")
	if response == nil {
		fmt.Println("No result returned")
		return nil
	}
	if err := json.NewDecoder(response).Decode(&supplierData); err != nil {
		fmt.Println(err)
		log(true, source, err.Error())
		return nil
	}
	return supplierData
}

// endregion " suppliers "

// region " filters "

func beginByE(SupplierOrders []data.PartialSupplierOrder) []data.PartialSupplierOrder {

	n := 0
	for _, supplierOrder := range SupplierOrders {
		if string([]rune(supplierOrder.Reference)[0]) == "e" {
			SupplierOrders[n] = supplierOrder
			n++
		}
	}

	SupplierOrders = SupplierOrders[:n]

	return SupplierOrders
}

// endregion " filters "

// endregion " data "

// region " events "

// region " table "
func tableOnSelected(cell widget.TableCellID, Columns []rtable.ColAttr, SupplierOrder *data.SupplierOrder, bottleNames []string, bottleMap map[string]int, buttonsContainer *fyne.Container) {
	var source = "WIDGETS.SupplierOrder.tableOnSelected"
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
		// Fetch individual supplierOrder to fill form
		response := data.AuthGetRequest("SupplierOrder/" + identifier)
		if err := json.NewDecoder(response).Decode(&SupplierOrder); err != nil {
			log(true, source, err.Error())
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		updateForm.form.Show()
		buttonsContainer.Show()

		updateForm.entryReference.SetText(SupplierOrder.Reference)

		if SupplierOrder != nil && SupplierOrder.Supplier != nil {
			updateForm.selectSupplier.SetSelected(SupplierOrder.Supplier.Name)
		} else {
			updateForm.selectSupplier.ClearSelected()
		}

		updateForm.gridContainerItems.RemoveAll()

		supplierOrderLineControls = make(map[*controls.SupplierOrderLineItem]int)

		//
		for _, bsl := range SupplierOrder.Lines {
			item := controls.NewSupplierOrderLineControl(bottleNames, bottleMap)
			item.Bind(bsl.Bottle.ID, bsl.SupplierOrder.ID)
			item.SelectBottle.Selected = bsl.Bottle.FullName
			item.EntryQuantity.Text = strconv.Itoa(bsl.Quantity)

			supplierOrderLineControls[item] = item.BottleId

			var deleteItemBtn *widget.Button
			deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
				func() {
					updateForm.gridContainerItems.Remove(deleteItemBtn)
					updateForm.gridContainerItems.Remove(item.SelectBottle)
					updateForm.gridContainerItems.Remove(item.EntryQuantity)
					delete(supplierOrderLineControls, item)
				})

			updateForm.gridContainerItems.Add(deleteItemBtn)
			updateForm.gridContainerItems.Add(item.SelectBottle)
			updateForm.gridContainerItems.Add(item.EntryQuantity)
		}
		updateForm.form.Refresh()
	} else {
		log(true, source, err.Error())
	}
}

func tableRefresh() {
	if table != nil && tableOptions != nil {
		getSupplierOrders()
		tableOptions.Bindings = bind
		table.Refresh()
	}
}

// endregion "table"

// endregion " events "

// region " Other functions "

func addSupplierOrderLineControl(bottleNames []string, bottleMap map[string]int, gridContainerItems *fyne.Container) {
	item := controls.NewSupplierOrderLineControl(bottleNames, bottleMap)

	supplierOrderLineControls[item] = len(supplierOrderLineControls) + 1

	var deleteItemBtn *widget.Button
	deleteItemBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(),
		func() {
			gridContainerItems.Remove(deleteItemBtn)
			gridContainerItems.Remove(item.SelectBottle)
			gridContainerItems.Remove(item.EntryQuantity)
			delete(supplierOrderLineControls, item)
		})

	gridContainerItems.Add(deleteItemBtn)
	gridContainerItems.Add(item.SelectBottle)
	gridContainerItems.Add(item.EntryQuantity)
}

// endregion " Other functions "
