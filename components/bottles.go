package components

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"time"
)

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

// Fetch API to return the list of all bottles
func retrieveBottles(_ fyne.Window) fyne.CanvasObject {
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
				label.SetText(bottles[id.Row].FullName)
			}
		})
	return table
}
