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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
	"github.com/rohanthewiz/rtable"
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
func displayAndUpdateProducers(win fyne.Window) fyne.CanvasObject {

	// retrieve structs from data package
	Individual := data.Individual
	ProducerData := data.ProducerData

	var identifier string
	var yPos, heightFields, widthForm float32

	yPos = 200
	heightFields = 50

	// DETAILS PRODUCER
	// declare elements (empty or hidden until an identifier in the table gets clicked on)
	instructions := widget.NewLabelWithStyle("Cliquez sur un identifiant dans le tableau pour afficher les détails", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	instructions.Resize(fyne.NewSize(widthForm, heightFields))
	instructions.Move(fyne.NewPos(0, yPos-500))
	productImg := canvas.NewImageFromFile("media/wineyard.jpeg")
	productImg.FillMode = canvas.ImageFillContain
	productImg.SetMinSize(fyne.NewSize(600, 340))
	productImg.Hidden = true
	productTitle := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	productTitle.Resize(fyne.NewSize(widthForm, heightFields))
	productTitle.Move(fyne.NewPos(0, yPos-300))
	productDesc := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	productDesc.Resize(fyne.NewSize(widthForm, heightFields))
	productDesc.Move(fyne.NewPos(0, yPos-250))

	// UPDATE FORM
	// declare form elements
	nameLabel := widget.NewLabel("Nom")
	nameProducer := widget.NewEntry()
	detailsLabel := widget.NewLabel("Description")
	detailsProducer := widget.NewMultiLineEntry()
	pictureLabel := widget.NewLabel("Image")
	pictureProducer := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	deleteBtn := widget.NewButtonWithIcon("Supprimer ce producteur", theme.WarningIcon(),
		func() {})

	// retrieve API url and fetch bottle data
	apiUrl := data.ProducerAPIConfig()
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.NewDecoder(res.Body).Decode(&ProducerData); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(ProducerData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		t := ProducerData[i]
		id := strconv.Itoa(t.Id)
		ProducerData[i].ID = id

		// binding producer data
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
		identifier, err = rtable.GetStrCellValue(cell, tableOptions)

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
		} else {
			instructions.Hidden = true
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

		productTitle.SetText(Individual.Name)
		productDesc.SetText(details)
	}

	updateForm := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameProducer},
			{Text: "", Widget: detailsLabel},
			{Text: "", Widget: detailsProducer},
			{Text: "", Widget: pictureLabel},
			{Text: "", Widget: pictureProducer},
		},
		OnSubmit: func() {
			producer := &data.Producer{
				Name:    nameProducer.Text,
				Details: detailsProducer.Text,
			}
			jsonValue, _ := json.Marshal(producer)
			updateEndpoint := data.UpdateProducerAPI()
			resp, err := http.Post(updateEndpoint+identifier, "application/json", bytes.NewBuffer(jsonValue))

			if err != nil {
				fmt.Println("Could not send form")
				return
			}
			if resp.StatusCode == 204 {
				fmt.Println("Response : no content (204)")
				return
			}
			fmt.Println("Producer updated")
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Envoyer",
		CancelText: "Annuler",
	}

	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewWithoutLayout(productTitle, productDesc))

	layoutDetailsTab := container.NewBorder(image, nil, nil, nil, textProduct, instructions)
	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithDelete := container.NewBorder(layoutUpdateForm, deleteBtn, nil, nil)

	// Define layout
	individualTabs := container.NewAppTabs(
		container.NewTabItem("Détails du producteur", layoutDetailsTab),
		container.NewTabItem("Modifier le producteur", layoutWithDelete),
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

	nameLabel := widget.NewLabel("Nom")
	nameProducer := widget.NewEntry()
	detailsLabel := widget.NewLabel("Description")
	detailsProducer := widget.NewMultiLineEntry()
	pictureLabel := widget.NewLabel("Image")
	pictureProducer := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	title := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: title},
			{Text: "", Widget: nameLabel},
			{Text: "", Widget: nameProducer},
			{Text: "", Widget: detailsLabel},
			{Text: "", Widget: detailsProducer},
			{Text: "", Widget: pictureLabel},
			{Text: "", Widget: pictureProducer},
		},
		OnSubmit: func() {

			producer := &data.Producer{
				Name:    nameProducer.Text,
				Details: detailsProducer.Text,
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
		SubmitText: "Envoyer",
	}
	mainContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(900, 600), form))

	return mainContainer
}

func contactProducers(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Messages échangés avec les producteurs (à implémenter)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	))
}
