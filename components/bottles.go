package components

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"time"
)

// Define the producer struct and associate json fields
type Bottle struct {
	ID                int         `json:"id"`
	FullName          string      `json:"full_Name"`
	Description       string      `json:"description"`
	Label             string      `json:"label"`
	Volume            int         `json:"volume"`
	Picture           string      `json:"picture"`
	YearProduced      int         `json:"year_Produced"`
	AlcoholPercentage int         `json:"alcohol_Percentage"`
	CurrentPrice      int         `json:"current_Price"`
	CreatedAt         time.Time   `json:"created_At"`
	UpdatedAt         time.Time   `json:"updated_At"`
	CreatedBy         string      `json:"created_By"`
	UpdatedBy         string      `json:"updated_By"`
	BottleLocations   interface{} `json:"bottleLocations"`
	BottleGrapes      interface{} `json:"bottleGrapes"`
	Producer          interface{} `json:"producer"`
}

var bottles []Bottle

// Call bottle API and return the list of all bottles
func fetchBottles() {
	env, err := LoadConfig(".")

	res, err := http.Get(env.SERVER + "/api/bottle")
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&bottles); err != nil {
		fmt.Println(err)
	}
}

// Display API call result in a table
func displayBottles(_ fyne.Window) fyne.CanvasObject {
	fetchBottles()
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
				label.SetText(fmt.Sprintf("%v", bottles[id.Row].ID))
			case 1:
				label.SetText(bottles[id.Row].FullName)
			case 2:
				label.SetText(bottles[id.Row].Label)
			case 3:
				label.SetText(bottles[id.Row].CreatedBy)
			case 4:
				label.SetText(fmt.Sprintf("%v", bottles[id.Row].CurrentPrice))
			case 5:
				label.SetText(fmt.Sprintf("%v", bottles[id.Row].Volume))
			}
		})

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 200)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 200)
	table.SetRowHeight(2, 50)

	return table
}
