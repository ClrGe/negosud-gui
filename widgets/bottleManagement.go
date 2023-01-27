package widgets

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
	"github.com/rohanthewiz/rtable"
	"negosud-gui/data"
	"strconv"
	"strings"
)

var BindBottle []binding.DataMap

// makeBottlesTabs creates a new set of tabs for bottles management
func makeBottlesTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", displayAndUpdateBottle(nil)),
		container.NewTabItem("Ajouter un produit", addNewBottle(nil)),
		container.NewTabItem("En stock", displayStock(nil)),
		container.NewTabItem("En rupture de stock", displayInventory(nil)),
		container.NewTabItem("Historique des inventaires", displayInventory(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// BottlesColumns defines the header row for the table
var BottlesColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 40},
	{ColName: "FullName", Header: "Nom", WidthPercent: 100},
	{ColName: "Label", Header: "Label", WidthPercent: 50},
	{ColName: "Year", Header: "Année", WidthPercent: 50},
}

// displayAndUpdateBottle implements a dynamic table bound to an editing form
func displayAndUpdateBottle(_ fyne.Window) fyne.CanvasObject {

	// retrieve structs from data package
	Individual := data.IndBottle
	BottleData := data.BottleData

	var identifier string
	var yPos, heightFields, widthForm float32

	yPos = 180
	heightFields = 50
	widthForm = 600

	// DETAILS PRODUCT
	instructions := widget.NewLabelWithStyle("Cliquez sur un identifiant dans le tableau pour afficher les détails", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	instructions.Resize(fyne.NewSize(widthForm, heightFields))
	instructions.Move(fyne.NewPos(0, yPos-500))

	productImg := canvas.NewImageFromFile("media/bouteille.jpeg")
	productImg.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		productImg.SetMinSize(fyne.NewSize(250, 300))
	} else {
		productImg.SetMinSize(fyne.NewSize(250, 300))
	}
	productImg.Hidden = true

	productTitle := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productTitle.Resize(fyne.NewSize(widthForm, heightFields))
	productTitle.Move(fyne.NewPos(0, yPos-400))

	productDesc := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productDesc.Resize(fyne.NewSize(widthForm, heightFields))
	productDesc.Move(fyne.NewPos(0, yPos-350))

	productLab := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productLab.Resize(fyne.NewSize(widthForm, heightFields))
	productLab.Move(fyne.NewPos(0, yPos-300))

	productVol := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productVol.Resize(fyne.NewSize(widthForm, heightFields))
	productVol.Move(fyne.NewPos(0, yPos-250))

	productYear := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productYear.Resize(fyne.NewSize(widthForm, heightFields))
	productYear.Move(fyne.NewPos(0, yPos-200))

	productPr := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productPr.Resize(fyne.NewSize(widthForm, heightFields))
	productPr.Move(fyne.NewPos(0, yPos-150))

	productAlc := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productAlc.Resize(fyne.NewSize(widthForm, heightFields))
	productAlc.Move(fyne.NewPos(0, yPos-100))

	productDetails := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "", Widget: productTitle},
			{Text: "", Widget: productDesc},
			{Text: "", Widget: productLab},
			{Text: "", Widget: productVol},
			{Text: "", Widget: productYear},
			{Text: "", Widget: productPr},
			{Text: "", Widget: productAlc},
		},
	}

	// UPDATE FORM

	// declare form elements
	nameLabel := widget.NewLabel("Nom")
	nameBottle := widget.NewEntry()
	detailsLabel := widget.NewLabel("Description")
	detailsBottle := widget.NewMultiLineEntry()
	labelLabel := widget.NewLabel("Label")
	labelBottle := widget.NewEntry()
	volumeLabel := widget.NewLabel("Volume (cL)")
	volumeBottle := widget.NewEntry()
	alcoholLabel := widget.NewLabel("Alcool (%)")
	alcoholBottle := widget.NewEntry()
	yearLabel := widget.NewLabel("Année")
	yearBottle := widget.NewEntry()
	priceLabel := widget.NewLabel("Prix")
	priceBottle := widget.NewEntry()
	pictureLabel := widget.NewLabel("Image")
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	deleteBtn := widget.NewButtonWithIcon("Supprimer ce produit", theme.WarningIcon(), func() { fmt.Print("Deleting producer") })
	deleteBtn.Resize(fyne.NewSize(600, 50))

	resultApi := data.AuthGetRequest("bottle")
	if err := json.NewDecoder(resultApi).Decode(&BottleData); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(BottleData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
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

		// binding bottle data
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
		_, err = rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		} else {
			instructions.Hidden = true
		}
		// Prevent app crash if other row than ID is clicked
		_, err = strconv.Atoi(identifier)
		if err == nil {
			resultApi := data.AuthGetRequest("bottle/" + identifier)
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
			// Display details
			productTitle.SetText("Nom: " + Individual.FullName)
			productDesc.SetText("Description: " + details)
			productLab.SetText("Label: " + Individual.Label)
			productYear.SetText("Année: " + strconv.Itoa(Individual.YearProduced))
			productVol.SetText("Volume : " + strconv.Itoa(Individual.Volume) + " cL")
			productPr.SetText("Prix HT : " + strconv.Itoa(Individual.CurrentPrice) + " €")
			productAlc.SetText("Alcool : " + strconv.Itoa(Individual.AlcoholPercentage) + " %")
		}
	}

	updateForm := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameBottle},
			{Text: "", Widget: detailsLabel},
			{Text: "", Widget: detailsBottle},
			{Text: "", Widget: labelLabel},
			{Text: "", Widget: labelBottle},
			{Text: "", Widget: volumeLabel},
			{Text: "", Widget: volumeBottle},
			{Text: "", Widget: alcoholLabel},
			{Text: "", Widget: alcoholBottle},
			{Text: "", Widget: yearLabel},
			{Text: "", Widget: yearBottle},
			{Text: "", Widget: priceLabel},
			{Text: "", Widget: priceBottle},
			{Text: "", Widget: pictureLabel},
			{Text: "", Widget: pictureBottle},
		},
		OnSubmit: func() {
			vol, _ := strconv.ParseInt(volumeBottle.Text, 10, 0)
			alc, _ := strconv.ParseInt(alcoholBottle.Text, 10, 0)
			year, _ := strconv.ParseInt(yearBottle.Text, 10, 0)
			pr, _ := strconv.ParseInt(priceBottle.Text, 10, 0)

			bottle := &data.Bottle{
				FullName:          nameBottle.Text,
				Description:       detailsBottle.Text,
				Label:             labelBottle.Text,
				Volume:            int(vol),
				AlcoholPercentage: int(alc),
				YearProduced:      int(year),
				CurrentPrice:      int(pr),
			}
			// Convert to JSON
			jsonValue, err := json.Marshal(bottle)
			if err != nil {
				fmt.Println(err)
			}
			// Send data to API
			postData := data.AuthPostRequest("bottle/"+identifier, jsonValue)
			if postData != 201|200 {
				fmt.Println("Error while sending data to API")
				fmt.Println("StatusCode " + strconv.Itoa(postData))
				return
			}
			fmt.Println("Bottle updated")
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Envoyer",
		CancelText: "Annuler",
	}

	// LAYOUT

	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 300), productDetails))

	layoutDetailsTab := container.NewBorder(image, nil, nil, nil, textProduct, instructions)
	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 750), updateForm))
	layoutWithDelete := container.NewBorder(layoutUpdateForm, deleteBtn, nil, nil)

	individualTabs := container.NewAppTabs(
		container.NewTabItem("Détails du produit", layoutDetailsTab),
		container.NewTabItem("Modifier le produit", layoutWithDelete),
	)

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)
	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new bottle to the API endpoint (POST /api/producer)
func addNewBottle(_ fyne.Window) fyne.CanvasObject {

	nameLabel := widget.NewLabel("Nom du produit")
	nameBottle := widget.NewEntry()
	descriptionLabel := widget.NewLabel("Description")
	descriptionBottle := widget.NewMultiLineEntry()
	labelLabel := widget.NewLabel("Label")
	labelBottle := widget.NewEntry()
	yearLabel := widget.NewLabel("Année")
	yearBottle := widget.NewEntry()
	volumeLabel := widget.NewLabel("Volume (cL)")
	volumeBottle := widget.NewEntry()
	alcoolLabel := widget.NewLabel("Alcool (%)")
	alcoholBottle := widget.NewEntry()
	currentPriceLabel := widget.NewLabel("Prix HT (€)")
	currentPriceBottle := widget.NewEntry()
	pictureLabel := widget.NewLabel("Image du produit")
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	form :=
		&widget.Form{
			Items: []*widget.FormItem{
				{Text: "", Widget: nameLabel},
				{Text: "", Widget: nameBottle},
				{Text: "", Widget: descriptionLabel},
				{Text: "", Widget: descriptionBottle},
				{Text: "", Widget: labelLabel},
				{Text: "", Widget: labelBottle},
				{Text: "", Widget: yearLabel},
				{Text: "", Widget: yearBottle},
				{Text: "", Widget: volumeLabel},
				{Text: "", Widget: volumeBottle},
				{Text: "", Widget: alcoolLabel},
				{Text: "", Widget: alcoholBottle},
				{Text: "", Widget: currentPriceLabel},
				{Text: "", Widget: currentPriceBottle},
				{Text: "", Widget: pictureLabel},
				{Text: "", Widget: pictureBottle},
			},
			OnSubmit: func() {
				// Convert strings to ints to match Bottle struct types
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
					FullName:          nameBottle.Text,
					Label:             labelBottle.Text,
					Volume:            volume,
					YearProduced:      year,
					AlcoholPercentage: alcohol,
					CurrentPrice:      price,
					Description:       descriptionBottle.Text,
				}
				// encode the value as JSON and send it to the API.
				jsonValue, err := json.Marshal(bottle)
				if err != nil {
					fmt.Println(err)
					return
				}
				postData := data.AuthPostRequest("bottle", jsonValue)
				if postData != 201|200 {
					fmt.Println("Error while sending data to API")
					fmt.Println("StatusCode " + strconv.Itoa(postData))
					return
				}
				fmt.Println("New product added with success")
			},
			SubmitText: "Envoyer",
		}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 800), form))

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
