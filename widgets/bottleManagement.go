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
	"negosud-gui/data"
	"net/http"
	"strconv"
	"strings"
)

var BindBottle []binding.DataMap

// makeBottlesTabs creates a new set of tabs for bottles management
func makeBottlesTabs(_ fyne.Window) fyne.CanvasObject {
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
func displayAndUpdateBottle(w fyne.Window) fyne.CanvasObject {

	// retrieve structs from data package
	Individual := data.IndBottle
	BottleData := data.BottleData

	var xPos, yPos, heightFields, heightLabels, widthForm float32

	xPos = -200
	yPos = 180
	heightFields = 50
	heightLabels = 20
	widthForm = 600

	// DETAILS PRODUCT

	instructions := canvas.NewText("Cliquez sur un identifiant dans le tableau pour afficher les détails du produit.", color.Black)
	instructions.TextSize = 15
	instructions.TextStyle = fyne.TextStyle{Bold: true}
	instructions.Resize(fyne.NewSize(widthForm, heightFields))
	instructions.Move(fyne.NewPos(0, yPos-500))

	productImg := canvas.NewImageFromFile("media/bouteille.jpeg")
	productImg.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		productImg.SetMinSize(fyne.NewSize(250, 350))
	} else {
		productImg.SetMinSize(fyne.NewSize(250, 350))
	}
	productImg.Hidden = true

	productTitle := widget.NewLabel("")
	productTitle.Resize(fyne.NewSize(widthForm, heightFields))
	productTitle.Move(fyne.NewPos(0, yPos-400))

	productDesc := widget.NewLabel("")
	productDesc.Resize(fyne.NewSize(widthForm, heightFields))
	productDesc.Move(fyne.NewPos(0, yPos-350))

	productLab := widget.NewLabel("")
	productLab.Resize(fyne.NewSize(widthForm, heightFields))
	productLab.Move(fyne.NewPos(0, yPos-300))

	productVol := widget.NewLabel("")
	productVol.Resize(fyne.NewSize(widthForm, heightFields))
	productVol.Move(fyne.NewPos(0, yPos-250))

	productYear := widget.NewLabel("")
	productYear.Resize(fyne.NewSize(widthForm, heightFields))
	productYear.Move(fyne.NewPos(0, yPos-200))

	productPr := widget.NewLabel("")
	productPr.Resize(fyne.NewSize(widthForm, heightFields))
	productPr.Move(fyne.NewPos(0, yPos-150))

	productAlc := widget.NewLabel("")
	productAlc.Resize(fyne.NewSize(widthForm, heightFields))
	productAlc.Move(fyne.NewPos(0, yPos-100))

	// UPDATE FORM

	formTitle := canvas.NewText("Modifier un produit", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(widthForm, heightFields))
	formTitle.Move(fyne.NewPos(xPos, yPos-600))

	nameLabel := canvas.NewText("Nom", color.Black)
	nameLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	nameLabel.Move(fyne.NewPos(xPos, yPos-520))
	nameBottle := widget.NewEntry()
	nameBottle.Resize(fyne.NewSize(widthForm, heightFields))
	nameBottle.Move(fyne.NewPos(xPos, yPos-500))

	detailsLabel := canvas.NewText("Description", color.Black)
	detailsLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	detailsLabel.Move(fyne.NewPos(xPos, yPos-420))
	detailsBottle := widget.NewMultiLineEntry()
	detailsBottle.Resize(fyne.NewSize(widthForm, heightFields*2))
	detailsBottle.Move(fyne.NewPos(xPos, yPos-400))

	labelLabel := canvas.NewText("Label", color.Black)
	labelLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	labelLabel.Move(fyne.NewPos(xPos, yPos-280))
	labelBottle := widget.NewEntry()
	labelBottle.Resize(fyne.NewSize(widthForm, heightFields))
	labelBottle.Move(fyne.NewPos(xPos, yPos-260))

	volumeLabel := canvas.NewText("Volume (cL)", color.Black)
	volumeLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	volumeLabel.Move(fyne.NewPos(xPos, yPos-180))
	volumeBottle := widget.NewEntry()
	volumeBottle.Resize(fyne.NewSize(widthForm, heightFields))
	volumeBottle.Move(fyne.NewPos(xPos, yPos-160))

	alcoholLabel := canvas.NewText("Alcool (%)", color.Black)
	alcoholLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	alcoholLabel.Move(fyne.NewPos(xPos, yPos-80))
	alcoholBottle := widget.NewEntry()
	alcoholBottle.Resize(fyne.NewSize(widthForm, heightFields))
	alcoholBottle.Move(fyne.NewPos(xPos, yPos-60))

	yearLabel := canvas.NewText("Année", color.Black)
	yearLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	yearLabel.Move(fyne.NewPos(xPos, yPos+20))
	yearBottle := widget.NewEntry()
	yearBottle.Resize(fyne.NewSize(widthForm, heightFields))
	yearBottle.Move(fyne.NewPos(xPos, yPos+40))

	priceLabel := canvas.NewText("Prix", color.Black)
	priceLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	priceLabel.Move(fyne.NewPos(xPos, yPos+120))
	priceBottle := widget.NewEntry()
	priceBottle.Resize(fyne.NewSize(widthForm, heightFields))
	priceBottle.Move(fyne.NewPos(xPos, yPos+140))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(widthForm, heightFields))
	submitBtn.Move(fyne.NewPos(xPos, yPos+220))

	apiUrl := data.BottleAPIConfig()

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
				binding.BindStruct(&data.Producer{Name: "Belle Ambiance",
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
		resultApi := data.FetchIndividualBottle(identifier)
		if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
			fmt.Println(err)
		} else {
			productImg.Hidden = false
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

		productTitle.SetText("Nom: " + Individual.FullName)
		productDesc.SetText("Description: " + details)
		productLab.SetText("Label: " + Individual.FullName)
		productYear.SetText("Année: " + strconv.Itoa(Individual.YearProduced))
		productVol.SetText("Volume (cL): " + strconv.Itoa(Individual.Volume))
		productPr.SetText("Prix HT: " + strconv.Itoa(Individual.CurrentPrice))
		productAlc.SetText("Alcool (%): " + strconv.Itoa(Individual.AlcoholPercentage))
	}

	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewWithoutLayout(productTitle, productDesc, productLab, productVol, productYear, productPr, productAlc))
	detailsProduct := container.NewBorder(image, instructions, nil, nil, textProduct)

	updateForm := container.NewCenter(container.NewWithoutLayout(formTitle, nameLabel, nameBottle, detailsLabel, detailsBottle, labelLabel, labelBottle, volumeLabel, volumeBottle, alcoholLabel, alcoholBottle, yearLabel, yearBottle, priceLabel, priceBottle, submitBtn))

	// Define layout

	individualTabs := container.NewAppTabs(
		container.NewTabItem("Détails du produit", detailsProduct),
		container.NewTabItem("Modifier le produit", updateForm),
	)

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// makeBottleTabs creates a new set of tabs for individual bottle management
func makeBottleTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Détails du produit", displayAndUpdateBottle(nil)),
		container.NewTabItem("Ajouter un produit", addNewBottle(nil)),
		container.NewTabItem("Produits en stock", displayStock(nil)),
		container.NewTabItem("Historique des inventaires", displayInventory(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// Form to add and send a new bottle to the API endpoint (POST /api/producer)
func addNewBottle(w fyne.Window) fyne.CanvasObject {

	apiUrl := data.BottleAPIConfig()

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
				bottle := &data.Bottle{
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
					data.BottleFailureDialog(w)
					fmt.Println(bottleJsonValue)
					return
				}
				data.BottleSuccessDialog(w)
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
