package components

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"net/http"
	"strings"
)

var individual Producer

func fetchIndividual(id string) io.ReadCloser {
	apiUrl := producerAPIConfig() + "/" + id
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	return res.Body
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
