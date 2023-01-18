package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"strconv"
)

func bottleForm(_ fyne.Window) fyne.CanvasObject {
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
				Description:       descriptionBottle.Text,
				Label:             labelBottle.Text,
				Volume:            volume,
				YearProduced:      year,
				AlcoholPercentage: alcohol,
				CurrentPrice:      price,

				CreatedBy: createdByBottle.Text,
			}

			// encode the value as JSON and send it to the API.
			bottleJsonValue, _ := json.Marshal(bottle)
			bottleResp, err := http.Post(bottleUrl, "application/json", bytes.NewBuffer(bottleJsonValue))
			if err != nil {
				fmt.Println("error while encoding response")
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error creating bottle: " + err.Error(),
				})
				return
			}
			if bottleResp.StatusCode != http.StatusCreated {
				fmt.Println(bottleJsonValue)
				fmt.Println("Erreur à l'envoi du formulaire")

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error creating producer: " + bottleResp.Status,
				})
				return
			}
			fmt.Println("New bottle added with success")

			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Content: "Producer created successfully!",
			})
		},
	}
	form.Append("Description", descriptionBottle)
	return form
}
