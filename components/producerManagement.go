package components

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
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Producer struct holds information about a producer
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

type PartialProducer struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_By"`
}

var individual Producer
var ProducerData []PartialProducer
var BindProducer []binding.DataMap
var producers []Producer

// makeProducerTabs function creates a new set of tabs
func makeProducerTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des producteurs", displayProducers(nil)),
		container.NewTabItem("Ajouter un producteur", producerForm(nil)),
		container.NewTabItem("Historique des commandes", displayOrders(nil)),
		container.NewTabItem("Passer une commande", producerOrdersForm(nil)),
		container.NewTabItem("Détails", widget.NewLabel("Content of tab 3")),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func producerAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	producerUrl := env.SERVER + "/api/producer"
	return producerUrl
}

// Call producer API and return the list of all producers
func fetchProducers() {
	apiUrl := producerAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&producers); err != nil {
		fmt.Println(err)
	}
}

// Call producer API and return producer matching ID
func fetchIndividual(id string) io.ReadCloser {
	apiUrl := producerAPIConfig() + "/" + id
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	return res.Body
}

// BETA Display API call result in a dynamic table
var ProducerColumns = []rtable.ColAttr{
	{ColName: "Name", Header: "Nom", WidthPercent: 150},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
}

func betaProducerTable(_ fyne.Window) fyne.CanvasObject {

	var id string

	chooseId := widget.NewEntryWithData(binding.BindString(&id))
	chooseId.SetText("Saississez un identifiant...")

	nameProducer := widget.NewEntry()
	nameProducer.Resize(fyne.NewSize(400, 35))
	nameProducer.Move(fyne.NewPos(10, 200))

	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(400, 35))
	formTitle.Move(fyne.NewPos(10, 150))

	detailsProducer := widget.NewMultiLineEntry()
	detailsProducer.Resize(fyne.NewSize(400, 100))
	detailsProducer.Move(fyne.NewPos(10, 250))

	createdByProducer := widget.NewEntry()
	createdByProducer.Resize(fyne.NewSize(400, 35))
	createdByProducer.Move(fyne.NewPos(10, 370))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(400, 50))
	submitBtn.Move(fyne.NewPos(10, 420))

	apiUrl := producerAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&ProducerData); err != nil {
		fmt.Println(err)
	}
	fmt.Print(ProducerData)

	for i := 0; i < len(ProducerData); i++ {
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
				binding.BindStruct(&Producer{Name: "Belle Ambiance",
					Details: "brown", CreatedBy: "170"}))
			tableOptions.Bindings = BindProducer
			table.Refresh()
			return
		}
		//Handle non-header row clicked
		str, err := rtable.GetStrCellValue(cell, tableOptions)
		str2, err := rtable.GetStrCellValue(cell, tableOptions)

		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}
		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]
		associatedRows := tableOptions.Bindings[cell.Col+1]

		colBinding, err := associatedRows.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		cellBinding, err := rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}

		err = colBinding.(binding.String).Set(rvsString(str2))
		err = cellBinding.(binding.String).Set(rvsString(str))
		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}

		fmt.Println("-->", str)
		fmt.Println(str2)

		nameProducer.SetText(str)
		details := string(individual.Details)
		details = strings.Replace(individual.Details, "\\n", "\n", -1)
		detailsProducer.SetText(details)
		createdByProducer.SetText(individual.CreatedBy)
	}

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewWithoutLayout(nameProducer, detailsProducer, createdByProducer, formTitle, submitBtn)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Display API call result in a table
func displayProducers(w fyne.Window) fyne.CanvasObject {

	//apiUrl := producerAPIConfig()

	idProducer := widget.NewEntry()
	idProducer.SetText("1")
	idProducer.Resize(fyne.NewSize(400, 35))
	idProducer.Move(fyne.NewPos(100, 100))

	nameProducer := widget.NewEntry()
	nameProducer.SetText("Belle Ambiance")
	nameProducer.Resize(fyne.NewSize(400, 35))
	nameProducer.Move(fyne.NewPos(100, 150))

	detailsProducer := widget.NewEntry()
	detailsProducer.SetText("Wine producer")
	detailsProducer.Resize(fyne.NewSize(400, 100))
	detailsProducer.Move(fyne.NewPos(100, 200))

	createdByProducer := widget.NewEntry()
	createdByProducer.SetText("negosud")
	createdByProducer.Resize(fyne.NewSize(400, 35))
	createdByProducer.Move(fyne.NewPos(100, 315))

	fetchProducers()
	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(400, 35))
	formTitle.Move(fyne.NewPos(100, 50))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(400, 50))
	submitBtn.Move(fyne.NewPos(100, 380))

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

	table.SetRowHeight(2, 50)

	deleteBtn := canvas.NewText("Supprimer ce producteur", color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 1,
	})
	deleteBtn.TextStyle = fyne.TextStyle{Bold: true}
	deleteBtn.TextSize = 20
	deleteBtn.Resize(fyne.NewSize(400, 50))
	deleteBtn.Move(fyne.NewPos(100, 500))

	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	//rightContainer := container.NewGridWithRows(2, form, deleteBtn)
	rightContainer := container.NewWithoutLayout(formTitle, idProducer, nameProducer, detailsProducer, createdByProducer, createdByProducer, submitBtn, deleteBtn)

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

// Form to add and send a new producer to the API endpoint (POST)
func producerForm(w fyne.Window) fyne.CanvasObject {

	apiUrl := producerAPIConfig()

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
				fmt.Println("Could not send form")
				producerFailureDialog(w)
				return
			}
			producerSuccessDialog(w)
			fmt.Println("New producer added with success")
		},
	}
	form.Append("Details", detailsProducer)
	mainContainer := container.NewVBox(title, form)

	return mainContainer
}

func displayIndividualProducer(w fyne.Window) fyne.CanvasObject {
	var id string

	chooseId := widget.NewEntryWithData(binding.BindString(&id))
	chooseId.SetText("Saississez un identifiant...")

	nameProducer := widget.NewEntry()
	nameProducer.Resize(fyne.NewSize(400, 35))
	nameProducer.Move(fyne.NewPos(10, 200))

	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(400, 35))
	formTitle.Move(fyne.NewPos(10, 150))

	detailsProducer := widget.NewMultiLineEntry()
	detailsProducer.Resize(fyne.NewSize(400, 100))
	detailsProducer.Move(fyne.NewPos(10, 250))

	createdByProducer := widget.NewEntry()
	createdByProducer.Resize(fyne.NewSize(400, 35))
	createdByProducer.Move(fyne.NewPos(10, 370))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(400, 50))
	submitBtn.Move(fyne.NewPos(10, 420))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Choix ID", Widget: chooseId},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			id = chooseId.Text

			resultApi := fetchIndividual(id)

			if err := json.NewDecoder(resultApi).Decode(&individual); err != nil {
				fmt.Println(err)
			}
			defer resultApi.Close()

			nameProducer.SetText(individual.Name)
			details := string(individual.Details)
			details = strings.Replace(individual.Details, "\\n", "\n", -1)
			detailsProducer.SetText(details)
			createdByProducer.SetText(individual.CreatedBy)
		},
	}

	form.Move(fyne.NewPos(400, 50))

	mainContainer := container.NewWithoutLayout(form, nameProducer, detailsProducer, createdByProducer, formTitle, submitBtn)

	return mainContainer
}

func individualPresentation() {
	details := string(individual.Details)
	details = strings.Replace(individual.Details, "\\n", "\n", -1)

	displayDetails := widget.NewTextGridFromString(details)
	displayName := canvas.NewText(individual.Name, color.Black)
	displayName.TextSize = 35
	displayName.Move(fyne.NewPos(10, 50))
	displayCreator := canvas.NewText("Ajouté par "+individual.CreatedBy, color.Black)
	displayName.TextSize = 35
	displayName.Move(fyne.NewPos(10, 50))

	displayCreator.TextSize = 15
	displayCreator.Move(fyne.NewPos(10, 130))
	displayDetails.Move(fyne.NewPos(10, 180))

}
