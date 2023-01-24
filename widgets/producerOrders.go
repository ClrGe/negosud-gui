package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"negosud-gui/config"
	"net/http"
	"strconv"
)

func orderAPIConfig() string {
	env, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	orderUrl := env.SERVER + "/api/producer"
	return orderUrl
}

// Display the list of orders fetched from API in a table
func displayOrders(w fyne.Window) fyne.CanvasObject {

	idProducer := widget.NewEntry()
	idProducer.SetText("1")

	nameProducer := widget.NewEntry()
	nameProducer.SetText("Belle Ambiance")

	detailsProducer := widget.NewEntry()
	detailsProducer.SetText("Wine producer")

	createdByProducer := widget.NewEntry()
	createdByProducer.SetText("negosud")

	config.FetchProducers()

	table := widget.NewTable(
		func() (int, int) { return 500, 150 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row >= len(Producers) {
				return
			}
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", Producers[id.Row].ID))
			case 1:
				label.SetText(Producers[id.Row].Name)
			case 2:
				label.SetText(Producers[id.Row].Details)
			case 3:
				label.SetText(Producers[id.Row].CreatedBy)
			case 4:
				label.SetText(fmt.Sprintf("%v", Producers[id.Row].CreatedAt))
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

// form to place a new order to a producer
func producerOrdersForm(w fyne.Window) fyne.CanvasObject {

	apiUrl := orderAPIConfig()

	idProducer := widget.NewEntry()
	nameProducer := widget.NewEntry()
	productToOrder := widget.NewEntry()
	quantity := widget.NewEntry()
	comment := widget.NewMultiLineEntry()

	title := widget.NewLabelWithStyle("PASSER UNE NOUVELLE COMMANDE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID producteur", Widget: idProducer},
			{Text: "Nom producteur", Widget: nameProducer},
			{Text: "Produits", Widget: productToOrder},
			{Text: "Quantité", Widget: quantity},
			{Text: "Commentaire", Widget: comment},
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		OnSubmit: func() {
			id, err := strconv.Atoi(idProducer.Text)
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error converting ID: " + err.Error(),
				})
				return
			}
			producer := &config.Producer{
				ID:        id,
				Name:      nameProducer.Text,
				Details:   comment.Text,
				CreatedBy: quantity.Text,
			}
			jsonValue, _ := json.Marshal(producer)
			resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonValue))

			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Content: "Error placing order: " + err.Error(),
				})
				return
			}
			if resp.StatusCode == 204 {
				fmt.Println("Could not send form")
				config.ProducerFailureDialog(w)
				return
			}
			config.ProducerSuccessDialog(w)
			fmt.Println("New order placed with success")
		},
	}
	mainContainer := container.NewVBox(title, form)

	return mainContainer
}
