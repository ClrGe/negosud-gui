package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"strconv"
	"time"
)

// Define the Bottle struct and associate json fields
type Bottle struct {
	ID                int         `json:"id"`
	FullName          string      `json:"full_Name"`
	Description       string      `json:"description"`
	Label             string      `json:"label"`
	Volume            int         `json:"volume"`
	Picture           string      `json:"picture"`
	YearProduced      int         `json:"year_Produced"`
	AlcoholPercentage int         `json:"alcohol_Percentage"`
	CurrentPrice      int         `json:"current_Price"`
	CreatedAt         time.Time   `json:"created_At"`
	UpdatedAt         time.Time   `json:"updated_At"`
	CreatedBy         string      `json:"created_By"`
	UpdatedBy         string      `json:"updated_By"`
	BottleLocations   interface{} `json:"bottleLocations"`
	BottleGrapes      interface{} `json:"bottleGrapes"`
	Producer          interface{} `json:"producer"`
}

func makeBottleTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", displayBottles(nil)),
		container.NewTabItem("Ajouter un produit", bottleForm(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

var bottles []Bottle

// Call bottle API and return the list of all bottles
func fetchBottles() {
	env, err := LoadConfig(".")
	res, err := http.Get(env.SERVER + "/api/bottle")

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&bottles); err != nil {
		fmt.Println(err)
	}
}

// Display API call result in a table
func displayBottles(w fyne.Window) fyne.CanvasObject {
	fetchBottles()
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}
	bottleUrl := env.SERVER + "/api/bottle"

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
			bottle := &Bottle{
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
			bottleResp, err := http.Post(bottleUrl, "application/json", bytes.NewBuffer(bottleJsonValue))
			if err != nil {
				fmt.Println("error while encoding response")
				return
			}
			if bottleResp.StatusCode == 204 {
				bottleFailureDialog(w)
				fmt.Println(bottleJsonValue)
				return
			}
			bottleSuccessDialog(w)
			fmt.Println("New bottle added with success")

		},
	}

	table := widget.NewTable(
		func() (int, int) { return 500, 150 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row >= len(producers) {
				return
			}
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%v", bottles[id.Row].ID))
			case 1:
				label.SetText(bottles[id.Row].FullName)
			case 2:
				label.SetText(bottles[id.Row].Label)
			case 3:
				label.SetText(bottles[id.Row].CreatedBy)
			case 4:
				label.SetText(fmt.Sprintf("%v", bottles[id.Row].CurrentPrice))
			case 5:
				label.SetText(fmt.Sprintf("%v", bottles[id.Row].Volume))
			}
		})

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 200)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 200)
	table.SetRowHeight(2, 50)

	dlt := widget.NewButton("Supprimer", func() {
		fmt.Println("Deleted")
	})

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewGridWithRows(2, form, dlt)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)
	return mainContainer
}

func bottleForm(w fyne.Window) fyne.CanvasObject {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	bottleUrl := env.SERVER + "/api/bottle"

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
				bottle := &Bottle{
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
				bottleResp, err := http.Post(bottleUrl, "application/json", bytes.NewBuffer(bottleJsonValue))
				if err != nil {
					fmt.Println("error while encoding response")
					return
				}
				if bottleResp.StatusCode == 204 {
					bottleFailureDialog(w)
					fmt.Println(bottleJsonValue)
					return
				}
				bottleSuccessDialog(w)
				fmt.Println("New bottle added with success")

			},
		}
	form.Append("Description", descriptionBottle)
	form.Size()
	mainContainer := container.NewVBox(title, form)

	return mainContainer
}
