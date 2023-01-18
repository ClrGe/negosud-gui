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

func producerForm(_ fyne.Window) fyne.CanvasObject {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}
	apiUrl := env.SERVER + "/api/producer"
	idProducer := widget.NewEntry()
	nameProducer := widget.NewEntry()
	detailsProducer := widget.NewEntry()
	createdByProducer := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idProducer},
			{Text: "Nom", Widget: nameProducer},
			{Text: "Created By", Widget: createdByProducer},
		},
		OnCancel: func() {
			fmt.Println("Annulation")
		},
		OnSubmit: func() {
			id, err := strconv.Atoi(idProducer.Text)
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error converting ID: " + err.Error(),
				})
				return
			}
			producer := &Producer{
				ID:        id,
				Name:      nameProducer.Text,
				Details:   detailsProducer.Text,
				CreatedBy: createdByProducer.Text,
			}
			jsonValue, _ := json.Marshal(producer)
			resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error creating producer: " + err.Error(),
				})
				return
			}
			if resp.StatusCode != http.StatusCreated {
				fmt.Println(jsonValue)
				fmt.Println("Erreur à l'envoi du formulaire")

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error creating producer: " + resp.Status,
				})
				return
			}
			fmt.Println("Nouveau producteur crée avec succès")

			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Content: "Producer created successfully!",
			})
		},
	}
	form.Append("Details", detailsProducer)
	return form
}
