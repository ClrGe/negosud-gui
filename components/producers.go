package components

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"time"
)

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

var producers []Producer

// Fetch API to return the list of all producers
func retrieveProducers(_ fyne.Window) fyne.CanvasObject {
	env, err := LoadConfig(".")
	res, err := http.Get(env.SERVER + "/api/producer")
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&producers); err != nil {
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
				label.SetText(fmt.Sprintf("%d", producers[id.Row].ID))
			case 1:
				label.SetText(producers[id.Row].Name)
			case 2:
				label.SetText(producers[id.Row].Details)
			case 3:
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
	table.SetColumnWidth(5, 200)

	table.SetRowHeight(2, 50)
	return table
}
