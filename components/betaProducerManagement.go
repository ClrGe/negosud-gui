package components

import (
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
	"net/http"
	"strings"
)

// Producer struct holds information about a producer
type PartialProducer struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_By"`
}

var individual Producer
var ProducerData []PartialProducer
var BindProducer []binding.DataMap

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
