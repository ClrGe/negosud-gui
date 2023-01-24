package widgets

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"negosud-gui/config"
	"strings"
)

var Producers []config.Producer

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

	config.FetchProducers()
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
			if id.Row >= len(Producers) {
				return
			}
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", Producers[id.Row].ID))
			case 1:
				label.SetText(Producers[id.Row].Name)
			case 2:
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

func individualPresentation() {
	Individual := config.Individual

	details := string(Individual.Details)
	details = strings.Replace(Individual.Details, "\\n", "\n", -1)

	displayDetails := widget.NewTextGridFromString(details)
	displayName := canvas.NewText(Individual.Name, color.Black)
	displayName.TextSize = 35
	displayName.Move(fyne.NewPos(10, 50))
	displayCreator := canvas.NewText("Ajouté par "+Individual.CreatedBy, color.Black)
	displayName.TextSize = 35
	displayName.Move(fyne.NewPos(10, 50))

	displayCreator.TextSize = 15
	displayCreator.Move(fyne.NewPos(10, 130))
	displayDetails.Move(fyne.NewPos(10, 180))
}

func displayIndividualProducer(w fyne.Window) fyne.CanvasObject {
	var id string
	Individual := config.Individual

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

			resultApi := config.FetchIndividual(id)

			if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
				fmt.Println(err)
			}
			defer resultApi.Close()

			nameProducer.SetText(Individual.Name)
			details := string(Individual.Details)
			details = strings.Replace(Individual.Details, "\\n", "\n", -1)
			detailsProducer.SetText(details)
			createdByProducer.SetText(Individual.CreatedBy)
		},
	}

	form.Move(fyne.NewPos(400, 50))

	mainContainer := container.NewWithoutLayout(form, nameProducer, detailsProducer, createdByProducer, formTitle, submitBtn)

	return mainContainer
}
