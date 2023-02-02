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
	"github.com/rohanthewiz/rtable"
	"negosud-gui/data"
	"os"
	"strconv"
	"strings"
	"time"
)

var BindBottles []binding.DataMap

func makeNewBottleTabs(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Liste des produits", displayBottle(nil)),
		container.NewTabItem("Ajouter un produit", addNewBottle(nil)),
		container.NewTabItem("En stock", displayStock(nil)),
		container.NewTabItem("En rupture de stock", displayInventory(nil)))
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

// displayAndUpdateBottle implements a dynamic table bound to an editing form
func displayBottle(_ fyne.Window) fyne.CanvasObject {
	var source = "WIDGETS.BOTTLE "

	// retrieve structs from data package
	Individual := data.IndBottle
	BottleData := data.BottleData

	var yPos, heightFields, widthForm float32
	var identifier string

	instructions, productImg, productName, productDesc, productLab, productVol, productYear, productPr, productAlc := createDetailsLabel(yPos, heightFields, widthForm)
	productDetails := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "", Widget: productName},
			{Text: "", Widget: productDesc},
			{Text: "", Widget: productLab},
			{Text: "", Widget: productVol},
			{Text: "", Widget: productYear},
			{Text: "", Widget: productPr},
			{Text: "", Widget: productAlc},
		},
	}

	// UPDATE FORM

	// declare form elements
	nameBottle := widget.NewEntry()
	detailsBottle := widget.NewMultiLineEntry()
	typeBottle := widget.NewEntry()
	volumeBottle := widget.NewEntry()
	alcoholBottle := widget.NewEntry()
	yearBottle := widget.NewEntry()
	priceBottle := widget.NewEntry()
	pictureBottle := widget.NewButtonWithIcon("Ajouter une image", theme.FileImageIcon(), func() { fmt.Print("Image was sent") })

	deleteBtn := widget.NewButtonWithIcon("Supprimer ce produit", theme.WarningIcon(), func() { fmt.Print("Deleting producer") })
	deleteBtn.Resize(fyne.NewSize(600, 50))

	response := data.AuthGetRequest("bottle")
	if response == nil {
		message := "Request body returned empty"
		fmt.Println(message)
		data.Logger(false, source, message)
		return widget.NewLabel("Le serveur n'a renvoyé aucun contenu")
	}
	if err := json.NewDecoder(response).Decode(&BottleData); err != nil {
		log(true, source, err.Error())
		fmt.Println(err)
	}

	for i := 0; i < len(BottleData); i++ {
		// converting 'int' to 'string' as rtable only accepts 'string' values
		id := strconv.Itoa(BottleData[i].Id)
		v := strconv.Itoa(BottleData[i].VolumeInt)
		a := strconv.Itoa(BottleData[i].AlcoholPercentage)
		p := strconv.Itoa(BottleData[i].CurrentPrice)
		y := strconv.Itoa(BottleData[i].YearProduced)
		BottleData[i].Price = p
		BottleData[i].Year = y
		BottleData[i].Volume = v
		BottleData[i].Alcohol = a
		BottleData[i].ID = id

		// binding bottle data
		BindBottles = append(BindBottles, binding.BindStruct(&BottleData[i]))
	}

	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: BottlesColumns,
		Bindings: BindBottles,
	}

	table := rtable.CreateTable(tableOptions)
	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 0 || cell.Row > len(BindBottles) { // 1st row is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(BottlesColumns) {
			fmt.Println("*-> Column out of limits")
			return
		}

		//Handle non-header row clicked
		identifier, err := rtable.GetStrCellValue(cell, tableOptions)
		if err != nil {
			fmt.Println(err.Error())
			log(true, source, err.Error())
			return
		}
		// Printout body cells
		rowBinding := tableOptions.Bindings[cell.Row-1]
		_, err = rowBinding.GetItem(tableOptions.ColAttrs[cell.Col].ColName)
		if err != nil {
			fmt.Println(err.Error())
			log(true, source, err.Error())
			return
		} else {
			instructions.Hidden = true
		}
		// Prevent app crash if other colunm than ID is clicked
		_, err = strconv.Atoi(identifier)
		if err != nil {
			resultApi := data.AuthGetRequest("bottle/" + identifier)
			if err := json.NewDecoder(resultApi).Decode(&Individual); err != nil {
				fmt.Println(err)
				log(true, source, err.Error())
			} else {
				productImg.Hidden = false
			}
			// Fill form fields with fetched data
			nameBottle.SetText(Individual.FullName)
			details := strings.Replace(Individual.Description, "\\n", "\n", -1)
			detailsBottle.SetText(details)
			typeBottle.SetText(Individual.WineType)
			volumeBottle.SetText(strconv.Itoa(Individual.Volume))
			yearBottle.SetText(strconv.Itoa(Individual.YearProduced))
			priceBottle.SetText(fmt.Sprintf("%f", Individual.CurrentPrice))
			alcoholBottle.SetText(fmt.Sprintf("%f", Individual.AlcoholPercentage))
			// Display details
			productName.SetText("Nom: " + Individual.FullName)
			productDesc.SetText("Description: " + details)
			productLab.SetText("Type : " + Individual.WineType)
			productYear.SetText("Année : " + strconv.Itoa(Individual.YearProduced))
			productVol.SetText("Volume : " + strconv.Itoa(Individual.Volume) + " cL")
			productPr.SetText("Prix HT : " + fmt.Sprintf("%f", Individual.CurrentPrice) + " €")
			productAlc.SetText("Alcool : " + fmt.Sprintf("%f", Individual.AlcoholPercentage) + " %")
		}

	}
	updateForm := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "Nom", Widget: nameBottle},
			{Text: "Description", Widget: detailsBottle},
			{Text: "Type", Widget: typeBottle},
			{Text: "Vol. (cL)", Widget: volumeBottle},
			{Text: "Alc. (%)", Widget: alcoholBottle},
			{Text: "Année", Widget: yearBottle},
			{Text: "Prix (€)", Widget: priceBottle},
			{Text: "", Widget: pictureBottle},
		},
		OnSubmit: func() {
			vol, _ := strconv.ParseInt(volumeBottle.Text, 10, 0)
			alc, _ := strconv.ParseInt(alcoholBottle.Text, 10, 0)
			year, _ := strconv.ParseInt(yearBottle.Text, 10, 0)
			pr, _ := strconv.ParseInt(priceBottle.Text, 10, 0)
			who, _ := os.Hostname()
			t, _ := time.Parse("2023-01-27T22:48:02.646Z", time.Now().String())
			bottle := &data.Bottle{
				FullName:          nameBottle.Text,
				Description:       detailsBottle.Text,
				WineType:          typeBottle.Text,
				Volume:            int(vol),
				AlcoholPercentage: float32(alc),
				CreatedAt:         t,
				UpdatedAt:         t,
				YearProduced:      int(year),
				CreatedBy:         who,
				UpdatedBy:         who,
				CurrentPrice:      float32(pr),
			}
			// Convert to JSON
			jsonValue, err := json.Marshal(bottle)
			if err != nil {
				log(true, source, err.Error())
				fmt.Println(err)
			}
			fmt.Print(bytes.NewBuffer(jsonValue))
			// Send data to API
			postData := data.AuthPostRequest("Bottle/UpdateBottle/"+identifier, bytes.NewBuffer(jsonValue))
			if postData != 200 {
				message := "Error on bottle " + identifier + " update " + " StatusCode " + strconv.Itoa(postData)
				log(true, source, message)
			}
			fmt.Println("Bottle updated")
		},
		OnCancel: func() {
			fmt.Println("Canceled")
		},
		SubmitText: "Envoyer",
		CancelText: "Annuler",
	}

	image := container.NewBorder(container.NewVBox(productImg), nil, nil, nil)
	textProduct := container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 100), productDetails))

	layoutDetailsTab := container.NewBorder(image, nil, nil, nil, textProduct, instructions)
	layoutUpdateForm := container.NewCenter(container.NewGridWrap(fyne.NewSize(600, 200), updateForm))
	layoutWithDelete := container.NewBorder(layoutUpdateForm, deleteBtn, nil, nil)

	individualTabs := container.NewAppTabs(
		container.NewTabItem("Détails du produit", layoutDetailsTab),
		container.NewTabItem("Modifier le produit", layoutWithDelete),
	)
	// LAYOUT
	mainContainer := container.New(layout.NewGridLayout(2))
	leftContainer := table
	rightContainer := container.NewBorder(nil, nil, nil, nil, individualTabs)
	mainContainer.Add(leftContainer)
	mainContainer.Add(rightContainer)

	return mainContainer
}

func createDetailsLabel(yPos float32, heightFields float32, widthForm float32) (*widget.Label, *canvas.Image, *widget.Label, *widget.Label, *widget.Label, *widget.Label, *widget.Label, *widget.Label, *widget.Label) {
	yPos = 180
	heightFields = 35
	widthForm = 600

	// DETAILS PRODUCT
	instructions := widget.NewLabelWithStyle("Cliquez sur un identifiant dans le tableau pour afficher les détails", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	instructions.Resize(fyne.NewSize(widthForm, heightFields))
	instructions.Move(fyne.NewPos(0, yPos-500))

	productImg := canvas.NewImageFromFile("media/bouteille.jpeg")
	productImg.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		productImg.SetMinSize(fyne.NewSize(250, 300))
	} else {
		productImg.SetMinSize(fyne.NewSize(250, 300))
	}
	productImg.Hidden = true

	productName := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productName.Resize(fyne.NewSize(widthForm, heightFields))
	productName.Move(fyne.NewPos(0, yPos-400))

	productDesc := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productDesc.Resize(fyne.NewSize(widthForm, heightFields))
	productDesc.Move(fyne.NewPos(0, yPos-350))

	productLab := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productLab.Resize(fyne.NewSize(widthForm, heightFields))
	productLab.Move(fyne.NewPos(0, yPos-300))

	productVol := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productVol.Resize(fyne.NewSize(widthForm, heightFields))
	productVol.Move(fyne.NewPos(0, yPos-250))

	productYear := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productYear.Resize(fyne.NewSize(widthForm, heightFields))
	productYear.Move(fyne.NewPos(0, yPos-200))

	productPr := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productPr.Resize(fyne.NewSize(widthForm, heightFields))
	productPr.Move(fyne.NewPos(0, yPos-150))

	productAlc := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	productAlc.Resize(fyne.NewSize(widthForm, heightFields))
	productAlc.Move(fyne.NewPos(0, yPos-100))

	return instructions, productImg, productName, productDesc, productLab, productVol, productYear, productPr, productAlc
}
