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
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
	"github.com/rohanthewiz/rtable"
	"image/color"
	"negosud-gui/config"
	"net/http"
	"strconv"
	"strings"
)

var BindBottle []binding.DataMap

// makeBottleTabs creates a new set of tabs for bottle management
func makeBottleTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", displayAndUpdateBottle(nil)),
		container.NewTabItem("Ajouter un produit", addNewBottle(nil)),
		container.NewTabItem("Produits en stock", displayStock(nil)),
		container.NewTabItem("Historique des inventaires", displayInventory(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// BottlesColumns defines the header row for the table
var BottlesColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 15},
	{ColName: "FullName", Header: "Nom", WidthPercent: 120},
	{ColName: "Description", Header: "Description", WidthPercent: 50},
	{ColName: "Label", Header: "Label", WidthPercent: 10},
	{ColName: "Volume", Header: "Volume (cL)", WidthPercent: 35},
	{ColName: "Alcohol", Header: "Alcool(%", WidthPercent: 35},
	{ColName: "Year", Header: "Année", WidthPercent: 35},
	{ColName: "Price", Header: "Prix (€)", WidthPercent: 50},
}

// displayAndUpdateBottle implements a dynamic table bound to an editing form
func displayAndUpdateBottle(_ fyne.Window) fyne.CanvasObject {

	// retrieve structs from config package
	Individual := config.IndBottle
	BottleData := config.BottleData

	formTitle := canvas.NewText("Modifier un produit", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(600, 20))
	formTitle.Move(fyne.NewPos(50, 40))

	nameLabel := canvas.NewText("Nom", color.Black)
	nameLabel.Resize(fyne.NewSize(600, 20))
	nameLabel.Move(fyne.NewPos(50, 100))
	nameBottle := widget.NewEntry()
	nameBottle.Resize(fyne.NewSize(600, 40))
	nameBottle.Move(fyne.NewPos(50, 120))

	detailsLabel := canvas.NewText("Description", color.Black)
	detailsLabel.Resize(fyne.NewSize(600, 20))
	detailsLabel.Move(fyne.NewPos(50, 180))
	detailsBottle := widget.NewMultiLineEntry()
	detailsBottle.Resize(fyne.NewSize(600, 100))
	detailsBottle.Move(fyne.NewPos(50, 200))

	labelLabel := canvas.NewText("Label", color.Black)
	labelLabel.Resize(fyne.NewSize(600, 20))
	labelLabel.Move(fyne.NewPos(50, 320))
	labelBottle := widget.NewEntry()
	labelBottle.Resize(fyne.NewSize(600, 40))
	labelBottle.Move(fyne.NewPos(50, 340))

	volumeLabel := canvas.NewText("Volume (cL)", color.Black)
	volumeLabel.Resize(fyne.NewSize(600, 20))
	volumeLabel.Move(fyne.NewPos(50, 400))
	volumeBottle := widget.NewEntry()
	volumeBottle.Resize(fyne.NewSize(600, 40))
	volumeBottle.Move(fyne.NewPos(50, 420))

	alcoholLabel := canvas.NewText("Alcool (%)", color.Black)
	alcoholLabel.Resize(fyne.NewSize(600, 20))
	alcoholLabel.Move(fyne.NewPos(50, 480))
	alcoholBottle := widget.NewEntry()
	alcoholBottle.Resize(fyne.NewSize(600, 40))
	alcoholBottle.Move(fyne.NewPos(50, 500))

	yearLabel := canvas.NewText("Année", color.Black)
	yearLabel.Resize(fyne.NewSize(600, 20))
	yearLabel.Move(fyne.NewPos(50, 560))
	yearBottle := widget.NewEntry()
	yearBottle.Resize(fyne.NewSize(600, 40))
	yearBottle.Move(fyne.NewPos(50, 580))

	priceLabel := canvas.NewText("Prix", color.Black)
	priceLabel.Resize(fyne.NewSize(600, 20))
	priceLabel.Move(fyne.NewPos(50, 640))
	priceBottle := widget.NewEntry()
	priceBottle.Resize(fyne.NewSize(600, 40))
	priceBottle.Move(fyne.NewPos(50, 660))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(600, 40))
	submitBtn.Move(fyne.NewPos(50, 720))

	apiUrl := config.BottleAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&BottleData); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(BottleData); i++ {
		id := strconv.Itoa(BottleData[i].Id)
		v := strconv.Itoa(BottleData[i].VolumeInt)
		a := strconv.Itoa(BottleData[i].AlcoholPercentage)
		p := strconv.Itoa(BottleData[i].CurrentPrice)
		y := strconv.Itoa(BottleData[i].YearProduced)

		BottleData[i].Price = p
		BottleData[i].Year = y
		BottleData[i].Volume = v
		BottleData[i].Alcohol = a
		BottleData[i].ID = id
		BindBottle = append(BindBottle, binding.BindStruct(&BottleData[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: BottlesColumns,
		Bindings: BindBottle,
	}

	table := rtable.CreateTable(tableOptions)

	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 0 || cell.Row > len(BindBottle) { // 1st col is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(BottlesColumns) {
			fmt.Println("*-> Column out of limits")
			return
		}
		// Handle header row clicked
		if cell.Row == 0 { // fmt.Println("-->", tblOpts.ColAttrs[cell.Col].Header)
			// Add a row
			BindProducer = append(BindProducer,
				binding.BindStruct(&config.Producer{Name: "Belle Ambiance",
					Details: "brown", CreatedBy: "170"}))
			tableOptions.Bindings = BindProducer
			table.Refresh()
			return
		}
		//Handle non-header row clicked
		identifier, err := rtable.GetStrCellValue(cell, tableOptions)

		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}
		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]

		cellBinding, err := rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}
		fmt.Println(cellBinding)

		fmt.Println("-->", identifier)

		// Fetch individual producer to fill form
		resultApi := config.FetchIndividualBottle(identifier)
		if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		nameBottle.SetText(Individual.FullName)
		details := strings.Replace(Individual.Description, "\\n", "\n", -1)
		detailsBottle.SetText(details)
		labelBottle.SetText(Individual.Label)
		volumeBottle.SetText(strconv.Itoa(Individual.Volume))
		yearBottle.SetText(strconv.Itoa(Individual.YearProduced))
		priceBottle.SetText(strconv.Itoa(Individual.CurrentPrice))
		alcoholBottle.SetText(strconv.Itoa(Individual.AlcoholPercentage))

	}

	// Define layout
	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewWithoutLayout(formTitle, nameLabel, nameBottle, detailsLabel, detailsBottle, labelLabel, labelBottle, volumeLabel, volumeBottle, alcoholLabel, alcoholBottle, yearLabel, yearBottle, priceLabel, priceBottle, submitBtn)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new bottle to the API endpoint (POST /api/producer)
func addNewBottle(w fyne.Window) fyne.CanvasObject {

	apiUrl := config.BottleAPIConfig()

	idBottle := widget.NewEntry()
	nameBottle := widget.NewEntry()
	descriptionBottle := widget.NewEntry()
	labelBottle := widget.NewEntry()
	yearBottle := widget.NewEntry()
	volumeBottle := widget.NewEntry()
	alcoholBottle := widget.NewEntry()
	currentPriceBottle := widget.NewEntry()
	createdByBottle := widget.NewEntry()

	title := widget.NewLabelWithStyle("AJOUTER UN NOUVEAU PRODUIT", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form :=
		&widget.Form{
			Items: []*widget.FormItem{
				{Text: "ID", Widget: idBottle},
				{Text: "Nom", Widget: nameBottle},
				{Text: "Description", Widget: descriptionBottle},
				{Text: "Label", Widget: labelBottle},
				{Text: "Year", Widget: yearBottle},
				{Text: "Volume", Widget: volumeBottle},
				{Text: "Alcohol percentage", Widget: alcoholBottle},
				{Text: "Current price", Widget: currentPriceBottle},
				{Text: "Created By", Widget: createdByBottle},
			},
			OnCancel: func() {
				fmt.Println("Annulé")
			},
			OnSubmit: func() {

				// Convert strings to ints to match Bottle struct types
				id, err := strconv.Atoi(idBottle.Text)
				if err != nil {
					return
				}
				volume, err := strconv.Atoi(volumeBottle.Text)
				if err != nil {
					return
				}
				year, err := strconv.Atoi(yearBottle.Text)
				if err != nil {
					return
				}
				alcohol, err := strconv.Atoi(alcoholBottle.Text)
				if err != nil {
					return
				}
				price, err := strconv.Atoi(currentPriceBottle.Text)
				if err != nil {
					return
				}

				// extract the value from the input widget and set the corresponding field in the Producer struct
				bottle := &config.Bottle{
					ID:                id,
					FullName:          nameBottle.Text,
					Label:             labelBottle.Text,
					Volume:            volume,
					YearProduced:      year,
					AlcoholPercentage: alcohol,
					CurrentPrice:      price,
					CreatedBy:         createdByBottle.Text,
					Description:       descriptionBottle.Text,
				}

				// encode the value as JSON and send it to the API.
				bottleJsonValue, _ := json.Marshal(bottle)
				bottleResp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(bottleJsonValue))
				if err != nil {
					fmt.Println("error while encoding response")
					return
				}
				if bottleResp.StatusCode == 204 {
					config.BottleFailureDialog(w)
					fmt.Println(bottleJsonValue)
					return
				}
				config.BottleSuccessDialog(w)
				fmt.Println("New bottle added with success")

			},
		}
	form.Append("Description", descriptionBottle)
	form.Size()
	mainContainer := container.NewVBox(title, form)

	return mainContainer
}

func displayStock(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Stock disponible (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}

func displayInventory(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Historique des inventaires entrepôt (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
