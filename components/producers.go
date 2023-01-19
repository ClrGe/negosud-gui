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

// Define the producer struct and associate json fields
type Producer struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Details   string      `json:"details"`
	CreatedAt interface{} `json:"created_At"`
	UpdatedAt time.Time   `json:"updated_At"`
	CreatedBy string      `json:"created_By"`
	UpdatedBy string      `json:"updated_By"`
	Bottles   interface{} `json:"bottles"`
	Region    interface{} `json:"region"`
}

var producers []Producer
var producer []Producer

// Display API call result in a table
func displayProducers(w fyne.Window) fyne.CanvasObject {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	apiUrl := env.SERVER + "/api/producer"

	idProducer := widget.NewEntry()
	idProducer.SetText("1")

	nameProducer := widget.NewEntry()
	nameProducer.SetText("Belle Ambiance")

	detailsProducer := widget.NewEntry()
	detailsProducer.SetText("Wine producer")

	createdByProducer := widget.NewEntry()
	createdByProducer.SetText("negosud")

	// call API - this function is defined in callAPI.go
	fetchProducers()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Id", Widget: idProducer},
			{Text: "Nom du producteur", Widget: nameProducer},
			{Text: "Ajouté par", Widget: createdByProducer},
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
			if resp.StatusCode == 204 {
				fmt.Println(jsonValue)
				fmt.Println("Erreur à l'envoi du formulaire")

				producerFailureDialog(w)
				return
			}
			producerSuccessDialog(w)
			fmt.Println("New producer added with success")
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
				label.SetText(fmt.Sprintf("%d", producers[id.Row].ID))
			case 1:
				label.SetText(producers[id.Row].Name)
			case 2:
				label.SetText(producers[id.Row].Details)
			case 3:
				label.SetText(producers[id.Row].CreatedBy)
			case 4:
				label.SetText(fmt.Sprintf("%v", producers[id.Row].CreatedAt))
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
