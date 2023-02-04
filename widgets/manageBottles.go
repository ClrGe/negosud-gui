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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
	"negosud-gui/data"
	"os"
	"strconv"
	"strings"
	"time"
)

var BindBottle []binding.DataMap

var BottleTableRefreshMethod func()

// BottlesColumns defines the header row for the table
var BottlesColumns = []rtable.ColAttr{

	{ColName: "ID", Header: "ID", WidthPercent: 40},
	{ColName: "FullName", Header: "Nom", WidthPercent: 80},
	{ColName: "WineType", Header: "Type", WidthPercent: 30},
	{ColName: "Year", Header: "Année", WidthPercent: 50},
	{ColName: "quantity", Header: "Quantité", WidthPercent: 50},
}

// makeBottlesPage creates a new set of tabs for bottles management
func makeBottlesPage(_ fyne.Window) fyne.CanvasObject {
	bottleListTab := container.NewTabItem("Liste des produits", displayAndUpdateBottle(nil))
	tabs := container.NewAppTabs(
		bottleListTab,
		container.NewTabItem("Ajouter un produit", addNewBottle(nil)),
		container.NewTabItem("En stock", displayStock(nil)),
		container.NewTabItem("En rupture de stock", displayInventory(nil)),
		container.NewTabItem("Historique des inventaires", displayInventory(nil)),
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		if ti == bottleListTab {
			BottleTableRefreshMethod()
		}
	}
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func getBottles() bool {
	var source = "WIDGETS.BOTTLE "
	// retrieve structs from data package

	BottleData := data.BottleData

	response := data.AuthGetRequest("bottle")
	if response == nil {
		message := "Request body returned empty"
		fmt.Println(message)
		data.Logger(false, "WIDGETS.BOTTLE", message)
		return false
	}
	if err := json.NewDecoder(response).Decode(&BottleData); err != nil {
		log(true, source, err.Error())
		fmt.Println(err)
	}

	BindBottle = nil

	for i := 0; i < len(BottleData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		id := strconv.Itoa(BottleData[i].Id)
		v := strconv.Itoa(BottleData[i].VolumeInt)
		alcoholPercentage := fmt.Sprintf("%f", BottleData[i].AlcoholPercentage)
		p := strconv.Itoa(BottleData[i].CurrentPrice)
		y := strconv.Itoa(BottleData[i].YearProduced)
		BottleData[i].Price = p
		BottleData[i].Year = y
		BottleData[i].Volume = v
		BottleData[i].Alcohol = alcoholPercentage
		BottleData[i].ID = id

		// binding bottle data
		BindBottle = append(BindBottle, binding.BindStruct(&BottleData[i]))
	}
	return true
}

// displayAndUpdateBottle implements a dynamic table bound to an editing form
func displayAndUpdateBottle(_ fyne.Window) fyne.CanvasObject {
	var source = "WIDGETS.BOTTLE "
	Bottle := data.IndBottle

	if !getBottles() {
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}

	var identifier string
	var yPos, heightFields, widthForm float32

	yPos = 180
	heightFields = 35
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
	nameBottle := widget.NewEntry()
	detailsBottle := widget.NewMultiLineEntry()
	typeBottle := widget.NewSelectEntry([]string{"Rouge", "Blanc", "Rosé", "Dessert", "Pétillant"})
	volumeBottle := widget.NewEntry()
	alcoholBottle := widget.NewEntry()
	yearBottle := widget.NewEntry()
	priceBottle := widget.NewEntry()
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	deleteBtn := widget.NewButtonWithIcon("Supprimer ce produit", theme.WarningIcon(), func() { fmt.Print("Deleting producer") })
	deleteBtn.Resize(fyne.NewSize(600, 50))

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: BottlesColumns,
		Bindings: BindBottle,
	}

	table := rtable.CreateTable(tableOptions)
	BottleTableRefreshMethod = func() {
		getBottles()
		tableOptions.Bindings = BindBottle
		table.Refresh()
	}
	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 1 || cell.Row > len(BindBottle) { // 1st col is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(BottlesColumns) {
			fmt.Println("*-> Column out of limits")
			return
		}

		//Handle non-header row clicked
		identifier, err := rtable.GetStrCellValue(cell, tableOptions)
		if err != nil {
			fmt.Println(err.Error())
			log(true, source, err.Error())

			return
		}
		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]
		_, err = rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(err.Error())
			log(true, source, err.Error())
			return
		} else {
			instructions.Hidden = true
		}
		// Prevent app crash if other row than ID is clicked
		_, err = strconv.Atoi(identifier)
		if err == nil {
			resultApi := data.AuthGetRequest("bottle/" + identifier)
			if err := json.NewDecoder(resultApi).Decode(&Bottle); err != nil {
				fmt.Println(err)
				log(true, source, err.Error())
			} else {
				productImg.Hidden = false
			}
			// Fill form fields with fetched data
			nameBottle.SetText(Bottle.FullName)
			details := strings.Replace(Bottle.Description, "\\n", "\n", -1)
			detailsBottle.SetText(details)
			typeBottle.SetPlaceHolder(Bottle.WineType)
			volumeBottle.SetText(strconv.Itoa(Bottle.Volume))
			yearBottle.SetText(strconv.Itoa(Bottle.YearProduced))
			priceBottle.SetText(fmt.Sprintf("%f", Bottle.CurrentPrice))
			alcoholBottle.SetText(fmt.Sprintf("%f", Bottle.AlcoholPercentage))
			// Display details
			productTitle.SetText("Nom: " + Bottle.FullName)
			productDesc.SetText("Description: " + details)
			productLab.SetText("Type : " + Bottle.WineType)
			productYear.SetText("Année : " + strconv.Itoa(Bottle.YearProduced))
			productVol.SetText("Volume : " + strconv.Itoa(Bottle.Volume) + " cL")
			productPr.SetText("Prix HT : " + fmt.Sprintf("%f", Bottle.CurrentPrice) + " €")
			productAlc.SetText("Alcool : " + fmt.Sprintf("%f", Bottle.AlcoholPercentage) + " %")
		}
	}

	updateForm := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "Nom", Widget: nameBottle},
			{Text: "Description", Widget: detailsBottle},
			{Text: "Type", Widget: typeBottle},
			{Text: "Vol. (cL)", Widget: volumeBottle},
			{Text: "Alc. (%)", Widget: alcoholBottle},
			{Text: "Année", Widget: yearBottle},
			{Text: "Prix (€)", Widget: priceBottle},
			{Text: "", Widget: pictureBottle},
		},
		OnSubmit: func() {
			vol, _ := strconv.ParseInt(volumeBottle.Text, 10, 0)
			alc, _ := strconv.ParseInt(alcoholBottle.Text, 10, 0)
			year, _ := strconv.ParseInt(yearBottle.Text, 10, 0)
			pr, _ := strconv.ParseInt(priceBottle.Text, 10, 0)
			who, _ := os.Hostname()
			t, _ := time.Parse("2023-01-27T22:48:02.646Z", time.Now().String())
			bottle := &data.Bottle{
				ID:                Bottle.ID,
				FullName:          nameBottle.Text,
				Description:       detailsBottle.Text,
				WineType:          typeBottle.Text,
				Volume:            int(vol),
				AlcoholPercentage: float32(alc),
				CreatedAt:         t,
				UpdatedAt:         t,
				YearProduced:      int(year),
				CreatedBy:         who,
				UpdatedBy:         who,
				CurrentPrice:      float32(pr),
			}
			// Convert to JSON
			jsonValue, err := json.Marshal(bottle)
			if err != nil {
				log(true, source, err.Error())
				fmt.Println(err)
			}
			fmt.Print(bytes.NewBuffer(jsonValue))
			// Send data to API
			postData := data.AuthPostRequest("Bottle/UpdateBottle/"+identifier, bytes.NewBuffer(jsonValue))
			if postData != 200 {
				message := "Error on bottle " + identifier + " update " + " StatusCode " + strconv.Itoa(postData)
				log(true, source, message)
			}
			fmt.Println("Bottle updated")
			BottleTableRefreshMethod()
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Envoyer",
		CancelText: "Annuler",
	}

	// LAYOUT

	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 100), productDetails))

	layoutDetailsTab := container.NewBorder(image, nil, nil, nil, textProduct, instructions)
	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
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
	var source = "WIDGETS.BOTTLE "

	nameBottle := widget.NewEntry()
	descriptionBottle := widget.NewMultiLineEntry()
	typeBottle := widget.NewSelectEntry([]string{"Red", "White", "Rosé", "Dessert", "Sparkling"})
	typeBottle.SetPlaceHolder("Veuillez sélectionner un type de vin...")
	yearBottle := widget.NewEntry()
	volumeBottle := widget.NewEntry()
	alcoholBottle := widget.NewEntry()
	currentPriceBottle := widget.NewEntry()
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	form :=
		&widget.Form{
			Items: []*widget.FormItem{
				{Text: "Nom", Widget: nameBottle},
				{Text: "Description", Widget: descriptionBottle},
				{Text: "Type", Widget: typeBottle},
				{Text: "Année", Widget: yearBottle},
				{Text: "Vol. (cL)", Widget: volumeBottle},
				{Text: "Alc. (%)", Widget: alcoholBottle},
				{Text: "Prix (€)", Widget: currentPriceBottle},
				{Text: "", Widget: pictureBottle},
			},
			OnSubmit: func() {
				// Convert strings to ints to match Bottle struct types
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
				volume, err := strconv.Atoi(volumeBottle.Text)
				if err != nil {
					return
				}
				// extract the value from the input widget and set the corresponding field in the Producer struct
				bottle := &data.Bottle{
					FullName:          nameBottle.Text,
					WineType:          typeBottle.Text,
					YearProduced:      year,
					AlcoholPercentage: float32(alcohol),
					CurrentPrice:      float32(price),
					Volume:            volume,
					Description:       descriptionBottle.Text,
				}
				// encode the value as JSON and send it to the API.
				jsonValue, err := json.Marshal(bottle)
				if err != nil {
					log(true, source, err.Error())
					fmt.Println(err)
					return
				}
				postData := data.AuthPostRequest("Bottle/AddBottle", bytes.NewBuffer(jsonValue))
				if postData != 201|200 {
					fmt.Println("Error while sending data to API")
					message := "Error while creating new Bottle. StatusCode " + strconv.Itoa(postData)
					log(true, source, message)
					return
				}
				fmt.Println("New product added with success")
				BottleTableRefreshMethod()
			},
			SubmitText: "Envoyer",
		}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))

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
