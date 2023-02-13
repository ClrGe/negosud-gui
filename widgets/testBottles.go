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

var bindingBottle []binding.DataMap
var source = "WIDGETS.BOTTLE "

// NewBottlesColumns  defines the header row for the table
var NewBottlesColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 40},
	{ColName: "FullName", Header: "Nom", WidthPercent: 80},
	{ColName: "WineType", Header: "Type", WidthPercent: 30},
	{ColName: "Year", Header: "Année", WidthPercent: 50},
	{ColName: "quantity", Header: "Quantité", WidthPercent: 50},
	{ColName: "Location", Header: "Emplacement", WidthPercent: 50},
}
var identifier string

// makeBottleTabs creates a new set of tabs for bottles management
func makeBottleTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", displayAndUpdateNewBottle(nil)),
		container.NewTabItem("Ajouter un produit", addNewTestBottle(nil)),
		container.NewTabItem("En stock", displayTestStock(nil)),
		container.NewTabItem("En rupture de stock", displayTestInventory(nil)),
		container.NewTabItem("Historique des inventaires", displayTestInventory(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// displayAndUpdateNewBottle implements a dynamic table bound to an editing form
func displayAndUpdateNewBottle(_ fyne.Window) fyne.CanvasObject {

	// retrieve structs from data package
	IndividualBottle := data.IndBottle
	BottleDataFomApi := data.Bottles
	storageLocationForm := data.IndStorageLocation

	// LAYOUT
	productImg := canvas.NewImageFromFile("media/bouteille.jpeg")
	productImg.FillMode = canvas.ImageFillContain

	if fyne.CurrentDevice().IsMobile() {
		productImg.SetMinSize(fyne.NewSize(250, 300))
	} else {
		productImg.SetMinSize(fyne.NewSize(250, 300))
	}
	productImg.Hidden = true

	location := getAllLocationName()
	locationName := GetLocationNameAsArray(location)

	nameBottle, wineTypeBottle, storageLocationData, quantityBottle, volumeBottle, alcoholBottle, yearBottle, detailsBottle, priceBottle, pictureBottle := CreateUpdateFrom(locationName)

	response := data.AuthGetRequest("bottle")
	if response == nil {
		message := "Request body returned empty"
		fmt.Println(message)
		data.Logger(false, source, message)
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&BottleDataFomApi); err != nil {
		log(true, source, err.Error())
		fmt.Println(err)
	}

	BottleToDisplay := make([]data.PartialBottle, len(BottleDataFomApi))

	for i := 0; i < len(BottleDataFomApi); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		BottleToDisplay[i].Price = fmt.Sprintf("%f", BottleDataFomApi[i].CurrentPrice)
		BottleToDisplay[i].FullName = BottleDataFomApi[i].FullName
		BottleToDisplay[i].Year = strconv.Itoa(BottleDataFomApi[i].YearProduced)
		BottleToDisplay[i].Volume = strconv.Itoa(BottleDataFomApi[i].Volume)
		BottleToDisplay[i].Alcohol = fmt.Sprintf("%f", BottleDataFomApi[i].AlcoholPercentage)
		BottleToDisplay[i].ID = strconv.Itoa(BottleDataFomApi[i].ID)

		// binding bottle data
		bindingBottle = append(bindingBottle, binding.BindStruct(&BottleToDisplay[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: NewBottlesColumns,
		Bindings: bindingBottle,
	}

	table := rtable.CreateTable(tableOptions)
	table.OnSelected = func(cell widget.TableCellID) {
		DisplaySelectedRow(cell, tableOptions, IndividualBottle, productImg, nameBottle, detailsBottle, wineTypeBottle, storageLocationData, storageLocationForm, volumeBottle, yearBottle, priceBottle, alcoholBottle)
	}

	updateForm := UpdateForm(nameBottle, storageLocationData, quantityBottle, detailsBottle, wineTypeBottle, volumeBottle, alcoholBottle, yearBottle, priceBottle, pictureBottle, location, IndividualBottle, identifier)

	mainContainer := SetDesignOptions(updateForm, table)

	return mainContainer
}

// DisplaySelectedRow Display selected row in the update form
func DisplaySelectedRow(cell widget.TableCellID, tableOptions *rtable.TableOptions, IndividualBottle data.Bottle, productImg *canvas.Image, nameBottle *widget.Entry, detailsBottle *widget.Entry, wineTypeBottle *widget.SelectEntry, storageLocationData *widget.SelectEntry, storageLocationForm data.StorageLocation, volumeBottle *widget.Entry, yearBottle *widget.Entry, priceBottle *widget.Entry, alcoholBottle *widget.Entry) {
	if cell.Row < 1 || cell.Row > len(bindingBottle) { // 1st col is header
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
	}

	// Prevent app crash if other row than ID is clicked
	_, err = strconv.Atoi(identifier)
	if err == nil {
		resultApi := data.AuthGetRequest("bottle/" + identifier)
		if err := json.NewDecoder(resultApi).Decode(&IndividualBottle); err != nil {
			fmt.Println(err)
			log(true, source, err.Error())
		} else {
			productImg.Hidden = false
		}

		FillUpdateForm(nameBottle, IndividualBottle, detailsBottle, wineTypeBottle, storageLocationData, storageLocationForm, volumeBottle, yearBottle, priceBottle, alcoholBottle)
	}
}

// LAYOUT
func SetDesignOptions(updateForm *widget.Form, table *widget.Table) *fyne.Container {
	deleteBtn := widget.NewButtonWithIcon("Supprimer ce produit", theme.WarningIcon(), func() { fmt.Print("Deleting producer") })
	deleteBtn.Resize(fyne.NewSize(600, 50))
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

func GetLocationNameAsArray(location map[string]int) []string {
	locationNameArray := make([]string, 0)
	for name := range location {
		locationNameArray = append(locationNameArray, name)
	}
	return locationNameArray
}

func CreateUpdateFrom(locationNameArray []string) (*widget.Entry, *widget.SelectEntry, *widget.SelectEntry, *widget.Entry, *widget.Entry, *widget.Entry, *widget.Entry, *widget.Entry, *widget.Entry, *widget.Button) {
	// UPDATE FORM
	// declare form elements
	nameBottle := widget.NewEntry()
	wineTypeBottle := widget.NewSelectEntry([]string{"Rouge", "Blanc", "Rosé", "Digestif", "Pétillant"})
	storageLocationData := widget.NewSelectEntry(locationNameArray)
	quantityBottle := widget.NewEntry()
	volumeBottle := widget.NewEntry()
	alcoholBottle := widget.NewEntry()
	yearBottle := widget.NewEntry()
	detailsBottle := widget.NewMultiLineEntry()
	priceBottle := widget.NewEntry()
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })
	return nameBottle, wineTypeBottle, storageLocationData, quantityBottle, volumeBottle, alcoholBottle, yearBottle, detailsBottle, priceBottle, pictureBottle
}

func UpdateForm(nameBottle *widget.Entry, storageLocationData *widget.SelectEntry, quantityBottle *widget.Entry, detailsBottle *widget.Entry, wineTypeBottle *widget.SelectEntry, volumeBottle *widget.Entry, alcoholBottle *widget.Entry, yearBottle *widget.Entry, priceBottle *widget.Entry, pictureBottle *widget.Button, location map[string]int, IndividualBottle data.Bottle, identifier string) *widget.Form {
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
			timeOfCreationOrUpdate, _ := time.Parse("2006-01-02 15:04:05", time.Now().String())
			quantity, _ := strconv.ParseInt(quantityBottle.Text, 10, 0)
			bottleStorageLocation := make([]data.BottleStorageLocation, 0)
			storageLocation := data.IndStorageLocation

			for i := 0; i < 1; i++ {
				storageLocation.Name = storageLocationData.Text
				storageLocation.ID = location[storageLocationData.Text]
				dataToSent := data.BottleStorageLocation{
					StorageLocation: storageLocation,
					Quantity:        int(quantity),
				}
				bottleStorageLocation = append(bottleStorageLocation, dataToSent)
			}

			bottle := &data.Bottle{
				ID:                     IndividualBottle.ID,
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

			jsonValue := convertToJson(bottle)

			updateBottle(identifier, jsonValue)
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Mettre à jour",
		CancelText: "Annuler",
	}
	return updateForm
}

// Fill form fields with fetched data
func FillUpdateForm(nameBottle *widget.Entry, IndividualBottle data.Bottle, detailsBottle *widget.Entry, wineTypeBottle *widget.SelectEntry, storageLocationData *widget.SelectEntry, storageLocationForm data.StorageLocation, volumeBottle *widget.Entry, yearBottle *widget.Entry, priceBottle *widget.Entry, alcoholBottle *widget.Entry) {
	nameBottle.SetText(IndividualBottle.FullName)
	detailsBottle.SetText(strings.Replace(IndividualBottle.Description, "\\n", "\n", -1))
	wineTypeBottle.SetPlaceHolder(IndividualBottle.WineType)
	storageLocationData.SetPlaceHolder(storageLocationForm.Name)
	volumeBottle.SetText(strconv.Itoa(IndividualBottle.Volume))
	yearBottle.SetText(strconv.Itoa(IndividualBottle.YearProduced))
	priceBottle.SetText(fmt.Sprintf("%f", IndividualBottle.CurrentPrice))
	alcoholBottle.SetText(fmt.Sprintf("%f", IndividualBottle.AlcoholPercentage))
}

// Send data to API
func updateBottle(identifier string, jsonValue []byte) {
	postData := data.AuthPostRequest("Bottle/UpdateBottle/"+identifier, bytes.NewBuffer(jsonValue))
	if postData != 200 {
		message := "Error on bottle " + identifier + " update " + " StatusCode " + strconv.Itoa(postData)
		log(true, source, message)
	} else if postData == 200 {
		fmt.Println("Bottle updated")
	}
}

// Convert to JSON
func convertToJson(bottle *data.Bottle) []byte {
	jsonValue, err := json.Marshal(bottle)
	if err != nil {
		log(true, source, err.Error())
		fmt.Println(err)
	}
	return jsonValue
}

// Form to add and send a new bottle to the API endpoint (POST /api/producer)
func addNewTestBottle(_ fyne.Window) fyne.CanvasObject {

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
					CurrentPrice:      float32(price),
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

func getAllLocationName() map[string]int {
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

	LocationNameIdMap := make(map[string]int)
	for _, v := range storageLocationData {
		LocationNameIdMap[v.Name] = v.Id
	}
	return LocationNameIdMap
}
