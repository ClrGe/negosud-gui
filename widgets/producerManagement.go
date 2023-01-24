package widgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
	"github.com/rohanthewiz/rtable"
	"image/color"
	"negosud-gui/config"
	"net/http"
	"strconv"
	"strings"
)

var BindProducer []binding.DataMap

// makeProducerTabs function creates a new set of tabs
func makeProducerTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des producteurs", betaProducerTable(nil)),
		container.NewTabItem("Ajouter un producteur", addNewProducer(nil)),
		container.NewTabItem("Historique des commandes", displayOrders(nil)),
		container.NewTabItem("Passer une commande", producerOrdersForm(nil)),
		container.NewTabItem("Détails", widget.NewLabel("Content of tab 3")),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// ProducerColumns defines the header row for the table
var ProducerColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 50},
	{ColName: "Name", Header: "Nom", WidthPercent: 120},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
}

// betaProducerTable implements a dynamic table bound to an editing form
func betaProducerTable(_ fyne.Window) fyne.CanvasObject {

	// retrieve structs from config package
	Individual := config.Individual
	ProducerData := config.ProducerData

	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(600, 50))
	formTitle.Move(fyne.NewPos(50, 100))

	nameLabel := canvas.NewText("Nom", color.Black)
	nameLabel.Resize(fyne.NewSize(600, 20))
	nameLabel.Move(fyne.NewPos(50, 180))
	nameProducer := widget.NewEntry()
	nameProducer.Resize(fyne.NewSize(600, 50))
	nameProducer.Move(fyne.NewPos(50, 200))

	detailsLabel := canvas.NewText("Description", color.Black)
	detailsLabel.Resize(fyne.NewSize(600, 20))
	detailsLabel.Move(fyne.NewPos(50, 280))
	detailsProducer := widget.NewMultiLineEntry()
	detailsProducer.Resize(fyne.NewSize(600, 150))
	detailsProducer.Move(fyne.NewPos(50, 300))

	createdByLabel := canvas.NewText("Ajouté par", color.Black)
	createdByLabel.Resize(fyne.NewSize(600, 20))
	createdByLabel.Move(fyne.NewPos(50, 480))
	createdByProducer := widget.NewEntry()
	createdByProducer.Resize(fyne.NewSize(600, 50))
	createdByProducer.Move(fyne.NewPos(50, 500))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(600, 50))
	submitBtn.Move(fyne.NewPos(50, 600))

	apiUrl := config.ProducerAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&ProducerData); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(ProducerData); i++ {
		t := ProducerData[i]
		id := strconv.Itoa(t.Id)
		ProducerData[i].ID = id
		BindProducer = append(BindProducer, binding.BindStruct(&ProducerData[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: ProducerColumns,
		Bindings: BindProducer,
	}

	table := rtable.CreateTable(tableOptions)

	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 0 || cell.Row > len(BindProducer) { // 1st col is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(ProducerColumns) {
			fmt.Println("*-> Column out of limits")
			return
		}
		// Handle header row clicked
		if cell.Row == 0 { // fmt.Println("-->", tblOpts.ColAttrs[cell.Col].Header)
			// Add a row
			BindProducer = append(BindProducer,
				binding.BindStruct(&config.Producer{Name: "Belle Ambiance",
					Details: "brown", CreatedBy: "170"}))
			tableOptions.Bindings = BindProducer
			table.Refresh()
			return
		}
		//Handle non-header row clicked
		identifier, err := rtable.GetStrCellValue(cell, tableOptions)

		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}
		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]

		cellBinding, err := rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}
		fmt.Println(cellBinding)

		fmt.Println("-->", identifier)

		// Fetch individual producer to fill form
		resultApi := config.FetchIndividualProducer(identifier)
		if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
			fmt.Println(err)
		}
		// Fill form fields with fetched data
		nameProducer.SetText(Individual.Name)
		details := string(Individual.Details)
		details = strings.Replace(Individual.Details, "\\n", "\n", -1)
		detailsProducer.SetText(details)
		createdByProducer.SetText(Individual.CreatedBy)
	}

	// Define layout
	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewWithoutLayout(nameLabel, nameProducer, detailsLabel, detailsProducer, createdByLabel, createdByProducer, formTitle, submitBtn)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new producer to the API endpoint (POST)
func addNewProducer(w fyne.Window) fyne.CanvasObject {
	apiUrl := config.ProducerAPIConfig()

	idProducer := widget.NewEntry()
	nameProducer := widget.NewEntry()
	detailsProducer := widget.NewEntry()
	createdByProducer := widget.NewEntry()

	title := widget.NewLabelWithStyle("AJOUTER UN NOUVEAU PRODUCTEUR", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idProducer},
			{Text: "Nom", Widget: nameProducer},
			{Text: "Created By", Widget: createdByProducer},
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
				fmt.Println("Could not send form")
				config.ProducerFailureDialog(w)
				return
			}
			config.ProducerSuccessDialog(w)
			fmt.Println("New producer added with success")
		},
	}
	form.Append("Details", detailsProducer)
	mainContainer := container.NewVBox(title, form)

	return mainContainer
}
