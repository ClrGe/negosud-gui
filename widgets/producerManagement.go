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
	"negosud-gui/data"
	"net/http"
	"strconv"
	"strings"
)

var BindProducer []binding.DataMap

// makeProducerTabs function creates a new set of tabs
func makeProducerTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des producteurs", displayAndUpdateProducers(nil)),
		container.NewTabItem("Ajouter un producteur", addNewProducer(nil)),
		container.NewTabItem("Contact producteurs", contactProducers(nil)),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// ProducerColumns defines the header row for the table
var ProducerColumns = []rtable.ColAttr{
	{ColName: "ID", Header: "ID", WidthPercent: 50},
	{ColName: "Name", Header: "Nom", WidthPercent: 120},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
}

// displayAndUpdateProducers implements a dynamic table bound to an editing form
func displayAndUpdateProducers(_ fyne.Window) fyne.CanvasObject {

	// retrieve structs from data package
	Individual := data.Individual
	ProducerData := data.ProducerData

	var xPos, yPos, heightFields, heightLabels, widthForm float32

	xPos = -200
	yPos = 200
	heightFields = 50
	heightLabels = 20
	widthForm = 600

	// DETAILS PRODUCER

	instructions := canvas.NewText("Cliquez sur un identifiant dans le tableau pour afficher les détails du producteur.", color.Black)
	instructions.TextSize = 15
	instructions.TextStyle = fyne.TextStyle{Bold: true}
	instructions.Resize(fyne.NewSize(widthForm, heightFields))
	instructions.Move(fyne.NewPos(0, yPos-500))

	productImg := canvas.NewImageFromFile("media/wineyard.jpeg")
	productImg.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		productImg.SetMinSize(fyne.NewSize(800, 340))
	} else {
		productImg.SetMinSize(fyne.NewSize(800, 340))
	}
	productImg.Hidden = true

	productTitle := widget.NewLabel("")
	productTitle.Resize(fyne.NewSize(widthForm, heightFields))
	productTitle.Move(fyne.NewPos(0, yPos-400))

	productDesc := widget.NewLabel("")
	productDesc.Resize(fyne.NewSize(widthForm, heightFields))
	productDesc.Move(fyne.NewPos(0, yPos-350))

	// UPDATE FORM

	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(widthForm, heightFields))
	formTitle.Move(fyne.NewPos(xPos, yPos-600))

	nameLabel := canvas.NewText("Nom", color.Black)
	nameLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	nameLabel.Move(fyne.NewPos(xPos, yPos-520))
	nameProducer := widget.NewEntry()
	nameProducer.Resize(fyne.NewSize(widthForm, heightFields))
	nameProducer.Move(fyne.NewPos(xPos, yPos-500))

	detailsLabel := canvas.NewText("Description", color.Black)
	detailsLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	detailsLabel.Move(fyne.NewPos(xPos, yPos-420))
	detailsProducer := widget.NewMultiLineEntry()
	detailsProducer.Resize(fyne.NewSize(widthForm, heightFields*2))
	detailsProducer.Move(fyne.NewPos(xPos, yPos-400))

	createdByLabel := canvas.NewText("Ajouté par", color.Black)
	createdByLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	createdByLabel.Move(fyne.NewPos(xPos, yPos-280))
	createdByProducer := widget.NewEntry()
	createdByProducer.Resize(fyne.NewSize(widthForm, heightFields))
	createdByProducer.Move(fyne.NewPos(xPos, yPos-260))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(widthForm, heightFields))
	submitBtn.Move(fyne.NewPos(xPos, yPos-150))

	apiUrl := data.ProducerAPIConfig()

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
				binding.BindStruct(&data.Producer{Name: "Belle Ambiance",
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
		resultApi := data.FetchIndividualProducer(identifier)
		if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
			fmt.Println(err)
		} else {
			productImg.Hidden = false
		}
		// Fill form fields with fetched data
		nameProducer.SetText(Individual.Name)
		details := string(Individual.Details)
		details = strings.Replace(Individual.Details, "\\n", "\n", -1)
		detailsProducer.SetText(details)
		createdByProducer.SetText(Individual.CreatedBy)

		productTitle.SetText("Nom: " + Individual.Name)
		productDesc.SetText("Description: " + details)
	}
	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewWithoutLayout(productTitle, productDesc))
	detailsProduct := container.NewBorder(image, instructions, nil, nil, textProduct)

	updateForm := container.NewCenter(container.NewWithoutLayout(nameLabel, nameProducer, detailsLabel, detailsProducer, createdByLabel, createdByProducer, formTitle, submitBtn))

	// Define layout
	individualTabs := container.NewAppTabs(
		container.NewTabItem("Détails du producteur", detailsProduct),
		container.NewTabItem("Modifier le producteur", updateForm),
	)

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new producer to the API endpoint (POST)
func addNewProducer(w fyne.Window) fyne.CanvasObject {
	apiUrl := data.ProducerAPIConfig()

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
			producer := &data.Producer{
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
				data.ProducerFailureDialog(w)
				return
			}
			data.ProducerSuccessDialog(w)
			fmt.Println("New producer added with success")
		},
	}
	form.Append("Details", detailsProducer)
	formContainer := container.NewVBox(title, form)
	mainContainer := container.NewCenter(formContainer)

	return mainContainer
}

func contactProducers(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Messages échangés avec les producteurs (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
