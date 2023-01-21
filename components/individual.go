package components

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net/http"
	"strings"
)

var individual Producer

func displayIndividualProducer(w fyne.Window) fyne.CanvasObject {
	apiUrl := producerAPIConfig() + "/100"
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&individual); err != nil {
		fmt.Println(err)
	}

	displayName := canvas.NewText(individual.Name, color.Black)
	displayName.TextSize = 35
	displayName.Move(fyne.NewPos(10, 50))

	displayCreator := canvas.NewText("Ajouté par "+individual.CreatedBy, color.Black)
	displayCreator.TextSize = 15
	displayCreator.Move(fyne.NewPos(10, 130))

	details := string(individual.Details)
	details = strings.Replace(individual.Details, "\\n", "\n", -1)
	displayDetails := widget.NewTextGridFromString(details)
	displayDetails.Move(fyne.NewPos(10, 180))

	//try and bind form

	formTitle := canvas.NewText("Modifier un producteur", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(400, 35))
	formTitle.Move(fyne.NewPos(10, 50))

	nameProducer := widget.NewEntry()
	nameProducer.SetText(individual.Name)
	nameProducer.Resize(fyne.NewSize(400, 35))
	nameProducer.Move(fyne.NewPos(10, 100))

	detailsProducer := widget.NewEntry()
	detailsProducer.SetText(details)
	detailsProducer.Resize(fyne.NewSize(400, 100))
	detailsProducer.Move(fyne.NewPos(10, 150))

	createdByProducer := widget.NewEntry()
	createdByProducer.SetText(individual.CreatedBy)
	createdByProducer.Resize(fyne.NewSize(400, 35))
	createdByProducer.Move(fyne.NewPos(10, 270))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(400, 50))
	submitBtn.Move(fyne.NewPos(10, 320))

	leftContainer := container.NewWithoutLayout(displayName, displayDetails, displayCreator)
	rightContainer := container.NewWithoutLayout(nameProducer, detailsProducer, createdByProducer, formTitle, submitBtn)
	mainContainer := container.New(layout.NewGridLayout(2))

	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)
	return mainContainer
}
