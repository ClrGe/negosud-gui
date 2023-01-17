package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/http"
)

var newBottle []Bottle

func postNewBottle(newBottle Bottle) error {
	env, err := LoadConfig(".")
	if err != nil {
		return err
	}
	// convert producer struct to json
	producerJSON, err := json.Marshal(newBottle)
	if err != nil {
		return err
	}
	// create http client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", env.SERVER+"/api/bottle", bytes.NewBuffer(producerJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	// make request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		return fmt.Errorf("error posting new bottle, status code: %d", res.StatusCode)
	}
	return nil
}

func bottleForm(_ fyne.Window) fyne.CanvasObject {
	name := widget.NewEntry()

	prod := widget.NewEntry()

	largeText := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Nom", Widget: name},
			{Text: "Producteur", Widget: prod},
		},
		OnCancel: func() {
			fmt.Println("Annulé")
		},
		OnSubmit: func() {
			fmt.Println("Formulaire envoyé")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Content: largeText.Text,
			})
		},
	}
	form.Append("Description", largeText)
	return form
}
