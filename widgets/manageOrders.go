package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"negosud-gui/data"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
)

// TODO : ADD FUNCTIONNALITY "ORDER AGAIN" ON THE LIST OF ORDERS
// TODO : SEND MAIL TO PRODUCER AUTOMATICALLY WHEN PLACING ORDER

var BindOrder []binding.DataMap

// makeOrdersTabs function creates a new set of tabs
func makeOrdersTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Historique des commandes", displayOrders(nil)),
		container.NewTabItem("Passer une commande", producerOrdersForm(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// OrdersColumns defines the header row for the table
var OrdersColumns = []rtable.ColAttr{
	{ColName: "Date", Header: "Date", WidthPercent: 80},
	{ColName: "Producer", Header: "ID Producteur", WidthPercent: 50},
	{ColName: "Product", Header: "ID Produit", WidthPercent: 50},
	{ColName: "Quantity", Header: "Quantité", WidthPercent: 50},
	{ColName: "Price", Header: "Prix unité (€)", WidthPercent: 50},
	{ColName: "Price", Header: "Prix total(€)", WidthPercent: 50},
	{ColName: "Status", Header: "Statut de la commande", WidthPercent: 120},
}

// Display the list of orders fetched from API in a table
func displayOrders(_ fyne.Window) fyne.CanvasObject {
	Orders := data.Orders

	response := data.AuthGetRequest("orders")
	if response == nil {
		message := "Request body returned empty"
		fmt.Println(message)
		data.Logger(false, "WIDGETS.ORDERS", message)
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&Orders); err != nil {
		data.Logger(true, "WIDGETS.ORDERS", err.Error())
		fmt.Println(err)
	}

	BindOrder = nil

	for i := 0; i < len(Orders); i++ {
		id := strconv.Itoa(Orders[i].Id)
		q := strconv.Itoa(Orders[i].QuantityInt)
		p := strconv.Itoa(Orders[i].ProducerId)
		b := strconv.Itoa(Orders[i].ProductId)
		Orders[i].Product = b
		Orders[i].Producer = p
		Orders[i].Quantity = q
		Orders[i].ID = id
		BindOrder = append(BindOrder, binding.BindStruct(&Orders[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: OrdersColumns,
		Bindings: BindOrder,
	}
	table := rtable.CreateTable(tableOptions)
	return table
}

// form to place a new order to a producer
func producerOrdersForm(_ fyne.Window) fyne.CanvasObject {
	nameLabel := widget.NewLabelWithStyle("Nom du producteur", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	nameProducer := widget.NewEntry()
	nameProducer.SetPlaceHolder("Jean Bon")
	productLabel := widget.NewLabelWithStyle("Nom du produit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	productToOrder := widget.NewEntry()
	productToOrder.SetPlaceHolder("Pinard")
	quantityLabel := widget.NewLabelWithStyle("Quantité", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	quantity := widget.NewSelectEntry([]string{"10", "20", "50", "100", "150", "200", "250", "300", "350", "400", "450"})
	quantity.SetPlaceHolder("Sélectionnez une quantité...")
	commentLabel := widget.NewLabelWithStyle("Commentaire (facultatif)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	comment := widget.NewMultiLineEntry()
	comment.SetPlaceHolder("Votre commentaire...")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameProducer},
			{Text: "", Widget: productLabel},
			{Text: " ", Widget: productToOrder},
			{Text: "", Widget: quantityLabel},
			{Text: "", Widget: quantity},
			{Text: "", Widget: commentLabel},
			{Text: "", Widget: comment},
		},
		OnSubmit: func() {

			// extract the value from the input widget and set the corresponding field in the Producer struct
			newOrder := &data.Order{
				Producer: nameProducer.Text,
				Product:  productToOrder.Text,
				Quantity: quantity.Text,
				Comment:  comment.Text,
			}
			// encode the value as JSON and send it to the API.
			jsonValue, _ := json.Marshal(newOrder)
			postData := data.AuthPostRequest("orders", bytes.NewBuffer(jsonValue))
			if postData != 201|200 {
				message := "Could not place order"
				fmt.Println(message)
				data.Logger(false, "WIDGETS.ORDERS", message)

			}
			message := "Order placed successfully"
			data.Logger(false, "WIDGETS.ORDERS", message)
		},
		SubmitText: "Envoyer",
	}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))
	return mainContainer
}
