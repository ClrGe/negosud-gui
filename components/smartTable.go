package components

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
	"github.com/rohanthewiz/rtable"
	"net/http"
)

type PartProducer struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_By"`
}

var SmartCols = []rtable.ColAttr{
	{ColName: "Name", Header: "Nom", WidthPercent: 150},
	{ColName: "CreatedBy", Header: "Crée par", WidthPercent: 50},
}

var ProdData []PartProducer

var ProdBindings []binding.DataMap

func bindTable(w fyne.Window) fyne.CanvasObject {

	apiUrl := producerAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&ProdData); err != nil {
		fmt.Println(err)
	}
	fmt.Print(ProdData)

	for i := 0; i < len(ProdData); i++ {
		ProdBindings = append(ProdBindings, binding.BindStruct(&ProdData[i]))
	}
	tableOptions := &rtable.TableOptions{
		RefWidth: "========================================",
		ColAttrs: SmartCols,
		Bindings: ProdBindings,
	}

	table := rtable.CreateTable(tableOptions)

	table.OnSelected = func(cell widget.TableCellID) {
		if cell.Row < 0 || cell.Row > len(ProdBindings) { // 1st col is header
			fmt.Println("*-> Row out of limits")
			return
		}
		if cell.Col < 0 || cell.Col >= len(SmartCols) {
			fmt.Println("*-> Column out of limits")
			return
		}
		// Handle header row clicked
		if cell.Row == 0 { // fmt.Println("-->", tblOpts.ColAttrs[cell.Col].Header)
			// Add a row
			ProdBindings = append(ProdBindings,
				binding.BindStruct(&Producer{Name: "Belle Ambiance",
					Details: "brown", CreatedBy: "170"}))
			tableOptions.Bindings = ProdBindings
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
		err = cellBinding.(binding.String).Set(rvsString(str))
		if err != nil {
			fmt.Println(rerr.StringFromErr(err))
			return
		}
		fmt.Println("-->", str)
	}

	return table
}
func rvsString(in string) (out string) {
	runes := []rune(in)
	ln := len(runes)
	halfLn := ln / 2

	for i := 0; i < halfLn; i++ {
		runes[i], runes[ln-1-i] = runes[ln-1-i], runes[i]
	}
	return string(runes)
}
