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

var bindingbottle []binding.DataMap

// makeBottles creates a new set of tabs for bottles management
func makeBottles(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", displayAndUpdateNewBottle(nil)),
		container.NewTabItem("Ajouter un produit", addNewTestBottle(nil)),
		container.NewTabItem("En stock", displayTestStock(nil)),
		container.NewTabItem("En rupture de stock", displayTestInventory(nil)),
		container.NewTabItem("Historique des inventaires", displayTestInventory(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// NewBottlesColumns  defines the header row for the table
var NewBottlesColumns = []rtable.ColAttr{

	{ColName: "ID", Header: "ID", WidthPercent: 40},
	{ColName: "FullName", Header: "Nom", WidthPercent: 80},
	{ColName: "WineType", Header: "Type", WidthPercent: 30},
	{ColName: "Year", Header: "Année", WidthPercent: 50},
	{ColName: "quantity", Header: "Quantité", WidthPercent: 50},
	{ColName: "Location", Header: "Emplacement", WidthPercent: 50},
}

// displayAndUpdateNewBottle implements a dynamic table bound to an editing form
func displayAndUpdateNewBottle(_ fyne.Window) fyne.CanvasObject {
	var source = "WIDGETS.BOTTLE "

	// retrieve structs from data package
	Individual := data.IndBottle
	BottleData := data.BottleData
	storageLocationForm := data.IndStorageLocation

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

	location := getAllLocationName()
	locationName := make([]string, 0)
	for i := 0; i < len(location); i++ {
		location[i].ID = strconv.Itoa(location[i].Id)
		name := location[i].Name
		locationName = append(locationName, name)
	}

	// UPDATE FORM
	// declare form elements
	nameBottle := widget.NewEntry()
	wineTypeBottle := widget.NewSelectEntry([]string{"Rouge", "Blanc", "Rosé", "Digestif", "Pétillant"})
	storageLocationData := widget.NewSelectEntry(locationName)
	quantityBottle := widget.NewEntry()
	volumeBottle := widget.NewEntry()
	alcoholBottle := widget.NewEntry()
	yearBottle := widget.NewEntry()
	detailsBottle := widget.NewMultiLineEntry()
	priceBottle := widget.NewEntry()
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	deleteBtn := widget.NewButtonWithIcon("Supprimer ce produit", theme.WarningIcon(), func() { fmt.Print("Deleting producer") })
	deleteBtn.Resize(fyne.NewSize(600, 50))

	response := data.AuthGetRequest("bottle")
	if response == nil {
		message := "Request body returned empty"
		fmt.Println(message)
		data.Logger(false, source, message)
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&BottleData); err != nil {
		log(true, source, err.Error())
		fmt.Println(err)
	}

	for i := 0; i < len(BottleData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		id := strconv.Itoa(BottleData[i].Id)
		volume := strconv.Itoa(BottleData[i].VolumeInt)
		alcoholPercentage := fmt.Sprintf("%f", BottleData[i].AlcoholPercentage)
		price := strconv.Itoa(BottleData[i].CurrentPrice)
		year := strconv.Itoa(BottleData[i].YearProduced)
		BottleData[i].Price = price
		BottleData[i].Year = year
		BottleData[i].Volume = volume
		BottleData[i].Alcohol = alcoholPercentage
		BottleData[i].ID = id

		// binding bottle data
		bindingbottle = append(bindingbottle, binding.BindStruct(&BottleData[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: BottlesColumns,
		Bindings: bindingbottle,
	}

	table := rtable.CreateTable(tableOptions)
	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 1 || cell.Row > len(bindingbottle) { // 1st col is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(BottlesColumns) {
			fmt.Println("*-> Column out of limits")
			return
		}

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
			if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
				fmt.Println(err)
				log(true, source, err.Error())
			} else {
				productImg.Hidden = false
			}
			// Fill form fields with fetched data
			nameBottle.SetText(Individual.FullName)
			details := strings.Replace(Individual.Description, "\\n", "\n", -1)
			detailsBottle.SetText(details)
			wineTypeBottle.SetPlaceHolder(Individual.WineType)
			storageLocationData.SetPlaceHolder(storageLocationForm.Name)

			volumeBottle.SetText(strconv.Itoa(Individual.Volume))
			yearBottle.SetText(strconv.Itoa(Individual.YearProduced))
			priceBottle.SetText(fmt.Sprintf("%f", Individual.CurrentPrice))
			alcoholBottle.SetText(fmt.Sprintf("%f", Individual.AlcoholPercentage))
		}
	}

	updateForm := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "Nom", Widget: nameBottle},
			{Text: "Emplacement", Widget: storageLocationData},
			{Text: "Quantité", Widget: quantityBottle},
			{Text: "Description", Widget: detailsBottle},
			{Text: "Type", Widget: wineTypeBottle},
			{Text: "Vol. (cL)", Widget: volumeBottle},
			{Text: "Alc. (%)", Widget: alcoholBottle},
			{Text: "Année", Widget: yearBottle},
			{Text: "Prix (€)", Widget: priceBottle},
			{Text: "", Widget: pictureBottle},
		},
		OnSubmit: func() {
			volume, _ := strconv.ParseInt(volumeBottle.Text, 10, 0)
			alcohol, _ := strconv.ParseFloat(alcoholBottle.Text, 32)
			year, _ := strconv.ParseInt(yearBottle.Text, 10, 0)
			price, _ := strconv.ParseFloat(priceBottle.Text, 32)
			who, _ := os.Hostname()
			timeOfCreationOrUpdate, _ := time.Parse("2023-01-27T22:48:02.646Z", time.Now().String())

			quantity, _ := strconv.ParseInt(quantityBottle.Text, 10, 0)

			bottleStorageLocation := make([]data.BottleStorageLocation, 0)
			storageLocation := data.IndStorageLocation

			for i := 0; i < 1; i++ {
				storageLocation.Name = storageLocationData.Text
				storageLocation.ID =
				dataToSent := data.BottleStorageLocation{
					StorageLocation: storageLocation,
					Quantity:        int(quantity),
				}
				bottleStorageLocation = append(bottleStorageLocation, dataToSent)
			}

			bottle := &data.Bottle{
				ID:                     Individual.ID,
				FullName:               nameBottle.Text,
				Description:            detailsBottle.Text,
				WineType:               wineTypeBottle.Text,
				Volume:                 int(volume),
				AlcoholPercentage:      float32(alcohol),
				CreatedAt:              timeOfCreationOrUpdate,
				UpdatedAt:              timeOfCreationOrUpdate,
				YearProduced:           int(year),
				CreatedBy:              who,
				UpdatedBy:              who,
				CurrentPrice:           float32(price),
				BottleStorageLocations: bottleStorageLocation,
			}

			// Convert to JSON
			jsonValue, err := json.Marshal(bottle)
			if err != nil {
				log(true, source, err.Error())
				fmt.Println(err)
			}
			jsonbuffered := bytes.NewBuffer(jsonValue)
			fmt.Print(jsonbuffered)
			// Send data to API
			postData := data.AuthPostRequest("Bottle/UpdateBottle/"+identifier, bytes.NewBuffer(jsonValue))
			if postData != 200 {
				message := "Error on bottle " + identifier + " update " + " StatusCode " + strconv.Itoa(postData)
				log(true, source, message)
			} else if postData == 200 {
				fmt.Println("Bottle updated")
			}
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Mettre à jour",
		CancelText: "Annuler",
	}

	// LAYOUT
	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithDelete := container.NewBorder(layoutUpdateForm, deleteBtn, nil, nil)

	individualTabs := container.NewAppTabs(
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
func addNewTestBottle(_ fyne.Window) fyne.CanvasObject {
	var source = "WIDGETS.BOTTLE "

	nameBottle := widget.NewEntry()
	descriptionBottle := widget.NewMultiLineEntry()
	typeBottle := widget.NewEntry()
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
				// extract the value from the input widget and set the corresponding field in the Producer struct
				bottle := &data.PartialBottle{
					FullName:          nameBottle.Text,
					WineType:          typeBottle.Text,
					YearProduced:      year,
					AlcoholPercentage: float32(alcohol),
					CurrentPrice:      price,
					Description:       descriptionBottle.Text,
				}
				// encode the value as JSON and send it to the API.
				jsonValue, err := json.Marshal(bottle)
				if err != nil {
					log(true, source, err.Error())
					fmt.Println(err)
					return
				}
				postData := data.AuthPostRequest("bottle", bytes.NewBuffer(jsonValue))
				if postData != 201|200 {
					fmt.Println("Error while sending data to API")
					message := "Error while creating new Bottle. StatusCode " + strconv.Itoa(postData)
					log(true, source, message)
					return
				}
				fmt.Println("New product added with success")
			},
			SubmitText: "Envoyer",
		}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))

	return mainContainer
}

func displayTestStock(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Stock disponible (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}

func displayTestInventory(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Historique des inventaires entrepôt (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}

func getAllLocationName() []data.PartialStorageLocation {
	var source = "WIDGETS Test bottle Get Locations "
	storageLocationData := data.StorageLocationsData
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
