package widgets

import (
	"encoding/json"
	"fmt"
	"negosud-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rtable"
)

var BindCustomerOrder []binding.DataMap

// makeOrdersTabs function creates a new set of tabs
func makeCusOrdersTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Historique des commandes clients", displayCustomersOrders(nil)),
		container.NewTabItem("Support clients", displayCustomersMessages(nil)),
	)

	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// CustomersOrdersColumns defines the header row for the table
var CustomersOrdersColumns = []rtable.ColAttr{
	{ColName: "Date", Header: "Date", WidthPercent: 80},
	{ColName: "Client", Header: "ID client", WidthPercent: 50},
	{ColName: "Product", Header: "ID Produit", WidthPercent: 50},
	{ColName: "Quantity", Header: "Quantité", WidthPercent: 50},
	{ColName: "PriceHT", Header: "Montant HT (€)", WidthPercent: 50},
	{ColName: "PriceHT", Header: "Montant total (€)", WidthPercent: 50},
	{ColName: "Status", Header: "Statut de la commande", WidthPercent: 120},
}

// Display the list of orders fetched from API in a table
func displayCustomersOrders(_ fyne.Window) fyne.CanvasObject {
	CustomerOrders := data.CustomerOrders

	response := data.AuthGetRequest("customers-orders")
	if response == nil {
		message := "Request body returned empty"
		fmt.Println(message)
		data.Logger(false, "WIDGETS.CUSTOMER-ORDERS", message)
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}

	if err := json.NewDecoder(response).Decode(&CustomerOrders); err != nil {
		fmt.Println(err)
		data.Logger(true, "WIDGETS.CUSTOMER-ORDERS", err.Error())

		return widget.NewLabel("Erreur de décodage du json")
	}

	BindCustomerOrder = nil

	for i := 0; i < len(CustomerOrders); i++ {
		BindCustomerOrder = append(BindCustomerOrder, binding.BindStruct(&CustomerOrders[i]))
	}
	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: CustomersOrdersColumns,
		Bindings: BindCustomerOrder,
	}
	table := rtable.CreateTable(tableOptions)

	return table
}

func displayCustomersMessages(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Messages envoyés par les clients du site e-commerce (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
