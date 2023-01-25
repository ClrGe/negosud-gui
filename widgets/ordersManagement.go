package widgets

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
	"negosud-gui/config"
	"net/http"
	"strconv"
)

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
func displayOrders(w fyne.Window) fyne.CanvasObject {

	Orders := config.Orders

	apiUrl := config.OrderAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&Orders); err != nil {
		fmt.Println(err)
	}

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
func producerOrdersForm(w fyne.Window) fyne.CanvasObject {

	nameProducer := widget.NewEntry()
	nameProducer.SetPlaceHolder("Nom du producteur...")
	productToOrder := widget.NewEntry()
	productToOrder.SetPlaceHolder("Nom du produit...")
	quantity := widget.NewSelectEntry([]string{"10", "20", "50", "100", "150", "200", "250", "300", "350", "400", "450"})
	quantity.SetPlaceHolder("Sélectionnez une quantité...")
	comment := widget.NewMultiLineEntry()
	comment.SetPlaceHolder("Ajouter un commentaire (facultatif)")

	title := widget.NewLabelWithStyle("PASSER UNE NOUVELLE COMMANDE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Nom du producteur", Widget: nameProducer},
			{Text: "Nom du produit", Widget: productToOrder},
			{Text: "Quantité", Widget: quantity},
			{Text: "Commentaire", Widget: comment},
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		OnSubmit: func() {
			fmt.Println("New order placed with success")
		},
	}
	form.Resize(fyne.NewSize(600, 1000))
	left := container.NewVBox(title, form)
	right := container.NewVBox()
	mainContainer := container.New(layout.NewGridLayout(2))

	mainContainer.Add(left)
	mainContainer.Add(right)

	return mainContainer
}
