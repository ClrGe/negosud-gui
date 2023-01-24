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

// ProducerColumns defines the header row for the table
var ProducerColumns = []rtable.ColAttr{
	{ColName: "Name", Header: "Nom", WidthPercent: 150},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
}

// betaProducerTable implements a dynamic table bound to an editing form
func betaProducerTable(_ fyne.Window) fyne.CanvasObject {
	nameProducer := widget.NewEntry()
	nameProducer.Resize(fyne.NewSize(400, 35))
	nameProducer.Move(fyne.NewPos(150, 200))

	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(400, 35))
	formTitle.Move(fyne.NewPos(150, 150))

	detailsProducer := widget.NewMultiLineEntry()
	detailsProducer.Resize(fyne.NewSize(400, 100))
	detailsProducer.Move(fyne.NewPos(150, 250))

	createdByProducer := widget.NewEntry()
	createdByProducer.Resize(fyne.NewSize(400, 35))
	createdByProducer.Move(fyne.NewPos(150, 370))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(400, 50))
	submitBtn.Move(fyne.NewPos(150, 420))

	apiUrl := config.ProducerAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&config.ProducerData); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(config.ProducerData); i++ {
		BindProducer = append(BindProducer, binding.BindStruct(&config.ProducerData[i]))
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
		str, err := rtable.GetStrCellValue(cell, tableOptions)

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
		//err = cellBinding.(binding.String).Set(rvsString(str))
		//if err != nil {
		//	fmt.Println(rerr.StringFromErr(err))
		//	return
		//}

		fmt.Println("-->", str)

		nameProducer.SetText(str)
		details := config.Individual.Details
		details = strings.Replace(config.Individual.Details, "\\n", "\n", -1)
		detailsProducer.SetText(details)
		createdByProducer.SetText(config.Individual.CreatedBy)
	}

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewWithoutLayout(nameProducer, detailsProducer, createdByProducer, formTitle, submitBtn)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

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
