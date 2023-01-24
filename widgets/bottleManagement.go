package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/config"
	"net/http"
	"strconv"
)

// The makeBottleTabs function creates a new set of tabs
func makeBottleTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", betaBottleTable(nil)),
		container.NewTabItem("Ajouter un produit", addNewBottle(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// The displayBottles function calls the fetchBottles function to display a table with the list of Bottles.
func displayBottles(w fyne.Window) fyne.CanvasObject {
	config.FetchBottles()

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := bottlesTable(w)
	rightContainer := container.NewGridWithRows(2, updateBottle(w), deleteBottle())

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)
	return mainContainer
}

func bottlesTable(w fyne.Window) fyne.CanvasObject {
	Bottles := config.Bottles

	table := widget.NewTable(
		func() (int, int) { return 500, 150 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row >= len(Bottles) {
				return
			}
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%v", Bottles[id.Row].ID))
			case 1:
				label.SetText(Bottles[id.Row].FullName)
			case 2:
				label.SetText(Bottles[id.Row].Label)
			case 3:
				label.SetText(Bottles[id.Row].CreatedBy)
			case 4:
				label.SetText(fmt.Sprintf("%v", Bottles[id.Row].CurrentPrice))
			case 5:
				label.SetText(fmt.Sprintf("%v", Bottles[id.Row].Volume))
			}
		})

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 200)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 200)
	table.SetRowHeight(2, 50)

	return table
}

// Form to add and send a new bottle to the API endpoint (POST)
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

func updateBottle(w fyne.Window) fyne.CanvasObject {
	idBottle := widget.NewEntry()
	idBottle.SetText("25")

	nameBottle := widget.NewEntry()
	nameBottle.SetText("Cooper Ranch")

	descriptionBottle := widget.NewEntry()
	descriptionBottle.SetText("Vin de caractère")

	labelBottle := widget.NewEntry()
	labelBottle.SetText("Biologique")

	yearBottle := widget.NewEntry()
	yearBottle.SetText("2020")

	volumeBottle := widget.NewEntry()
	volumeBottle.SetText("75")

	alcoholBottle := widget.NewEntry()
	alcoholBottle.SetText("13")

	currentPriceBottle := widget.NewEntry()
	currentPriceBottle.SetText("16")

	createdByBottle := widget.NewEntry()
	createdByBottle.SetText("negosud")

	form := &widget.Form{
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

			// encode the value as JSON
			bottleJsonValue, _ := json.Marshal(bottle)
			apiUrl := config.BottleAPIConfig()
			// send it to the API
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
	return form
}

func deleteBottle() fyne.CanvasObject {
	deleteButton := widget.NewButton("Supprimer", func() {
		fmt.Println("Deleted")
	})
	return deleteButton
}
